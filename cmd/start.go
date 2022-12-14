package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/hhiroshell/boat/daemon"
)

var (
	updateKubeconfig bool

	startCmd = &cobra.Command{
		Use:          "start",
		Short:        "Start local Kubernetes API server",
		Long:         `Start local Kubernetes API server`,
		SilenceUsage: true,

		RunE: start,
	}
)

func init() {
	rootCmd.AddCommand(startCmd)

	setServeFlags(startCmd)
	startCmd.Flags().BoolVar(&updateKubeconfig, "update-kubeconfig", true, "Set the kubectl context to point to kube-apiserver started by this tool")
}

func start(_ *cobra.Command, _ []string) error {
	client, err := daemon.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create kube-boat daemon client: %w", err)
	}

	err = client.Readyz()
	if err == nil {
		return errors.New("local Kubernetes API Server is already running")
	}
	var errno syscall.Errno
	if !errors.As(err, &errno) {
		return errors.New("local Kubernetes API Server is already running and is in an unhealthy state")
	}

	boat, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failted to find kube-boat executable to start kube-boat daemon: %w", err)
	}

	serveOptions := []string{"serve"}
	for _, path := range crdPaths {
		serveOptions = append(serveOptions, "--"+crdPathFlag+"="+path)
	}
	for _, path := range webhookConfigurationPaths {
		serveOptions = append(serveOptions, "--"+webhookConfigurationPathFlag+"="+path)
	}

	serve := exec.Command(boat, serveOptions...)
	if err := serve.Start(); err != nil {
		return fmt.Errorf("failted to start kube-boat daemon or kube-apiserver: %w", err)
	}
	fmt.Println("Starting local Kubernetes API server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println()
			return errors.New("failed to start kube-boat daemon or kube-apiserver before timed out")
		case <-ticker.C:
			fmt.Print(" 🚤")
			if err := client.Readyz(); err == nil {
				fmt.Println("\n...Done.")
				goto exit
			}
		}
	}

exit:

	if updateKubeconfig {
		if err := setKubectlContext(client); err != nil {
			return fmt.Errorf("failed to update kubeconfig: %w", err)
		}
	}

	return nil
}
