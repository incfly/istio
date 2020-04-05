package collections

import (
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/resource"
)

var (
	JwksUpdateResourceSchema = collection.Builder{
		Name:         "k8s/security.istio.io/v1beta1/requestauthenticationsupdater",
		VariableName: "K8SSecurityIstioIoV1Beta1Requestauthenticationsupdater",
		Disabled:     false,
		Resource: resource.Builder{
			Group:        "security.istio.io",
			Kind:         "RequestAuthenticationUpdater",
			Plural:       "requestauthenticationsUpdater",
			Version:      "v1beta1",
			Proto:        "istio.security.v1beta1.RequestAuthentication", // TODO: some other proto?
			ProtoPackage: "istio.io/api/security/v1beta1",
			// ClusterScoped: false,
			// ValidateProto: validation.ValidateRequestAuthentication,
		}.MustBuild(),
	}.MustBuild()
)
