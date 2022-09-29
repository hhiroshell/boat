package cmd

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/hhiroshell/boat/daemon"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of local Kubernetes API server",
	Long:  "Check the status of local Kubernetes API server",
	RunE:  status,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

const (
	boatd   = "kube-boat daemon"
	kubeapi = "Kubernetes API Server"
)

func status(_ *cobra.Command, _ []string) error {
	client, err := daemon.NewClient()
	if err != nil {
		return err
	}

	status := map[string]string{
		boatd:   "⛔ Not yet running",
		kubeapi: "⛔ Not yet running",
	}
	defer func() {
		for k, v := range status {
			fmt.Println(k + ": " + v)
		}
	}()

	err = client.Readyz()
	if err == nil {
		status[boatd] = "✅ Running"
		status[kubeapi] = "✅ Running"
		return nil
	}

	var errno syscall.Errno
	if !errors.As(err, &errno) {
		// Readyz endpoint of Kube API Server returned some unhealthy status.
		status[boatd] = "✅ Running"
		status[kubeapi] = "⚠️ Unhealthy"
		return nil
	}

	return nil
}
