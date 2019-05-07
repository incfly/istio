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
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"

	"istio.io/istio/istioctl/pkg/meshexp"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/log"
	"k8s.io/client-go/kubernetes"
)

var (
	labels      []string
	annotations []string
	svcAcctAnn  string
)

func createClients() (kubernetes.Interface, *crd.Client, error) {
	client, err := createInterface(kubeconfig)
	if err != nil {
		return nil, nil, err
	}
	seClient, err := crd.NewClient(
		kubeconfig, configContext,
		model.ConfigDescriptor{model.ServiceEntry},
		model.IstioAPIGroupDomain)
	if err != nil {
		return nil, nil, err
	}
	return client, seClient, nil
}

func register() *cobra.Command {
	registerCmd := &cobra.Command{
		Use:   "register <svcname> <ip> [protocol:]port1 [protocol:]port2 ...",
		Short: "Register a service instance (e.g. VM) joining the mesh",
		Long: `register creates Kuberentes services and ServiceEntry for mesh expansion services.

For example
Create a http service listening on port 8080, hosting on a VM with address as 10.0.0.1
	istioctl register httpbin-vm 10.0.0.1 http:8080

Create a service with three ports and in namespace foo, with service account bar:
	istioctl register http-grpc-vm 10.0.0.1 http:8080 grpc:9090 http:9000 -n foo -s bar
`,
		Args: cobra.MinimumNArgs(3),
		RunE: func(c *cobra.Command, args []string) error {
			client, seClient, err := createClients()
			if err != nil {
				return fmt.Errorf("failed to create client, check kubeconfig,kubecontext are correct, %v", err)
			}
			svcName := args[0]
			ip := args[1]
			portsListStr := args[2:]
			ports, err := meshexp.ConverPortList(portsListStr)
			if err != nil {
				return fmt.Errorf("failed to convert port list %v", err)
			}
			ns, _ := handleNamespaces(namespace)
			opts := &meshexp.VMServiceOpts{
				Name:           svcName,
				Namespace:      ns,
				PortList:       ports,
				IP:             []string{ip},
				ServiceAccount: svcAcctAnn,
			}
			se, err := meshexp.GetServiceEntry(opts)
			if err != nil {
				return err
			}
			svc, err := meshexp.GetKubernetesService(opts)
			if err != nil {
				return err
			}
			fmt.Printf("jianfeih debug \n%+v\n%+v\n", proto.MarshalTextString(svc), proto.MarshalTextString(se.Spec))
			if err := meshexp.Add(client, seClient, ns, se, svc); err != nil {
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
