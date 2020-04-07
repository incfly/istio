package jwks

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	secv1 "istio.io/api/security/v1beta1"
	"istio.io/istio/galley/pkg/config/processing"
	"istio.io/istio/galley/pkg/config/scope"
	"istio.io/istio/pkg/config/event"
	"istio.io/istio/pkg/config/resource"
	"istio.io/istio/pkg/config/schema/collections"
)

var (
	// Issuer name capitalized means jwks already confiugred without requiring conversion, lower case
	// means requiring.
	policyRegistry = map[string]*secv1.RequestAuthentication{
		"a": &secv1.RequestAuthentication{
			JwtRules: []*secv1.JWTRule{
				&secv1.JWTRule{
					Issuer:  "a-iss",
					JwksUri: "a-uri",
				},
			},
		},
		"A": &secv1.RequestAuthentication{
			JwtRules: []*secv1.JWTRule{
				&secv1.JWTRule{
					Issuer: "a-iss",
					Jwks:   "a-pubkey",
				},
			},
		},
	}
)

type fakeJwksresolver struct {
	jwksMap  map[string]string
	updateFn JwksUpdateHandler
}

func (r *fakeJwksresolver) ResolveJwks(jwksURI string) string {
	return r.jwksMap[jwksURI]
}

func (r *fakeJwksresolver) SetUpdateFunc(fn JwksUpdateHandler) {
	r.updateFn = fn
}

func (r *fakeJwksresolver) update(jwksUri, jwks string) {
	val, ok := r.jwksMap[jwksUri]
	if !ok || val != jwks {
		r.jwksMap[jwksUri] = jwks
		if r.updateFn != nil {
			r.updateFn()
		}
	}
}

type transformState struct {
	jwksMap  map[string]string
	policies map[string]*secv1.RequestAuthentication
}

type jwksEntry struct {
	jwksURI string
	jwks    string
}

type jwksUpdate struct {
	policyEvent *event.Event
	jwksUpdate  *jwksEntry
}

type fakeHandler struct {
	events []*event.Event
}

func (fh *fakeHandler) Handle(e event.Event) {
	scope.Processing.Infof("incfly debug fake handler event %v", e)
	fh.events = append(fh.events, &e)
}

func (fh *fakeHandler) validateEvents(t *testing.T, events []*event.Event) {
	t.Helper()
	if diff := cmp.Diff(fh.events, events); diff != "" {
		t.Errorf("handler received different envents, diff %v", diff)
	}
}

func TestJwksTransformer(t *testing.T) {
	testCases := []struct {
		name string
		// initial state of the transformer.
		initial transformState
		// updates is the changes we applied sequentially.
		updates jwksUpdate
		// want is expected events passed by the transformer.
		want []*event.Event
	}{
		{
			name: "BasicTransform",
			initial: transformState{
				jwksMap: map[string]string{
					"a-uri": "a-pubkey",
				},
			},
			updates: jwksUpdate{
				// Add single policy for "a".
				policyEvent: &event.Event{
					Kind:   event.Added,
					Source: collections.IstioSecurityV1Beta1Requestauthentications,
					Resource: &resource.Instance{
						Message: policyRegistry["a"],
					},
				},
			},
			want: []*event.Event{
				&event.Event{
					Kind:   event.Added,
					Source: collections.IstioSecurityV1Beta1Requestauthentications,
					Resource: &resource.Instance{
						Message: &secv1.RequestAuthentication{
							JwtRules: []*secv1.JWTRule{
								&secv1.JWTRule{
									Issuer:  "a-iss",
									JwksUri: "a-uri",
									Jwks:    "a-pubkey",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "IgnoreResolverEmptyResponse",
			// want: []*event.Event{},
		},
		{
			name: "OriginalJwksRespectEvenUpdated",
		},
		{
			name: "PolicyAddedTransform",
		},
		{
			name: "JwksUpdate",
		},
		{
			name: "JwksUpdateOnlyAffected",
		},
		{
			name: "PolicyDeletionAndJwksUpdate",
		},
	}
	for _, tc := range testCases {
		// fill in source schema of event.
		t.Run(tc.name, func(t *testing.T) {
			c := tc
			fh := &fakeHandler{}
			res := &fakeJwksresolver{
				jwksMap: c.initial.jwksMap,
			}
			// init the state.
			xform := newJwksTransformer(res, processing.ProcessorOptions{})
			for _, p := range c.initial.policies {
				xform.Handle(event.Event{
					Kind:   event.Added,
					Source: collections.IstioSecurityV1Beta1Requestauthentications,
					Resource: &resource.Instance{
						Message: p,
					},
				})
			}
			xform.DispatchFor(collections.IstioSecurityV1Beta1Requestauthentications, fh)
			// Apply operation.
			if c.updates.policyEvent != nil {
				xform.Handle(c.updates.policyEvent.Clone())
			}
			if jwt := c.updates.jwksUpdate; jwt != nil {
				res.update(jwt.jwksURI, jwt.jwks)
			}
			// Check the events output.
			fh.validateEvents(t, c.want)
		})
	}
}
