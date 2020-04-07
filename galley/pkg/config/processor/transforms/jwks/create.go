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
		return newJwksTransformer(&fakeJwksresolver{}, o)
	}
	return []xformer.Provider{transformer.NewProvider(inputs, outputs, createFn)}
}

type jwksTransformer struct {
	inputs   collection.Schemas
	outputs  collection.Schemas
	options  processing.ProcessorOptions
	handler  event.Handler
	policies map[string]*resource.Instance
	// jwksMap  map[string]string
	resolver JwksResolverHelper
}
type JwksUpdateHandler func() error

// JwksResolverHelper is a wrapper interface around actual jwks resolving implementation, intended used
// by Galley transformer.
type JwksResolverHelper interface {
	SetUpdateFunc(fn JwksUpdateHandler)
	ResolveJwks(jwksURI string) string
}

func newJwksTransformer(resolver JwksResolverHelper, opt processing.ProcessorOptions) *jwksTransformer {
	inputs := collection.NewSchemasBuilder().
		MustAdd(collections.K8SSecurityIstioIoV1Beta1Requestauthentications).
		Build()

	outputs := collection.NewSchemasBuilder().
		MustAdd(collections.IstioSecurityV1Beta1Requestauthentications).
		Build()

	xform := &jwksTransformer{
		inputs:   inputs,
		outputs:  outputs,
		options:  opt,
		policies: map[string]*resource.Instance{},
		resolver: resolver,
	}
	resolver.SetUpdateFunc(xform.jwksUpdateHandler)
	return xform
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

func (t *jwksTransformer) overridePolicy(policy *secv1.RequestAuthentication) bool {
	updated := false
	for _, r := range policy.GetJwtRules() {
		uri := r.GetJwksUri()
		jwks := t.resolver.ResolveJwks(uri)
		if jwks != "" {
			r.Jwks = jwks
			updated = true
		} else {
			scope.Processing.Errorf("jwks transforming resolver failed, empty jwks for jwks %v", uri)
		}
	}
	return updated
}

// Handle implements event.Transformer.
func (t *jwksTransformer) Handle(e event.Event) {
	switch e.Kind {
	case event.Added, event.Updated:
		updated := t.overridePolicy(e.Resource.Message.(*secv1.RequestAuthentication))
		scope.Processing.Debugf("incfly/init add, updated %v", updated)
		t.policies[e.Resource.Metadata.FullName.String()] = e.Resource
	case event.Deleted:
		delete(t.policies, e.Resource.Metadata.FullName.String())
	case event.Reset:
		t.policies = map[string]*resource.Instance{}
	default:
	}
	t.dispatch(e)
}

// jwksUpdateHandler should be invoked by resolver when a jwks is updated.
func (t *jwksTransformer) jwksUpdateHandler() error {
	// Iterate all the policies.
	for _, p := range t.policies {
		msg := p.Message.(*secv1.RequestAuthentication)
		updated := t.overridePolicy(msg)
		if updated {
			t.dispatch(event.Event{
				Kind:     event.Updated,
				Resource: p,
				Source:   collections.IstioSecurityV1Beta1Requestauthentications,
			})
		}
	}
	return nil
}

func (t *jwksTransformer) dispatch(e event.Event) {
	if t.handler != nil {
		t.handler.Handle(e)
	}
}

// DispatchFor implements event.Transformer
func (t *jwksTransformer) DispatchFor(c collection.Schema, h event.Handler) {
	switch c.Name() {
	case collections.IstioSecurityV1Beta1Requestauthentications.Name():
		t.handler = event.CombineHandlers(t.handler, h)
	}
}
