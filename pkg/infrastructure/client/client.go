package client

import (
	"context"
	"net"
	"net/http"

	"github.com/hhiroshell/kube-boat/pkg/infrastructure/socket"
)

func NewSocketClient() (*http.Client, error) {
	sock, err := socket.NewSocket()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sock.Path())
			},
		},
	}
	return client, nil
}
