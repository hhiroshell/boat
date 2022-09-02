package daemon

import (
	"fmt"
	"os"
	"path/filepath"
)

type Socket struct {
	path string
}

func NewSocket() (*Socket, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get file path of unix socket: %w", err)
	}

	dir := filepath.Join(home, ".kube-boat")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0700)
	}

	path := filepath.Join(dir, "daemon.socket")

	return &Socket{path: path}, nil
}

func (s *Socket) Path() string {
	return s.path
}

func (s *Socket) Close() error {
	if err := os.Remove(s.path); err != nil {
		return fmt.Errorf("failed remove unix socket file: %w", err)
	}

	return nil
}
