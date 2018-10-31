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
	"time"
	"fmt"
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
func TestAutoMtlsCDS(t *testing.T) {
	initLocalPilotTestEnv(t)
	adsc := adsConnectAndWait(t, 0x0a0a0a0a)
	server := util.MockTestServer
	defer adsc.Close()

	endpoints := []*model.IstioEndpoint{
		newEndpointWithAccount("127.0.0.1", "sa1", "v1"),
		newEndpointWithAccount("127.0.0.2", "sa1", "v1"),
	}

	server.EnvoyXdsServer.MemRegistry.SetEndpoints("cds.test.svc.cluster.local", endpoints)
	adsc.WaitClear()
	adsc.Wait("", time.Second*5)
	fmt.Println("jianfeih debug cluster is ", adsc.Clusters)
	//cdsr, err := connectADS(util.MockPilotGrpcAddr)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//err = sendCDSReq(sidecarId(app3Ip, "app3"), cdsr)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//res, err := cdsr.Recv()
	//if err != nil {
	//	t.Fatal("Failed to receive CDS", err)
	//	return
	//}
	//
	//strResponse, _ := model.ToJSONWithIndent(res, " ")
	//_ = ioutil.WriteFile(env.IstioOut+"/cdsv2_sidecar.json", []byte(strResponse), 0644)
	//
	//t.Log("CDS response", strResponse)
	//if len(res.Resources) == 0 {
	//	t.Fatal("No response")
	//}
}