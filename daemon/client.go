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

func (c *Client) Readyz() error {
	res, err := c.client.Get("http://localhost" + readyz)
	if err != nil {
		return fmt.Errorf("kube-boat daemon or local Kubernetes API Server is not ready: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	} else {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("kube-boat daemon or local Kubernetes API Server is not ready: %w", err)
		}
		return fmt.Errorf("kube-boat daemon or local Kubernetes API Server is not ready: %s", string(body))
	}
}

func (c *Client) Kubeconfig() (*Kubeconfig, error) {
	res, err := c.client.Get("http://localhost" + kubeconfig)
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

func (c *Client) WebhookConfig() (*WebhookConfig, error) {
	res, err := c.client.Get("http://localhost" + webhookConfig)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	webhookConfig := &WebhookConfig{}
	if err := json.Unmarshal(body, webhookConfig); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal response from the kube-boat daemon: %w", err)
	}

	return webhookConfig, nil
}

func (c *Client) StopDaemon() (string, error) {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost"+base, nil)
	if err != nil {
		return "", err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	msg := &Message{}
	if err := json.Unmarshal(body, msg); err != nil {
		return "", fmt.Errorf("failed to Unmarshal response from the kube-boat daemon: %w", err)
	}

	return msg.Message, nil
}
