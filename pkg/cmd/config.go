package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"

	"github.com/hhiroshell/kube-boat/pkg/daemon"
)

var configCmd = &cobra.Command{
	Use:           "config",
	Short:         "",
	Long:          ``,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          config,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func config(_ *cobra.Command, _ []string) error {
	c, err := daemon.NewSocketClient()
	if err != nil {
		return err
	}

	res, err := c.Get("http://localhost/kube-boat")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(res.Status)
	fmt.Println(string(body))

	return nil
}
