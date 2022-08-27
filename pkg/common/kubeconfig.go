package common

type KubeConfig struct {
	Server     string `json:"server"`
	ClientCert string `json:"client-certificate-data"`
	ClientKey  string `json:"client-key-data"`
}
