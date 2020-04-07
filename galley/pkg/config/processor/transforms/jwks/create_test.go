package jwks

import (
	"sort"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	secv1 "istio.io/api/security/v1beta1"
	"istio.io/istio/galley/pkg/config/processing"
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
		"a-jwks-filled": &secv1.RequestAuthentication{
			JwtRules: []*secv1.JWTRule{
				&secv1.JWTRule{
					Issuer: "a-iss",
					Jwks:   "a-pubkey-filled",
				},
			},
		},
		"b": &secv1.RequestAuthentication{
			JwtRules: []*secv1.JWTRule{
				&secv1.JWTRule{
					Issuer:  "b-iss",
					JwksUri: "b-uri",
				},
				&secv1.JWTRule{
					Issuer:  "b-iss-v2",
					JwksUri: "b-uri",
				},
			},
		},
		"a-from-header-foo": &secv1.RequestAuthentication{
			JwtRules: []*secv1.JWTRule{
				&secv1.JWTRule{
					Issuer:  "a-iss",
					JwksUri: "a-uri",
					FromHeaders: []*secv1.JWTHeader{
						&secv1.JWTHeader{
							Name: "x-foo",
						},
					},
				},
			},
		},
	}
)

func policyByName(name string) *secv1.RequestAuthentication {
	p, ok := policyRegistry[name]
	if !ok {
		return nil
	}
	return proto.Clone(p).(*secv1.RequestAuthentication)
}

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
	fh.events = append(fh.events, &e)
}

// ByAge implements sort.Interface based on the Age field.
type byResourceName []*event.Event

func (a byResourceName) Len() int { return len(a) }
func (a byResourceName) Less(i, j int) bool {
	return a[i].Resource.Metadata.FullName.String() < a[j].Resource.Metadata.FullName.String()
}
func (a byResourceName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (fh *fakeHandler) validateEvents(t *testing.T, events []*event.Event) {
	t.Helper()
	// We sort the events because the order of the updates of individual policy does not matter.
	sort.Sort(byResourceName(events))
	sort.Sort(byResourceName(fh.events))
	if diff := cmp.Diff(fh.events, events); diff != "" {
		t.Errorf("handler received different envents, diff %v", diff)
	}
}

func policyAddEvent(t *testing.T, p *secv1.RequestAuthentication) *event.Event {
	if p == nil {
		t.Fatalf("unexpected input nil policy")
	}
	return &event.Event{
		Kind:   event.Added,
		Source: collections.IstioSecurityV1Beta1Requestauthentications,
		Resource: &resource.Instance{
			Message: p,
		},
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
					"a-uri": "a-pubkey-by-resolver",
				},
			},
			updates: jwksUpdate{
				policyEvent: policyAddEvent(t, policyByName("a")),
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
									Jwks:    "a-pubkey-by-resolver",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TransformMultiRules",
			initial: transformState{
				jwksMap: map[string]string{
					"b-uri": "b-pubkey-by-resolver",
				},
			},
			updates: jwksUpdate{
				policyEvent: policyAddEvent(t, policyByName("b")),
			},
			want: []*event.Event{
				&event.Event{
					Kind:   event.Added,
					Source: collections.IstioSecurityV1Beta1Requestauthentications,
					Resource: &resource.Instance{
						Message: &secv1.RequestAuthentication{
							JwtRules: []*secv1.JWTRule{
								&secv1.JWTRule{
									Issuer:  "b-iss",
									JwksUri: "b-uri",
									Jwks:    "b-pubkey-by-resolver",
								},
								&secv1.JWTRule{
									Issuer:  "b-iss-v2",
									JwksUri: "b-uri",
									Jwks:    "b-pubkey-by-resolver",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "OriginalJwksRespectEvenUpdated",
			initial: transformState{
				jwksMap: map[string]string{
					"a-uri": "a-pubkey-by-resolver",
				},
			},
			updates: jwksUpdate{
				policyEvent: policyAddEvent(t, policyByName("a-jwks-filled")),
			},
			want: []*event.Event{
				&event.Event{
					Kind:   event.Added,
					Source: collections.IstioSecurityV1Beta1Requestauthentications,
					Resource: &resource.Instance{
						Message: &secv1.RequestAuthentication{
							JwtRules: []*secv1.JWTRule{
								&secv1.JWTRule{
									Issuer: "a-iss",
									Jwks:   "a-pubkey-filled",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "IgnoreResolverEmptyResponse",
			initial: transformState{
				jwksMap: map[string]string{
					"a-uri": "",
				},
			},
			updates: jwksUpdate{
				policyEvent: policyAddEvent(t, policyByName("a")),
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
									Jwks:    "", // leave it as empty, unmodified.
								},
							},
						},
					},
				},
			},
		},
		{
			name: "BasicTransformByRefresh",
			initial: transformState{
				jwksMap: map[string]string{
					"a-uri": "a-pubkey-by-resolver-v1",
				},
				policies: map[string]*secv1.RequestAuthentication{
					"a":                 policyByName("a"),
					"b":                 policyByName("b"),
					"a-from-header-foo": policyByName("a-from-header-foo"),
				},
			},
			updates: jwksUpdate{
				jwksUpdate: &jwksEntry{
					jwksURI: "a-uri",
					jwks:    "a-pubkey-by-resolver-v2",
				},
			},
			want: []*event.Event{
				&event.Event{
					Kind:   event.Updated,
					Source: collections.IstioSecurityV1Beta1Requestauthentications,
					Resource: &resource.Instance{
						Metadata: resource.Metadata{
							FullName: resource.FullName{
								Name: resource.LocalName("a"),
							},
						},
						Message: &secv1.RequestAuthentication{
							JwtRules: []*secv1.JWTRule{
								&secv1.JWTRule{
									Issuer:  "a-iss",
									JwksUri: "a-uri",
									Jwks:    "a-pubkey-by-resolver-v2",
								},
							},
						},
					},
				},
				&event.Event{
					Kind:   event.Updated,
					Source: collections.IstioSecurityV1Beta1Requestauthentications,
					Resource: &resource.Instance{
						Metadata: resource.Metadata{
							FullName: resource.FullName{
								Name: resource.LocalName("a-from-header-foo"),
							},
						},
						Message: &secv1.RequestAuthentication{
							JwtRules: []*secv1.JWTRule{
								&secv1.JWTRule{
									Issuer:  "a-iss",
									JwksUri: "a-uri",
									Jwks:    "a-pubkey-by-resolver-v2",
									FromHeaders: []*secv1.JWTHeader{
										&secv1.JWTHeader{
											Name: "x-foo",
										},
									},
								},
							},
						},
					},
				},
			},
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
			for k, p := range c.initial.policies {
				xform.Handle(event.Event{
					Kind:   event.Added,
					Source: collections.IstioSecurityV1Beta1Requestauthentications,
					Resource: &resource.Instance{
						Metadata: resource.Metadata{
							FullName: resource.FullName{
								Name: resource.LocalName(k),
							},
						},
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
