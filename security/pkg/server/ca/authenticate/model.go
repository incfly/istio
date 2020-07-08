// Copyright Istio Authors
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

package authenticate

import "context"

// Caller carries the identity and authentication source of a caller.
type Caller struct {
	AuthSource AuthSource
	Identities []string
}

type Authenticator interface {
	Authenticate(ctx context.Context) (*Caller, error)
	AuthenticatorType() string
}
