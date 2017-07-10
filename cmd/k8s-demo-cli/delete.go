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

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

// The description is a variable to make editing and keeping the max width
// easier in an editor that has guides.
var deleteDesc = `
The delete command deletes one or more pods. If a namespace is passed in the
delete will be sure to scope the delete request to the namespace. Otherwise
the namespace will be automatically determined.

For example,

	$ k8s-demo-cli delete hello-world-938614450-7jfcg
	Deleted the pod "hello-world-938614450-7jfcg" in namespace "default"

Here the namespace was looked up for the pod name. Alternately, the namespace
can be passed in. For example,

	$ k8s-demo-cli delete --namespace default hello-world-938614450-7jfcg
	Deleted the pod "hello-world-938614450-7jfcg" in namespace "default"

Multiple pods can be deleted with a single command and they can be across
multiple namespaces. For example,

	$ k8s-demo-cli delete hello-world-938614450-7jfcg hello-world-1484529432-x2jck
	Deleted the pod "hello-world-938614450-7jfcg" in namespace "default"
	Deleted the pod "hello-world-1484529432-x2jck" in namespace "example"
`

type deleteCmd struct {
	out       io.Writer
	client    *kubernetes.Clientset
	name      string
	namespace string
	cacheMap  map[string]string
}

// newDeleteCmd is a constructor to create a cobra command to delete pods
// in a Kubernetes cluster
func newDeleteCmd(out io.Writer, client *kubernetes.Clientset) *cobra.Command {

	delete := &deleteCmd{
		out:      out,
		client:   client,
		cacheMap: make(map[string]string),
	}

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Short:   "Delete a pod running in a Kubernetes cluster",
		Long:    deleteDesc,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			ensureEnvFlag("namespace", cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete.run(args)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&delete.namespace, "namespace", "n", "", "Namespace to look for the pod in")

	return cmd
}

func (d *deleteCmd) run(names []string) error {
	options := &v1.DeleteOptions{}

	var ns string
	var found bool

	// A simple first pass at deleting pods is done in serial
	// TODO(mattfarina): Deleting multiplt pods can be done concurrently
	for _, v := range names {

		ns = d.namespace
		if ns == "" {
			// A namespace is required for deleting and none was passed in. So,
			// we figure it out.
			ns, found = d.findNamespace(v)
			if !found {
				return fmt.Errorf("Unable to delete %q. Pod not found", v)
			}
		}

		err := d.client.CoreV1().Pods(ns).Delete(v, options)
		if err != nil {
			return fmt.Errorf("Unable to delete %q. %s", v, err)
		}

		fmt.Fprintf(d.out, "Deleted the pod %q in namespace %q\n", v, ns)
	}

	return nil
}

// Pass in a name of a pod and a namespace along with if a namespace was found
// are returned.
func (d *deleteCmd) findNamespace(name string) (string, bool) {
	if len(d.cacheMap) == 0 {
		// Load the namespace to pod mapping as a cache
		pods, err := d.client.CoreV1().Pods("").List(v1.ListOptions{})
		if err != nil {
			fmt.Fprintf(d.out, "Error listing pods: %s\n", err)
			return "", false
		}

		for _, v := range pods.Items {
			d.cacheMap[v.Name] = v.Namespace
		}
	}

	ret, f := d.cacheMap[name]
	return ret, f
}
