// Package meshexp contains test suite for mesh expansion.
package meshexp

import (
	"fmt"
	"strings"
	"testing"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/environment"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/rawvm"
	"istio.io/istio/pkg/test/framework/resource"
)

var (
	vmInstance rawvm.Instance
	// ports is the port the VM service exposes.
	ports = model.PortList{
		&model.Port{
			Name:     "http",
			Port:     8080,
			Protocol: model.ProtocolHTTP,
		},
	}
)

const (
	serviceName = ""
)

// How to run this test suite locally
// go test -v ./tests/integration/meshexp   -istio.test.env  kube \
// -istio.test.hub "gcr.io/istio-release" -istio.test.tag "master-latest-daily" \
// --project_number=<your-gcp-project-number>  --project_id=<your-gcp-project>  \
// --log_output_level=tf:debug,CI:debug  --zone=us-central1-a \
// --deb_url=https://storage.googleapis.com/istio-release/releases/1.1.3/deb
func TestMain(m *testing.M) {
	framework.
		NewSuite("meshexp_test", m).
		// Restrict the test to the K8s environment only, tests will be skipped in native environment.
		RequireEnvironment(environment.Kube).
		// Deploy Istio on the cluster
		Setup(istio.SetupOnKube(nil, setupMeshExpansionInstall)).
		// Create a VM instance before running the test.
		Setup((resource.SetupFn)(setupVMInstance)).
		Run()
}

func setupMeshExpansionInstall(cfg *istio.Config) {
	cfg.Values["global.meshExpansion.enabled"] = "true"
}

// setupVMInstance runs necessary setup on the VM instance and create service
// entry for VM application.
func setupVMInstance(ctx framework.SuiteContext) error {
	var err error
	vmInstance, err = rawvm.New(ctx, rawvm.Config{
		Type: rawvm.GCE,
	})
	if err != nil {
		return fmt.Errorf("failed to create VM service %v", err)
	}
	rawvm.Register(serviceName, ports)
	return nil
}

func TestPilotIsReachable(t *testing.T) {
	output, err := vmInstance.Execute("/bin/sh -c curl localhost:15000/clusters")
	if err != nil {
		t.Errorf("VM instance failed to get Envoy CDS, %v", err)
	}
	// Examine sidecar CDS config to see if control plane exists or not.
	for _, cluster := range []string{"istio-pilot", "istio-citadel"} {
		if !strings.Contains(output, cluster) {
			t.Errorf("%v not found in VM sidecar CDS config", cluster)
		}
	}
}

// TestKubernetesToVM sends a request to a pod in Kubernetes cluster, then the pod sends the request
// to app runs on the VM, returns success if VM app returns success result.
// TODO(incfly): implemets it.
// func TestKubernetesToVM(t *testing.T) {
// }
