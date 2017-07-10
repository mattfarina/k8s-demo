# Kubernetes Demo CLI

This is a simple CLI demo application that does two things:

1. Provides the ability to list the pods in a Kubernetes cluster. Includes limited filtering capabilities.
2. One or more pods can be deleted.

## Assumptions

1. All pods available to a user, without filtering of any kind, should be included by default. That includes the system pods.
2. The Kubernetes go client can be used to interact with the API.
3. Pods are the only thing deleted. New pods may be created to keep the application functioning.

## Design Decisions

This section lists the reasons behind design decisions within the code. They are listed in the Readme to consolidate them to one easy to find location.

1. The use of `github.com/spf13/cobra` for the CLI package is due to it being used in the Kubernetes project and related projects, such as Kubeless. It is being used for consistency.
2. By default, `cobra` CLI commands setup a project to use `github.com/spf13/viper`. Opting out because `viper` uses `fsnotify` which requires libc. This means if you use `viper` you cannot compile an application using `$  CGO_ENABLED=0 go build` (cgo disabled).
3. To find the home directory the `github.com/mitchellh/go-homedir` package is used as it avoids the need for cgo.
4. All pods are lists including those in the "kube-system" namespace. No filtering is provided by default. All filtering is opt-in.

## Installation

The application can be installed using `go get`. For example, `go get -u github.com/mattfarina/k8s-demo/cmd/k8s-demo-cli`. This will download the application into your `GOPATH` and install it.

## Dependencies

[kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) is installed and you have authenticated to a cluster. The current cluster is the one the application will interact with.

For development, [Glide](https://glide.sh) is used to manage the application dependencies. The dependencies themselves have been vendored to the `vendor` folder so the application can be installed with `go get`.
