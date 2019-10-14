#!/usr/bin/env bash

# Copyright Istio Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euxo pipefail

# Extract Istio add-on version.
export TAG=$(grep -Pom1 'image:.*:\K\d{1,3}\.\d{1,3}.\d{1,3}' "$MANIFEST_FILE")

if [ -z "$TAG" ]; then
  echo "Unable to parse Istio version tag" >&2
  exit 1
fi

# Clone Istio at this tag.
git -C "$GOPATH/src" clone --branch "$TAG" "https://github.com/istio/istio.git" "istio.io/istio"

# Run smoke test.
JUNIT_E2E_XML="${ARTIFACTS}/junit.xml" \
  TARGET=e2e_bookinfo_envoyv2_v1alpha3 \
  E2E_ARGS='--skip_setup --namespace=istio-system --test.run=Test[^DbRoutingMysql]' \
  make -C "$GOPATH/src/istio.io/istio" with_junit_report

# Test add-on disable.
gcloud beta container clusters update istio-e2e --project="$PROJECT" --zone="$ZONE" --update-addons=Istio=DISABLED
