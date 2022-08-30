package daemon

import (
	"context"
	"net"
	"net/http"
)

func NewSocketClient() (*http.Client, error) {
	sock, err := NewSocket()
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
