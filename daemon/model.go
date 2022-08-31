package daemon

type Kubeconfig struct {
	Server     string `json:"server"`
	ClientCert string `json:"client-certificate-data"`
	ClientKey  string `json:"client-key-data"`
}

type Message struct {
	Message string `json:"message"`
}
