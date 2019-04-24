// Package register implements mesh expansion service registry for istioctl.
package register

import (
	"istio.io/api/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
)

// VMServiceOpts contains the options of a mesh exapnsion VM service.
type VMServiceOpts struct {
	Name     string
	PortList []string
	Labels   []string
	IP       []string
}

// GetServiceEntry returns a service entry for mesh expansion service.
func GetServiceEntry(vs *VMServiceOpts) (*v1alpha3.ServiceEntry, error) {
	se := v1alpha3.ServiceEntry{
		Hosts: []string{vs.Name},
	}
	for _, p := range vs.PortList {
		se.Ports = append(se.Ports, &v1alpha3.Port{
			Number:   123,
			Protocol: p,
		})
	}
	for _, endpoint := range vs.IP {
		se.Endpoints = append(se.Endpoints, &v1alpha3.ServiceEntry_Endpoint{
			Address: endpoint,
		})
	}
	return &se, nil
}

// GetKubernetesService returns the Kubernetes service object based on the
func GetKubernetesService() (*corev1.Service, error) {
	svc := &corev1.Service{}
	return svc, nil
}
