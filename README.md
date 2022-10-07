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

Quick Start
---
Make sure [environment variables for Envtest Binaries](https://book.kubebuilder.io/reference/envtest.html#environment-variables)
has been set appropriately.
Kube Boat executes Kubernetes API Server using Envtest framework. So these variables affects to behavior of the Server.

Then, you can start a Kubernetes API Server by `boat start`.

```console
$ boat start
Starting local Kubernetes API server...
 🚤 🚤 🚤 🚤 🚤 🚤 🚤 🚤 🚤 🚤
...Done.
```

Your kubectl context can be updated by `kubeconfig` sub command.

```console
$ boat kubeconfig

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
$ boat stop
shutting down the API server...
```

Advanced Usage
---

### Start with CRDs installed
When you start local Kubernetes API Server, you can set paths to the directory or file containing manifests of
CustomResourceDefinition with `--crd-path` flag.

With this flag, the API Server will start with the CRDs installed.

```console
$ boat start --crd-path=./example/crd/crontabs.stable.example.com.yaml

$ kubectl api-resources --api-group=stable.example.com
NAME       SHORTNAMES   APIVERSION              NAMESPACED   KIND
crontabs   ct           stable.example.com/v1   true         CronTab
```

LICENSE
---
Kube Boat is licensed under the Apache License 2.0, and includes works distributed under same one or the MIT License.
