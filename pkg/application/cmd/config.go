package cmd

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/hhiroshell/kube-boat/pkg/infrastructure/socket"
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
	sock, err := socket.NewSocket()
	if err != nil {
		return err
	}

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sock.Path())
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
