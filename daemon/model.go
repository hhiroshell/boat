package daemon

type Kubeconfig struct {
	Server     string `json:"server"`
	ClientCert string `json:"client-certificate-data"`
	ClientKey  string `json:"client-key-data"`
}

type WebhookConfig struct {
	LocalServingHost    string `json:"local-serving-host"`
	LocalServingPort    int    `json:"local-serving-port"`
	LocalServingCertDir string `json:"local-serving-cert-dir"`
}

type Message struct {
	Message string `json:"message"`
}
