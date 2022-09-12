package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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
	boat, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(boat, "serve")
	if err := cmd.Start(); err != nil {
		return err
	}
	fmt.Println("Starting local Kubernetes API server...")

	client, err := daemon.NewClient()
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Print("time out")
		case <-ticker.C:
			fmt.Print(" 🚤")
			if err := client.Readyz(); err == nil {
				fmt.Println("\n...Done.")
				return nil
			}
		}
	}

	return nil
}
