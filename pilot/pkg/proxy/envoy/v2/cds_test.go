// Copyright 2018 Istio Authors
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
package v2_test

import (
	"io/ioutil"
	"testing"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/test/env"
	"istio.io/istio/tests/util"
	"fmt"
	"time"
	"istio.io/istio/pkg/adsc"
)

func TestCDS(t *testing.T) {
	initLocalPilotTestEnv(t)

	cdsr, err := connectADS(util.MockPilotGrpcAddr)
	if err != nil {
		t.Fatal(err)
	}

	err = sendCDSReq(sidecarId(app3Ip, "app3"), cdsr)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cdsr.Recv()
	if err != nil {
		t.Fatal("Failed to receive CDS", err)
		return
	}

	strResponse, _ := model.ToJSONWithIndent(res, " ")
	_ = ioutil.WriteFile(env.IstioOut+"/cdsv2_sidecar.json", []byte(strResponse), 0644)

	t.Log("CDS response", strResponse)
	if len(res.Resources) == 0 {
		t.Fatal("No response")
	}

	// TODO: dump the response resources, compare with some golden once it's stable
	// check that each mocked service and destination rule has a corresponding resource

	// TODO: dynamic checks ( see EDS )
}


// TestAutoMtlsCDS tests the auto mtls feature. If a service consists of an endpoints all have
// mtls_ready label, we configure the Cluster's TLS settings to be tls.
// TODO: TestEnvoy fails because local envoy unable to load `/etc/certs/` path. Solution can be make that configurable
// discovery request thus to do dependency injection.
func TestAutoMtlsCDS(t *testing.T) {
	initLocalPilotTestEnv(t)
	server := util.MockTestServer

	endpoints := []*model.IstioEndpoint{
		newEndpointWithAccount("127.0.0.1", "sa1", "v1"),
		newEndpointWithAccount("127.0.0.2", "sa1", "v1"),
	}
	for _, ep := range endpoints {
		ep.Labels["authentication.istio.io/mtls_ready"] = "true"
	}
	svcName := "cds.test.svc.cluster.local"
	server.EnvoyXdsServer.MemRegistry.AddHTTPService(svcName, "10.0.0.1", 8000)
	server.EnvoyXdsServer.MemRegistry.SetEndpoints(svcName, endpoints)

	adsc, err := adsc.Dial(util.MockPilotGrpcAddr, "", &adsc.Config{
		IP: testIp(uint32(0x0a0a0a0a)),
	})
	if err != nil {
		t.Fatal("Error connecting ", err)
	}
	defer adsc.Close()

	tlsChecker := func() {
		adsc.Wait("cds", time.Second*5)
		for name, cluster := range adsc.Clusters {
			if  name == "outbound|8000||cds.test.svc.cluster.local" {
				fmt.Println("cluster name ", name, "\nCluster", cluster)
			}
		}
	}

	fmt.Println("jianfeih debug, first endpoints setup done")
	adsc.Watch()
	tlsChecker()

	// Now adds an IstioEndpoint with annotation, should still see TLS settings.
	//epNew := newEndpointWithAccount("127.0.0.3", "sa1", "v1")
	//epNew.Labels["authentication.istio.io/mtls_ready"] = "true"
	//endpoints = append(endpoints, epNew)
	//server.EnvoyXdsServer.MemRegistry.SetEndpoints(svcName, endpoints)
	//
	//adsc.Wait("cds", time.Second*5)
	//tlsChecker()

	fmt.Println("jianfeih debug added un annodated endpoint")
	// Add an endpoint without annotation, expect to see cluster without TLS settings.
	epNotReady := newEndpointWithAccount("127.0.0.4", "sa1", "v1")
	endpoints = append(endpoints, epNotReady)
	server.EnvoyXdsServer.MemRegistry.SetEndpoints(svcName, endpoints)

	tlsChecker()
}