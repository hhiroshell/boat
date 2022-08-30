package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/hhiroshell/kube-boat/pkg/daemon"
)

var stopCmd = &cobra.Command{
	Use:           "stop",
	Short:         "",
	Long:          ``,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          stop,
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
