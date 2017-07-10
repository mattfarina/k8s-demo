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
	"io"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

// The description is a variable to make editing and keeping the max width
// easier in an editor that has guides.
var listDesc = `
The list command provides the ability to list the pods within a cluster.

Using this command without any flags will list all of the pods within a cluster.
There are flags to filter based on namespace and labels.

For example, to filter based on namespace:

    $ k8s-demo-cli list --namespace default

This will filter to the pods in the default namespace

Filtering based on label is another option. For example:

	$ k8s-demo-cli list --labels "version=v1.6.1"

In this case the label key of version is equal to v1.6.1. Filtering can use the
=, ==, and != as possible options. Multiple labels can be provided in a comma
separated list.
`

type listCmd struct {
	out       io.Writer
	client    *kubernetes.Clientset
	namespace string
	labels    string
}

// newListCmd is a constructor to create a cobra command to list the pods
// in a Kubernetes cluster
func newListCmd(out io.Writer, client *kubernetes.Clientset) *cobra.Command {

	list := &listCmd{
		out:    out,
		client: client,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the pods running in a Kubernetes cluster",
		Long:    listDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&list.namespace, "namespace", "n", "", "Namespace to use as a filter")
	f.StringVarP(&list.labels, "labels", "l", "", "Labels to use as a filter")

	return cmd
}

func (l *listCmd) run() error {

	// TODO(mattfarina): Add more options for filtering lists
	options := v1.ListOptions{
		LabelSelector: l.labels,
	}
	pods, err := l.client.CoreV1().Pods(l.namespace).List(options)

	// TODO (mattfarina): The error provided here is not going to be user friendly.
	// Figure out how to provide a more user friendly message.
	if err != nil {
		return err
	}

	if len(pods.Items) == 0 {
		fmt.Fprintln(l.out, "No pods found")
		return nil
	}

	// For easier viewing, using a table
	table := uitable.New()
	table.Wrap = true
	table.AddRow("NAME", "NAMESPACE", "LABELS")
	for _, p := range pods.Items {

		// Get the labels as a string
		var tempL []string
		for k, v := range p.Labels {
			tempL = append(tempL, fmt.Sprintf("%s=%s", k, v))
		}

		table.AddRow(p.Name, p.Namespace, strings.Join(tempL, ", "))
	}

	fmt.Fprintln(l.out, table.String())

	return nil
}
