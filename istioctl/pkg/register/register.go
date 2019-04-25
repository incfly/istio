// Package register implements mesh expansion service registry for istioctl.
package register

import (
	"fmt"

	"istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/model"
	kube_registry "istio.io/istio/pilot/pkg/serviceregistry/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// VMServiceOpts contains the options of a mesh exapnsion VM service.
type VMServiceOpts struct {
	Name           string
	Namespace      string
	ServiceAccount string
	// TODO: support more than one IPs when needed.
	IP          []string
	PortList    model.PortList
	Labels      map[string]string
	Annotations map[string]string
}

// ConverPortList converst a list of string to the `model.PortList`.
func ConverPortList(ports []string) (model.PortList, error) {
	return nil, nil
}

// GetServiceEntry returns a service entry for mesh expansion service.
func GetServiceEntry(vs *VMServiceOpts) (*v1alpha3.ServiceEntry, error) {
	if vs == nil {
		return nil, fmt.Errorf("empty VmServiceOpts")
	}
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
	host := fmt.Sprintf("%v.%v.svc.cluster.local", vs.Name, vs.Namespace)
	return &v1alpha3.ServiceEntry{
		Hosts:      []string{host},
		Ports:      ports,
		Endpoints:  eps,
		Resolution: v1alpha3.ServiceEntry_STATIC,
	}, nil
}

// GetKubernetesService returns the Kubernetes service object.
func GetKubernetesService(vs *VMServiceOpts) (*corev1.Service, error) {
	if vs == nil {
		return nil, fmt.Errorf("empty VmServiceOpts")
	}
	ports := []corev1.ServicePort{}
	for _, p := range vs.PortList {
		ports = append(ports, corev1.ServicePort{
			Name: p.Name,
			Port: int32(p.Port),
		})
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vs.Name,
			Namespace: vs.Namespace,
			Annotations: map[string]string{
				kube_registry.KubeServiceAccountsOnVMAnnotation: vs.ServiceAccount,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: ports,
		},
	}, nil
}

// Apply creates service entry and kubernetes service object in order to register vm service.
func Apply(client *kubernetes.Clientset, se *v1alpha3.ServiceEntry, svc *corev1.Service) error {
	if se == nil || svc == nil {
		return fmt.Errorf("failed to create vm service")
	}
	return nil
}
