// Copyright © 2017 Matthew Farina <matt@mattfarina.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string // Path to the kube config file

// The description is a variable to make editing and keeping the max width
// easier in an editor that has guides.
var rootDesc = `
This is a demo application illustrating how to interact with Kubernetes, the
open-source system for container management. This application enables pods to be
listed and deleted.
`

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k8s-demo-cli",
		Short: "A kubernetes demo CLI",
		Long:  rootDesc,
	}

	f := cmd.PersistentFlags()
	f.StringVar(&kubeconfig, "config", "", "Location of the Kubernetes config file (defaults to ~/.kube/config)")

	client, err := kubernetesClient()
	if err != nil {
		log.Fatalf("Unable to locate Kubernetes config file. err: %s", err)
		os.Exit(1)
	}

	out := cmd.OutOrStdout()

	cmd.AddCommand(
		newListCmd(out, client),
		newDeleteCmd(out, client),
	)

	return cmd
}

func main() {

	cmd := newRootCmd()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Create a client that can talk with a Kubernetes cluster. Get the details for
// the cluster location from the kubectl (local) configuration file.
// TODO(mattfarina): Provide the ability to pass in the cluster location and auth
// as flags or environment variables.
func kubernetesClient() (*kubernetes.Clientset, error) {

	cfg := kubeconfig
	if cfg == "" {

		// Using homedir because it can location the home directory without
		// using libc.
		d, err := homedir.Dir()
		if err != nil {
			return nil, err
		}

		// Using filepath.Join so that the path separator works for both *nix
		// and Windows systems.
		cfg = filepath.Join(d, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", cfg)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
