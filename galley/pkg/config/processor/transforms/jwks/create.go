// Copyright 2019 Istio Authors
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

package jwks

import (
	"fmt"
	"reflect"

	secv1 "istio.io/api/security/v1beta1"
	"istio.io/istio/galley/pkg/config/processing/transformer"
	"istio.io/istio/galley/pkg/config/scope"
	"istio.io/istio/pkg/config/event"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
)

var (
	count = 1
)

// GetProviders returns transformer providers for auth policy transformers
func GetProviders() transformer.Providers {
	scope.Processing.Infof("incfly/create.go invoked")
	return []transformer.Provider{
		transformer.NewSimpleTransformerProvider(
			collections.K8SSecurityIstioIoV1Beta1Requestauthentications,          // k8s version schema
			collections.IstioSecurityV1Beta1Requestauthentications,               // istio version schema.
			handler(collections.K8SSecurityIstioIoV1Beta1Requestauthentications), // GALLEY-NOTE: werid... type, what's from to for?
		),
	}
}

func handler(destination collection.Schema) func(e event.Event, h event.Handler) {
	return func(e event.Event, h event.Handler) {
		scope.Processing.Infof("incfly/jwks/transformer invoked, event %v", e)
		if e.Resource.Metadata.FullName.Namespace == "asm-jwks-internal-event" {
			scope.Processing.Infof("incfly/jwks internal asm jws event")
			return
		}
		e = e.WithSource(destination)
		if e.Resource != nil && e.Resource.Message != nil {
			policy, ok := e.Resource.Message.(*secv1.RequestAuthentication)
			if !ok {
				scope.Processing.Errorf("unexpected proto found when converting authn.Policy: %v", reflect.TypeOf(e.Resource.Message))
				return
			}
			policy.GetJwtRules()[0].Jwks = fmt.Sprintf("jwtver-%v", count)
			count++
		}

		h.Handle(e)
	}
}
