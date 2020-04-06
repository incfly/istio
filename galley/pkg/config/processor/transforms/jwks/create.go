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
	"istio.io/istio/galley/pkg/config/processing"
	"istio.io/istio/galley/pkg/config/processing/transformer"
	xformer "istio.io/istio/galley/pkg/config/processing/transformer"
	"istio.io/istio/galley/pkg/config/scope"
	"istio.io/istio/pkg/config/event"
	"istio.io/istio/pkg/config/resource"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
)

var (
	count = 1
)

// GetProviders returns transformer providers for auth policy transformers
// func GetProviders() transformer.Providers {
// 	return []transformer.Provider{
// 		transformer.NewSimpleTransformerProvider(
// 			collections.K8SSecurityIstioIoV1Beta1Requestauthentications,          // k8s version schema
// 			collections.IstioSecurityV1Beta1Requestauthentications,               // istio version schema.
// 			handler(collections.K8SSecurityIstioIoV1Beta1Requestauthentications), // GALLEY-NOTE: werid... type, what's from to for?
// 		),
// 	}
// }

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
	scope.Processing.Infof("incfly jwks transformer started")
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
	scope.Processing.Infof("incfly/after updating, jwksmap %v", t.jwksMap)
}

// TODO: only the affecting working set is pushed, not persistent?
// next working set, others fall back to original.
// Handle implements event.Transformer
func (t *jwksTransformer) Handle(e event.Event) {
	scope.Processing.Infof("incfly transfomr handle event invoked %v", e)
	if e.Resource != nil &&
		e.Resource.Metadata.FullName.Namespace == "asm-jwks-internal-event" {
		t.updateJwksMap(e.Resource.Message.(*secv1.RequestAuthentication))
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

	// update internal cache of the req authn map
	switch e.Kind {
	case event.Added, event.Updated:
		// TODO(here): DO transform here!
		updated := t.updateJwks(e.Resource.Message.(*secv1.RequestAuthentication))
		scope.Processing.Infof("incfly/init add, updated %v", updated)
		t.policies[e.Resource.Metadata.FullName.String()] = e.Resource
	case event.Deleted:
		delete(t.policies, e.Resource.Metadata.FullName.String())
	case event.Reset:
		t.policies = map[string]*resource.Instance{}
	default:
	}

	// always doing nothing.
	t.dispatch(e)
	// switch e.Kind {
	// case event.FullSync:
	// 	t.dispatch(event.FullSyncFor(collections.IstioNetworkingV1Alpha3SyntheticServiceentries))
	// 	return

	// case event.Reset:
	// 	t.dispatch(event.Event{Kind: event.Reset})
	// 	return

	// case event.Added, event.Updated, event.Deleted:
	// 	// fallthrough

	// default:
	// 	panic(fmt.Errorf("transformer.Handle: Unexpected event received: %v", e))
	// }

	// switch e.Source.Name() {
	// default:
	// 	// panic(fmt.Errorf("received event with unexpected collection: %v", e.Source.Name()))
	// }
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
	scope.Processing.Infof("incfly/jwks, DispatchFor handler %v, c %v", h, c)
	switch c.Name() {
	case collections.IstioSecurityV1Beta1Requestauthentications.Name():
		t.handler = event.CombineHandlers(t.handler, h)
	}
}

func handler(destination collection.Schema) func(e event.Event, h event.Handler) {
	return func(e event.Event, h event.Handler) {
		scope.Processing.Infof("incfly/jwks/transformer invoked, event %v", e)
		if e.Resource.Metadata.FullName.Namespace == "asm-jwks-internal-event" {
			scope.Processing.Infof("incfly/jwks internal asm jws event")
			// TODO: workaround using ADD && DELETE AGAIN does not work.
			// ne := e.Clone()
			// ne.Kind = event.Deleted
			// scope.Processing.Infof("incfly/jwks internal asm jws event, deleting it, before %v, after %v", e, ne)
			// h.Handle(ne)
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
