//  Copyright 2019 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

// Package basic contains an example test suite for showcase purposes.
package security

import (
	"testing"

	authnv1alpha "istio.io/api/authentication/v1alpha1"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/test/framework/runtime/components/environment/native"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/api/descriptors"
	"istio.io/istio/pkg/test/framework/api/ids"
	"istio.io/istio/pkg/test/framework/api/lifecycle"
)

// TestMtlsHealthCheck verifies Kubernetes HTTP health check can work when mTLS
// is enabled.
func TestMtlsHealthCheck(t *testing.T) {
	ctx := framework.GetContext(t)
	// TODO(incfly): make test able to run both on k8s and native when galley is ready.
	ctx.RequireOrSkip(t, lifecycle.Test, &descriptors.KubernetesEnvironment, &ids.Apps)
	env := native.GetEnvironmentOrFail(ctx, t)
	_, err := env.ServiceManager.ConfigStore.Create(
		model.Config{
			ConfigMeta: model.ConfigMeta{
				Type:      model.AuthenticationPolicy.Type,
				Name:      "default",
				Namespace: "istio-system",
			},
			Spec: &authnv1alpha.Policy{
				// TODO: make policy work just applied to service a.
				// Targets: []*authn.TargetSelector{
				// 	{
				// 		Name: "a.istio-system.svc.local",
				// 	},
				// },
				Peers: []*authnv1alpha.PeerAuthenticationMethod{{
					Params: &authnv1alpha.PeerAuthenticationMethod_Mtls{
						Mtls: &authnv1alpha.MutualTls{
							Mode: authnv1alpha.MutualTls_PERMISSIVE,
						},
					},
				}},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	// apps := components.GetApps(ctx, t)
	// a := apps.GetAppOrFail("a", t)
}
