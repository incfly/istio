package meshexp

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	messagediff "gopkg.in/d4l3k/messagediff.v1"
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
	for _, tc := range []struct {
		name   string
		vmopts VMServiceOpts
		want   string
	}{
		{
			name:   "basic",
			vmopts: VMServiceOpts{},
			want: `{
"hosts": [
	"httpbin-vm.default.svc.cluster.local"
],
"resolution": 1,
"endpoints": [
	{
		"address": "10.0.0.1"
	}
]
}`,
		},
	} {
		se, err := GetServiceEntry(&tc.vmopts)
		if err != nil {
			t.Errorf("[%v] failed %v", tc.name, err)
		}
		got := proto.MarshalTextString(se.Spec)
		if diff, equal := messagediff.PrettyDiff(got, tc.want); !equal {
			t.Errorf("[%v] unexpected service entry %v", tc.name, diff)
		}
	}
}

func TestGetKubernetesService(t *testing.T) {
	for _, tc := range []struct {
		name   string
		vmopts VMServiceOpts
		want   string
	}{
		{
			name:   "basic",
			vmopts: VMServiceOpts{},
			want: `{
"hosts": [
	"httpbin-vm.default.svc.cluster.local"
],
"resolution": 1,
"endpoints": [
	{
		"address": "10.0.0.1"
	}
]
}`,
		},
	} {
		svc, err := GetKubernetesService(&tc.vmopts)
		if err != nil {
			t.Errorf("[%v] failed %v", tc.name, err)
		}
		got := proto.MarshalTextString(svc)
		if diff, equal := messagediff.PrettyDiff(got, tc.want); !equal {
			t.Errorf("[%v] unexpected service entry %v", tc.name, diff)
		}
	}
}

// 	opt := &VMServiceOpts{
// 		Name:           "vmhttp",
// 		Namespace:      "default",
// 		ServiceAccount: "foo",
// 	}
// 	want := v1alpha3.ServiceEntry{
// 		Hosts: []string{"vmhttp.default.svc.cluster.local"},
// 	}
// 	got, err := GetServiceEntry(opt)
// 	if err != nil {
// 		t.Errorf("failed to returned a service entry, expect succeed, err %v", err)
// 	}
// 	gotBytes, err := yaml.Marshal(got)
// 	if err != nil {
// 		t.Errorf("failed to convert to yaml %v", err)
// 	}
// 	wantBytes, err := yaml.Marshal(want)
// 	if err != nil {
// 		t.Errorf("failed to convert to yaml %v", err)
// 	}
// 	// Compare to golden service entry.
// 	if !reflect.DeepEqual(gotBytes, want) {
// 		t.Errorf("unexpected service entry, got %v, want %v", string(gotBytes), wantBytes)
// 	}
// }
