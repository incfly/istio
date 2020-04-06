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
// Package jwks provides jwt public key transformation for request auethentication.
// - start galley, `gupdate && gglog`
// - port forwarding, kpfn istio-system $(kpidn istio-system -listio=galley) 9901
// - run client:
// go run galley/tools/mcpc/main.go
//   --collections istio/security/v1beta1/requestauthentications  --output long  | grep 'issuer' -A 2`
package jwks

import (
	secv1 "istio.io/api/security/v1beta1"
	"istio.io/istio/galley/pkg/config/processing"
	"istio.io/istio/galley/pkg/config/processing/transformer"
	xformer "istio.io/istio/galley/pkg/config/processing/transformer"
	"istio.io/istio/galley/pkg/config/scope"
	"istio.io/istio/pkg/config/event"
	"istio.io/istio/pkg/config/resource"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
)

func GetProviders() transformer.Providers {
	inputs := collection.NewSchemasBuilder().
		MustAdd(collections.K8SSecurityIstioIoV1Beta1Requestauthentications).
		Build()

	outputs := collection.NewSchemasBuilder().
		MustAdd(collections.IstioSecurityV1Beta1Requestauthentications).
		Build()

	createFn := func(o processing.ProcessorOptions) event.Transformer {
		return &jwksTransformer{
			inputs:   inputs,
			outputs:  outputs,
			options:  o,
			policies: map[string]*resource.Instance{},
			jwksMap: map[string]string{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
				"5": "5",
			},
		}
	}
	return []xformer.Provider{transformer.NewProvider(inputs, outputs, createFn)}
}

// Start implements event.Transformer
func (t *jwksTransformer) Start() {
}

// Stop implements event.Transformer
func (t *jwksTransformer) Stop() {
}

// Inputs implements event.Transformer
func (t *jwksTransformer) Inputs() collection.Schemas {
	return t.inputs
}

// Outputs implements event.Transformer
func (t *jwksTransformer) Outputs() collection.Schemas {
	return t.outputs
}

func (t *jwksTransformer) updateJwks(policy *secv1.RequestAuthentication) bool {
	updated := false
	for _, r := range policy.GetJwtRules() {
		iss := r.GetIssuer()
		jwks, ok := t.jwksMap[iss]
		if !ok {
			continue
		}
		r.Jwks = jwks
		updated = true
	}
	return updated
}

func (t *jwksTransformer) updateJwksMap(policy *secv1.RequestAuthentication) {
	for _, rule := range policy.GetJwtRules() {
		iss := rule.GetIssuer()
		jwks := rule.GetJwks()
		t.jwksMap[iss] = jwks
	}
}

// TODO: only the affecting working set is pushed, not persistent?
// next working set, others fall back to original.
// Handle implements event.Transformer
func (t *jwksTransformer) Handle(e event.Event) {
	if e.Resource != nil &&
		e.Resource.Metadata.FullName.Namespace == "asm-jwks-internal-event" {
		t.updateJwksMap(e.Resource.Message.(*secv1.RequestAuthentication))
		// iterate all the policies.
		for k, p := range t.policies {
			msg := p.Message.(*secv1.RequestAuthentication)
			updated := t.updateJwks(msg)
			if updated {
				scope.Processing.Infof("incfly/transform, perform update %v, policy %v", k, msg)
				t.dispatch(event.Event{
					Kind:     event.Updated,
					Resource: p,
					Source:   collections.IstioSecurityV1Beta1Requestauthentications,
				})
			}
		}
		return
	}

	switch e.Kind {
	case event.Added, event.Updated:
		updated := t.updateJwks(e.Resource.Message.(*secv1.RequestAuthentication))
		scope.Processing.Infof("incfly/init add, updated %v", updated)
		t.policies[e.Resource.Metadata.FullName.String()] = e.Resource
	case event.Deleted:
		delete(t.policies, e.Resource.Metadata.FullName.String())
	case event.Reset:
		t.policies = map[string]*resource.Instance{}
	default:
	}

	t.dispatch(e)
}

func (t *jwksTransformer) dispatch(e event.Event) {
	if t.handler != nil {
		t.handler.Handle(e)
	}
}

type jwksTransformer struct {
	inputs   collection.Schemas
	outputs  collection.Schemas
	options  processing.ProcessorOptions
	handler  event.Handler
	policies map[string]*resource.Instance
	jwksMap  map[string]string
}

// DispatchFor implements event.Transformer
func (t *jwksTransformer) DispatchFor(c collection.Schema, h event.Handler) {
	// scope.Processing.Infof("incfly/jwks, DispatchFor handler %v, c %v", h, c)
	switch c.Name() {
	case collections.IstioSecurityV1Beta1Requestauthentications.Name():
		t.handler = event.CombineHandlers(t.handler, h)
	}
}
