// Package register implements mesh expansion service registry for istioctl.
package register

import (
	"fmt"

	"istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/config/kube/crd"
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
	IP             []string
	PortList       model.PortList
	Labels         map[string]string
	Annotations    map[string]string
}

// ConverPortList converst a list of string to the `model.PortList`.
func ConverPortList(ports []string) (model.PortList, error) {
	portList := model.PortList{}
	for _, p := range ports {
		np, err := kube_registry.Str2NamedPort(p)
		if err != nil {
			return nil, fmt.Errorf("invalid port format %v", p)
		}
		portList = append(portList, &model.Port{
			Port:     int(np.Port),
			Protocol: model.Protocol(np.Name),
		})
	}
	return portList, nil
}

// GetServiceEntry returns a service entry for mesh expansion service.
// TODO(incfly): change to model.Config such that the metadata is also included.
func GetServiceEntry(vs *VMServiceOpts) (*v1alpha3.ServiceEntry, error) {
	if vs == nil {
		return nil, fmt.Errorf("empty vm service options")
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
	for _, ip := range vs.IP {
		eps = append(eps, &v1alpha3.ServiceEntry_Endpoint{
			Address: ip,
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
		return nil, fmt.Errorf("empty vm service opts")
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

// ServiceEntryName returns the service entry crd name.
func ServiceEntryName(hostShortName string) string {
	return fmt.Sprintf("mesh-expansion-%v", hostShortName)
}

// Apply creates service entry and kubernetes service object in order to register vm service.
func Apply(client kubernetes.Interface, kubecfg, ns string,
	se *v1alpha3.ServiceEntry, svc *corev1.Service) error {
	// args.Config.ControllerOptions.DomainSuffix)
	if se == nil || svc == nil {
		return fmt.Errorf("failed to create vm service")
	}
	configClient, err := crd.NewClient(kubecfg, "", model.ConfigDescriptor{model.ServiceEntry}, "istio.io")
	if err != nil {
		return fmt.Errorf("failed to init client, check kubeconfig is provided %v", err)
	}
	// Pre-check Kubernetes service and service entry does not exist.
	_, err = client.CoreV1().Services(ns).Get(svc.Name, metav1.GetOptions{
		IncludeUninitialized: true,
	})
	if err == nil {
		return fmt.Errorf("service already exists, skip")
	}
	seName := ServiceEntryName(svc.Name)
	if oldServiceEntry := configClient.Get(
		model.ServiceEntry.Type, seName, ns); oldServiceEntry != nil {
		return fmt.Errorf("service entry already exists, skip")
	}
	// Create Kubernetes service and service entry.
	if _, err := client.CoreV1().Services(ns).Create(svc); err != nil {
		return fmt.Errorf("failed to create kuberenetes service %v", err)
	}
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:      model.ServiceEntry.Type,
			Group:     model.ServiceEntry.Group,
			Version:   model.ServiceEntry.Version,
			Name:      seName,
			Namespace: ns,
			Domain:    model.IstioAPIGroupDomain,
		},
		Spec: se,
	}
	if _, err := configClient.Create(config); err != nil {
		return fmt.Errorf("failed to create service entry %v", err)
	}
	return nil
}
