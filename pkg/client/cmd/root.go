package cmd

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:           "boat",
		Short:         "",
		Long:          ``,
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE:          run,
	}

	err := rootCmd.Execute()
	if err != nil {
		//Log.Info("exit with error", zap.Error(err))
		os.Exit(1)
	}
}

func run(_ *cobra.Command, _ []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", filepath.Join(home, ".kube-boat", "daemon.socket"))
			},
		},
	}

	res, err := client.Get("http://localhost/testenv")
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
