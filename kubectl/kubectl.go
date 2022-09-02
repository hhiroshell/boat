package kubectl

import (
	"os/exec"
)

func SetContext(server, cert, key string, use bool) error {
	setCertCmd := exec.Command(
		"kubectl",
		"config",
		"set",
		"users.kube-boat.client-certificate-data", cert,
	)
	if err := setCertCmd.Run(); err != nil {
		return err
	}

	setKeyCmd := exec.Command(
		"kubectl",
		"config",
		"set",
		"users.kube-boat.client-key-data", key,
	)
	if err := setKeyCmd.Run(); err != nil {
		return err
	}

	setClusterCmd := exec.Command(
		"kubectl",
		"config",
		"set-cluster",
		"kube-boat",
		"--server="+server,
		"--insecure-skip-tls-verify=true",
	)
	if err := setClusterCmd.Run(); err != nil {
		return err
	}

	setContextCommand := exec.Command(
		"kubectl",
		"config",
		"set-context",
		"kube-boat",
		"--cluster=kube-boat",
		"--user=kube-boat",
		"--namespace=default",
	)
	if err := setContextCommand.Run(); err != nil {
		return err
	}

	if use {
		useContextCmd := exec.Command(
			"kubectl",
			"config",
			"use-context",
			"kube-boat",
		)
		if err := useContextCmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
