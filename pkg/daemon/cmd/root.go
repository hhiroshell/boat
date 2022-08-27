package cmd

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/hhiroshell/kube-boat/pkg/common"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:           "boat-daemon",
		Short:         "Start a standalone, local Kubernetes api server.",
		Long:          `Start a standalone, local Kubernetes api server.`,
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
	sock := filepath.Join(home, ".kube-boat", "daemon.socket")

	testEnv := &envtest.Environment{}
	config, err := testEnv.Start()
	if err != nil {
		return err
	}

	engine := gin.Default()
	engine.GET("/testenv", func(c *gin.Context) {
		c.JSON(http.StatusOK, common.KubeConfig{
			Server:     config.Host,
			ClientCert: base64.StdEncoding.EncodeToString(config.CertData),
			ClientKey:  base64.StdEncoding.EncodeToString(config.KeyData),
		})
	})
	engine.DELETE("/testenv", func(c *gin.Context) {
		if err := testEnv.Stop(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to stop the testenv server",
			})
		}

		if err := os.Remove(sock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to clean up the socket file. please delete \"./daemon.socket\"",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "stopped",
		})

		os.Exit(0)
	})

	if err := engine.RunUnix(sock); err != nil {
		fmt.Print(err)
		if err := testEnv.Stop(); err != nil {
			return err
		}
		return err
	}

	return nil
}
