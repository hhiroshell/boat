package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

type Client struct {
	client *http.Client
}

func NewClient() (*Client, error) {
	sc, err := newSocketClient()
	if err != nil {
		return nil, err
	}

	return &Client{client: sc}, nil
}

func newSocketClient() (*http.Client, error) {
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

func (c *Client) Kubeconfig() (*Kubeconfig, error) {
	res, err := c.client.Get("http://localhost/kube-boat")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	kubeconfig := &Kubeconfig{}
	if err := json.Unmarshal(body, kubeconfig); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal response from the kube-boat daemon: %w", err)
	}

	return kubeconfig, nil
}
