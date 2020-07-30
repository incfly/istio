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

package v2

import (
	"time"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"

	"istio.io/istio/pilot/pkg/gcpmonitoring"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/util"
)

func (s *DiscoveryServer) pushLds(con *XdsConnection, push *model.PushContext, version string) error {
	// TODO: Modify interface to take services, and config instead of making library query registry
	pushStart := time.Now()
	rawListeners := s.ConfigGenerator.BuildListeners(con.node, push)

	if s.DebugConfigs {
		con.LDSListeners = rawListeners
	}
	response := ldsDiscoveryResponse(rawListeners, version, push.Version, con.RequestedTypes.LDS)
	err := con.send(response)
	ldsPushTime.Record(time.Since(pushStart).Seconds())
	if err != nil {
		adsLog.Warnf("LDS: Send failure %s: %v", con.ConID, err)
		recordSendError("LDS", ldsSendErrPushes, err)
		return err
	}
	ldsPushes.Increment()
	gcpmonitoring.IncrementConfigPushMeasuare("LDS", true)

	adsLog.Infof("LDS: PUSH for node:%s listeners:%d", con.node.ID, len(rawListeners))
	return nil
}

// LdsDiscoveryResponse returns a list of listeners for the given environment and source node.
func ldsDiscoveryResponse(ls []*xdsapi.Listener, version, noncePrefix, typeURL string) *xdsapi.DiscoveryResponse {
	resp := &xdsapi.DiscoveryResponse{
		TypeUrl:     typeURL,
		VersionInfo: version,
		Nonce:       nonce(noncePrefix),
	}
	for _, ll := range ls {
		if ll == nil {
			adsLog.Errora("Nil listener ", ll)
			totalXDSInternalErrors.Increment()
			continue
		}
		lr := util.MessageToAny(ll)
		lr.TypeUrl = typeURL
		resp.Resources = append(resp.Resources, lr)
	}

	return resp
}
