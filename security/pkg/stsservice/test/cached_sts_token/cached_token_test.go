// Copyright 2020 Istio Authors
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

package cachedststoken

import (
	"testing"

	"github.com/onsi/gomega"

	testID "istio.io/istio/mixer/test/client/env"
	xdsService "istio.io/istio/security/pkg/stsservice/mock"
	stsTest "istio.io/istio/security/pkg/stsservice/test"
)

// TestCachedToken verifies when proxy reconnects XDS server and sends token over
// the stream, if the original token is not expired, gRPC library does not call
// STS server and provides cached token to proxy.
func TestCachedToken(t *testing.T) {
	// Enable this test when gRPC fix is picked by Istio Proxy
	// https://github.com/grpc/grpc/pull/21641
	t.Skip("https://github.com/istio/istio/issues/20133")
	// Sets up callback that verifies token on new XDS stream.
	cb := xdsService.CreateXdsCallback(t)
	numCloseStream := 3
	// Force XDS server to close streams 3 times and keep the 4th stream open.
	cb.SetNumberOfStreamClose(numCloseStream, 0)
	// Start all test servers and proxy
	setup := stsTest.SetUpTest(t, cb, testID.STSCacheTest)
	// Explicitly set token life time to a long duration.
	setup.AuthServer.SetTokenLifeTime(3600)
	// Explicitly set auth server to return different access token to each call.
	setup.AuthServer.EnableDynamicAccessToken(true)
	// Verify that initially XDS stream is not set up, stats are not incremented.
	g := gomega.NewWithT(t)
	g.Expect(cb.NumStream()).To(gomega.Equal(0))
	g.Expect(cb.NumTokenReceived()).To(gomega.Equal(0))
	// Get initial number of calls to auth server. They are not zero due to STS flow test
	// in the test setup phase, which is to make sure the servers are up and ready.
	initialNumFederatedTokenCall := setup.AuthServer.NumGetFederatedTokenCalls()
	initialNumAccessTokenCall := setup.AuthServer.NumGetAccessTokenCalls()
	setup.StartProxy(t)
	setup.ProxySetUp.WaitEnvoyReady()
	// Verify that proxy re-connects XDS server after each stream close, and the
	// same token is received.
	g.Expect(cb.NumStream()).To(gomega.Equal(numCloseStream + 1))
	g.Expect(cb.NumTokenReceived()).To(gomega.Equal(1))
	// Verify there is only one extra call for each token.
	g.Expect(setup.AuthServer.NumGetFederatedTokenCalls()).To(gomega.Equal(initialNumFederatedTokenCall + 1))
	g.Expect(setup.AuthServer.NumGetAccessTokenCalls()).To(gomega.Equal(initialNumAccessTokenCall + 1))
	setup.TearDown()
}
