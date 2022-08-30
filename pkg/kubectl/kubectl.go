package kubectl

import (
	"os/exec"
)

func SetUser(cert, key string) error {
	setCertCmd := exec.Command("kubectl", "config", "set", "users.kube-boat.client-certificate-data", cert)
	if err := setCertCmd.Run(); err != nil {
		return err
	}

	setKeyCmd := exec.Command("kubectl", "config", "set", "users.kube-boat.client-key-data", key)
	if err := setKeyCmd.Run(); err != nil {
		return err
	}

	return nil
}

func SetCluster(server string) error {
	cmd := exec.Command("kubectl", "config", "set-cluster", "kube-boat", "--server="+server, "--insecure-skip-tls-verify=true")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
