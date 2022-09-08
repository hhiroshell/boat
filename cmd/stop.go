package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hhiroshell/boat/daemon"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running local Kubernetes API server",
	Long:  `Stop the running local Kubernetes API server`,
	RunE:  stop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stop(_ *cobra.Command, _ []string) error {
	client, err := daemon.NewClient()
	if err != nil {
		return err
	}

	msg, err := client.StopDaemon()
	if err != nil {
		return err
	}

	fmt.Println(msg)
	return nil
}
