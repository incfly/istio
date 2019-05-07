package register

import (
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v2"
	"istio.io/api/networking/v1alpha3"
)

var (
	output = `{
  "hosts": [
    "httpbin-vm.default.svc.cluster.local"
  ],
  "resolution": 1,
  "endpoints": [
    {
      "address": "10.0.0.1"
    }
  ]
}`
)

func TestGetServiceEntry(t *testing.T) {
	opt := &VMServiceOpts{
		Name:           "vmhttp",
		Namespace:      "default",
		ServiceAccount: "foo",
	}
	want := v1alpha3.ServiceEntry{
		Hosts: []string{"vmhttp.default.svc.cluster.local"},
	}
	got, err := GetServiceEntry(opt)
	if err != nil {
		t.Errorf("failed to returned a service entry, expect succeed, err %v", err)
	}
	gotBytes, err := yaml.Marshal(got)
	if err != nil {
		t.Errorf("failed to convert to yaml %v", err)
	}
	wantBytes, err := yaml.Marshal(want)
	if err != nil {
		t.Errorf("failed to convert to yaml %v", err)
	}
	// Compare to golden service entry.
	if !reflect.DeepEqual(gotBytes, want) {
		t.Errorf("unexpected service entry, got %v, want %v", string(gotBytes), wantBytes)
	}
}

// func TestGetKubernetesService(t *testing.T) {
// 	opt := &VMServiceOpts{
// 		Name:           "vmhttp",
// 		Namespace:      "default",
// 		ServiceAccount: "foo",
// 	}
// 	got, err := GetKubernetesService(opt)
// 	if err != nil {
// 		t.Errorf("failed to return a kubernetes service, expect suceed, err %v", er)
// 	}
// 	out, err := yaml.Marshal(got)
// 	if err != nil {
// 		t.Errorf("failed to convert to yaml", err)
// 	}
// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("unexpected service entry, got %v, want %v", got, want)
// 	}
// }
