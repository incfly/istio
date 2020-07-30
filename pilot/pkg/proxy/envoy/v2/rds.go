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

package v2

import (
	"time"

	"istio.io/istio/pkg/util/protomarshal"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"

	"istio.io/istio/pilot/pkg/gcpmonitoring"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/util"
)

func (s *DiscoveryServer) pushRoute(con *XdsConnection, push *model.PushContext, version string) error {
	pushStart := time.Now()
	rawRoutes := s.ConfigGenerator.BuildHTTPRoutes(con.node, push, con.Routes)
	if s.DebugConfigs {
		for _, r := range rawRoutes {
			con.RouteConfigs[r.Name] = r
			if adsLog.DebugEnabled() {
				resp, _ := protomarshal.ToJSONWithIndent(r, " ")
				adsLog.Debugf("RDS: Adding route:%s for node:%v", resp, con.node.ID)
			}
		}
	}

	response := routeDiscoveryResponse(rawRoutes, version, push.Version, con.RequestedTypes.RDS)
	err := con.send(response)
	rdsPushTime.Record(time.Since(pushStart).Seconds())
	if err != nil {
		adsLog.Warnf("RDS: Send failure for node:%v: %v", con.node.ID, err)
		recordSendError("RDS", rdsSendErrPushes, err)
		return err
	}
	rdsPushes.Increment()
	gcpmonitoring.IncrementConfigPushMeasuare("RDS", true)

	adsLog.Infof("RDS: PUSH for node:%s routes:%d", con.node.ID, len(rawRoutes))
	return nil
}

func routeDiscoveryResponse(rs []*xdsapi.RouteConfiguration, version, noncePrefix, typeURL string) *xdsapi.DiscoveryResponse {
	resp := &xdsapi.DiscoveryResponse{
		TypeUrl:     typeURL,
		VersionInfo: version,
		Nonce:       nonce(noncePrefix),
	}
	for _, rc := range rs {
		rr := util.MessageToAny(rc)
		rr.TypeUrl = typeURL
		resp.Resources = append(resp.Resources, rr)
	}

	return resp
}
