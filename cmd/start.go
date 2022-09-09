package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
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
	boat, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(boat, "serve")
	if err := cmd.Start(); err != nil {
		return err
	}
	fmt.Println("Starting local Kubernetes API server...")

	// TODO: wait for the daemon become running
	return nil
}
