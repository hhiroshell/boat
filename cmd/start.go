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

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start local Kubernetes API server",
	Long:  `Start local Kubernetes API server`,
	RunE:  start,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func start(_ *cobra.Command, _ []string) error {
	client, err := daemon.NewClient()
	if err != nil {
		return err
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
		return err
	}

	cmd := exec.Command(boat, "serve")
	if err := cmd.Start(); err != nil {
		return err
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
			return errors.New("failed to start kube-boat daemon or Kubernetes API Server before timed out")
		case <-ticker.C:
			fmt.Print(" 🚤")
			if err := client.Readyz(); err == nil {
				fmt.Println("\n...Done.")
				return nil
			}
		}
	}
}
