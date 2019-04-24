// Package register implements mesh expansion service registry for istioctl.
package register

import (
	"fmt"

	"istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/model"
	corev1 "k8s.io/api/core/v1"
)

// VMServiceOpts contains the options of a mesh exapnsion VM service.
type VMServiceOpts struct {
	Name        string
	PortList    model.PortList
	Labels      map[string]string
	IP          []string
	Annotations []string
}

// ConverPortList converst a list of string to the `model.PortList`.
func ConverPortList(ports []string) (model.PortList, error) {
	return nil, nil
}

// GetServiceEntry returns a service entry for mesh expansion service.
func GetServiceEntry(vs *VMServiceOpts) (*v1alpha3.ServiceEntry, error) {
	ports := []*v1alpha3.Port{}
	for _, p := range vs.PortList {
		ports = append(ports, &v1alpha3.Port{
			Number:   uint32(p.Port),
			Protocol: string(p.Protocol),
			Name:     p.Name,
		})
	}
	eps := []*v1alpha3.ServiceEntry_Endpoint{}
	for _, endpoint := range vs.IP {
		eps = append(eps, &v1alpha3.ServiceEntry_Endpoint{
			Address: endpoint,
			Labels:  vs.Labels,
		})
	}
	return &v1alpha3.ServiceEntry{
		Hosts:      []string{vs.Name},
		Ports:      ports,
		Endpoints:  eps,
		Resolution: v1alpha3.ServiceEntry_STATIC,
	}, nil
}

// GetKubernetesService returns the Kubernetes service object based on the
func GetKubernetesService(vs *VMServiceOpts) (*corev1.Service, error) {
	return &corev1.Service{}, nil
}

// Apply creates service entry and kubernetes service object in order to register vm service.
func Apply(se *v1alpha3.ServiceEntry, svc *corev1.Service) error {
	if se == nil || svc == nil {
		return fmt.Errorf("failed to create vm service")
	}
	return nil
}
