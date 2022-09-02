🚤 Kube Boat 🚤
===
![ci](https://github.com/hhiroshell/kube-boat/actions/workflows/ci.yaml/badge.svg)

Local and standalone Kubernetes API Server, always by your side.


Setup
---

### Requirements
Kube Boat requires following tools are installed in local machines.

- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- [Evntest Binaries](https://book.kubebuilder.io/reference/envtest.html)

### Installation

```console
$ go install github.com/hhiroshell/kube-boat@latest
```

Usage
---
Make sure [environment variables for Envtest Binaries](https://book.kubebuilder.io/reference/envtest.html#environment-variables)
has been set appropriately.
Kube Boat executes Kubernetes API Server using Envtest framework. So these variables affects to behavior of the Server.

Then, you can start a Kubernetes API Server localy by `kube-boat start`.

```console
$ kube-boat start
Starting local Kubernetes API server...
```

Your kubectl context can be updated by `kubeconfig` sub command.

```console
$ kube-boat kubeconfig

$ kubectl config current-context
kube-boat

$ kubectl get namespaces
NAME              STATUS   AGE
default           Active   22s
kube-node-lease   Active   27s
kube-public       Active   27s
kube-system       Active   27s
```

To stop the local API Server, use `stop` sub command.

```console
$ kube-boat stop
shutting down the API server...
```

LICENSE
---
Kube Boat is licensed under the Apache License 2.0, and includes works distributed under same one or the MIT License.
