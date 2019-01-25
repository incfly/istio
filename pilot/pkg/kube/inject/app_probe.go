// Copyright 2018 Istio Authors
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

package inject

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"istio.io/istio/pilot/cmd/pilot-agent/status"
	"istio.io/istio/pkg/log"
)

const (
	// StatusPortCmdFlagName is the name of the command line flag passed to pilot-agent for sidecar readiness probe.
	// We reuse it for taking over application's readiness probing as well.
	// TODO: replace the hardcoded statusPort elsewhere by this variable as much as possible.
	StatusPortCmdFlagName = "statusPort"

	// TODO: any constant refers to this container's name?
	istioProxyContainerName = "istio-proxy"
)

var (
	// regex pattern for to extract the pilot agent probing port.
	// Supported format, --statusPort, -statusPort, --statusPort=15020.
	statusPortPattern = regexp.MustCompile(fmt.Sprintf(`^-{1,2}%s(=(?P<port>\d+))?$`, StatusPortCmdFlagName))
)

// extractStatusPort accepts the sidecar container spec and returns its port for healthiness probing.
func extractStatusPort(sidecar *corev1.Container) int {
	for i, arg := range sidecar.Args {
		// Skip for unrelated args.
		match := statusPortPattern.FindAllStringSubmatch(strings.TrimSpace(arg), -1)
		if len(match) != 1 {
			continue
		}
		groups := statusPortPattern.SubexpNames()
		portStr := ""
		for ind, s := range match[0] {
			if groups[ind] == "port" {
				portStr = s
				break
			}
		}
		// Port not found from current arg, extract from next arg.
		if portStr == "" {
			// Matches the regex pattern, but without actual values provided.
			if len(sidecar.Args) <= i+1 {
				log.Errorf("No statusPort value provided, skip app probe rewriting")
				return -1
			}
			portStr = sidecar.Args[i+1]
		}
		p, err := strconv.Atoi(portStr)
		if err != nil {
			log.Errorf("Failed to convert statusPort to int %v, err %v", portStr, err)
			return -1
		}
		return p
	}
	return -1
}

func calculateRewrite(probe *corev1.Probe, newURL string, statusPort int) *corev1.HTTPGetAction {
	if probe == nil || probe.HTTPGet == nil {
		return nil
	}
	c := proto.Clone(probe.HTTPGet).(*corev1.HTTPGetAction)
	// Change the application container prober config.
	c.Port = intstr.FromInt(statusPort)
	c.Path = newURL
	return c
}

// createProbeRewritePatch generates the patch for webhook.
// TODO(here): problems, the webhook injector when see a pod, may or maynot have istio-proxy sidecar
// container there.
func createProbeRewritePatch(podSpec *corev1.PodSpec, spec *SidecarInjectionSpec) []rfc6902PatchOperation {
	patch := []rfc6902PatchOperation{}
	if spec == nil || podSpec == nil || !spec.RewriteAppHTTPProbe {
		return patch
	}
	var sidecar *corev1.Container
	for i := range spec.Containers {
		if podSpec.Containers[i].Name == istioProxyContainerName {
			sidecar = &podSpec.Containers[i]
			break
		}
	}
	if sidecar == nil {
		return nil
	}

	statusPort := extractStatusPort(sidecar)
	// Pilot agent statusPort is not defined, skip changing application http probe.
	if statusPort == -1 {
		return nil
	}
	// kubeProbers := &status.KubeAppProbers{}
	for i, c := range podSpec.Containers {
		// Skip sidecar container.
		if c.Name == istioProxyContainerName {
			continue
		}
		portMap := map[string]int32{}
		for _, p := range c.Ports {
			portMap[p.Name] = p.ContainerPort
		}
		readyz, livez := status.FormatProberURL(c.Name)
		if after := calculateRewrite(c.ReadinessProbe, readyz, statusPort); after != nil {
			patch = append(patch, rfc6902PatchOperation{
				Op:    "replace",
				Path:  fmt.Sprintf("/spec/containers/%v/readinessProbe/httpGet", i),
				Value: *after,
			})
		}
		if after := calculateRewrite(c.LivenessProbe, livez, statusPort); after != nil {
			patch = append(patch, rfc6902PatchOperation{
				Op:    "replace",
				Path:  fmt.Sprintf("/spec/containers/%v/livenessProbe/httpGet", i),
				Value: *after,
			})
		}
	}
	return patch
}

// extractKubeAppProbers returns a pointer to the KubeAppProbers.
// Also update the probers so that all the references to the named port will be resolved to integer.
// TODO: reasoning as following
// - we need to update istio-proxy container args.
// - webhook we append that before we call calculate rewrites.
func extractKubeAppProbers(podspec *corev1.PodSpec) *status.KubeAppProbers {
	out := status.KubeAppProbers{}
	updateNamedPort := func(p *corev1.Probe, portMap map[string]int32) *corev1.HTTPGetAction {
		if p == nil || p.HTTPGet == nil {
			return nil
		}
		h := p.HTTPGet
		if h.Port.Type == intstr.String {
			port, exists := portMap[h.Port.StrVal]
			if !exists {
				return nil
			}
			h.Port = intstr.FromInt(int(port))
		}
		return h
	}
	for _, c := range podspec.Containers {
		if c.Name == ProxyContainerName {
			continue
		}
		readyz, livez := status.FormatProberURL(c.Name)
		portMap := map[string]int32{}
		for _, p := range c.Ports {
			if p.Name != "" {
				portMap[p.Name] = p.ContainerPort
			}
		}
		if h := updateNamedPort(c.ReadinessProbe, portMap); h != nil {
			out[readyz] = h
		}
		if h := updateNamedPort(c.LivenessProbe, portMap); h != nil {
			out[livez] = h
		}
	}
	return &out
}

func rewriteAppHTTPProbe(podSpec *corev1.PodSpec, spec *SidecarInjectionSpec) {
	if spec == nil || podSpec == nil {
		return
	}
	if !spec.RewriteAppHTTPProbe {
		return
	}
	var sidecar *corev1.Container
	for i := range podSpec.Containers {
		if podSpec.Containers[i].Name == istioProxyContainerName {
			sidecar = &podSpec.Containers[i]
			break
		}
	}
	if sidecar == nil {
		return
	}

	statusPort := extractStatusPort(sidecar)
	// Pilot agent statusPort is not defined, skip changing application http probe.
	if statusPort == -1 {
		return
	}

	appProberInfo := extractKubeAppProbers(podSpec)
	// Finally propagate app prober config to `istio-proxy` through command line flag.
	b, err := json.Marshal(appProberInfo)
	if err != nil {
		log.Errorf("failed to serialize the app prober config %v", err)
		return
	}
	// We don't have to escape json encoding here when using golang libraries.
	sidecar.Args = append(sidecar.Args, []string{fmt.Sprintf("--%v", status.KubeAppProberCmdFlagName), string(b)}...)

	// Now time to modify the container probers.
	for _, c := range podSpec.Containers {
		// Skip sidecar container.
		if c.Name == istioProxyContainerName {
			continue
		}
		portMap := map[string]int32{}
		for _, p := range c.Ports {
			portMap[p.Name] = p.ContainerPort
		}
		readyz, livez := status.FormatProberURL(c.Name)
		after := calculateRewrite(c.ReadinessProbe, readyz, statusPort)
		if after != nil {
			*c.ReadinessProbe.HTTPGet = *after
		}
		after = calculateRewrite(c.LivenessProbe, livez, statusPort)
		if after != nil {
			*c.LivenessProbe.HTTPGet = *after
		}
	}

}
