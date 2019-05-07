// Copyright 2017 Istio Authors
//
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

package cmd

import (
	"github.com/spf13/cobra"

	"istio.io/istio/istioctl/pkg/meshexp"
	"istio.io/istio/pkg/log"
)

var (
	labels      []string
	annotations []string
	svcAcctAnn  string
)

func register() *cobra.Command {
	registerCmd := &cobra.Command{
		Use:   "register <svcname> <ip> [name1:]port1 [name2:]port2 ...",
		Short: "Registers a service instance (e.g. VM) joining the mesh",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(c *cobra.Command, args []string) error {
			svcName := args[0]
			ip := args[1]
			portsListStr := args[2:]
			ports, err := meshexp.ConverPortList(portsListStr)
			if err != nil {
				log.Errorf("failed to convert port list %v", err)
			}
			client, err := createInterface(kubeconfig)
			if err != nil {
				return err
			}
			ns, _ := handleNamespaces(namespace)
			opts := &meshexp.VMServiceOpts{
				Name:           svcName,
				Namespace:      ns,
				PortList:       ports,
				IP:             []string{ip},
				ServiceAccount: "default",
			}
			se, err := meshexp.GetServiceEntry(opts)
			if err != nil {
				return err
			}
			svc, err := meshexp.GetKubernetesService(opts)
			if err != nil {
				return err
			}
			// fmt.Printf("jianfeih debug \n%+v\n%+v\n", proto.MarshalTextString(svc), proto.MarshalTextString(se))
			if err := meshexp.Add(client, kubeconfig, ns, se, svc); err != nil {
				log.Errorf("failed to create service enetry and k8s svc: %v", err)
			}
			return nil
		},
	}

	registerCmd.PersistentFlags().StringSliceVarP(&labels, "labels", "l",
		nil, "List of labels to apply if creating a service/endpoint; e.g. -l env=prod,vers=2")
	registerCmd.PersistentFlags().StringSliceVarP(&annotations, "annotations", "a",
		nil, "List of string annotations to apply if creating a service/endpoint; e.g. -a foo=bar,test,x=y")
	registerCmd.PersistentFlags().StringVarP(&svcAcctAnn, "serviceaccount", "s",
		"default", "Service account to link to the service")

	return registerCmd
}
