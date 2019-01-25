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

// calculateRewrite returns a pointer to the HTTPGetAction.
func calculateRewrite(probe *corev1.Probe, newURL string,
	statusPort int, portMap map[string]int32, kubeProbers *status.KubeAppProbers) *corev1.HTTPGetAction {
	if probe == nil || probe.HTTPGet == nil {
		return nil
	}
	copyProber := func() *corev1.HTTPGetAction {
		c := proto.Clone(probe.HTTPGet).(*corev1.HTTPGetAction)
		c.Port = probe.HTTPGet.Port
		return c
	}
	original := copyProber()
	after := copyProber()
	// A named port, resolve by looking at port map.
	if original.Port.Type == intstr.String {
		port, exists := portMap[original.Port.StrVal]
		if !exists {
			log.Errorf("named port not found in the map skip rewriting probing %v", *probe)
			return after
		}
		original.Port = intstr.FromInt(int(port))
	}
	(*kubeProbers)[newURL] = original
	// Change the application container prober config.
	after.Port = intstr.FromInt(statusPort)
	after.Path = newURL
	return after
}

func calculateRewrite2(probe *corev1.Probe, newURL string, statusPort int, portMap map[string]int32) *corev1.HTTPGetAction {
	if probe == nil || probe.HTTPGet == nil {
		return nil
	}
	copyProber := func() *corev1.HTTPGetAction {
		c := proto.Clone(probe.HTTPGet).(*corev1.HTTPGetAction)
		c.Port = probe.HTTPGet.Port
		return c
	}
	original := copyProber()
	after := copyProber()
	// A named port, resolve by looking at port map.
	if original.Port.Type == intstr.String {
		port, exists := portMap[original.Port.StrVal]
		if !exists {
			log.Errorf("named port not found in the map skip rewriting probing %v", *probe)
			return after
		}
		original.Port = intstr.FromInt(int(port))
	}
	// Change the application container prober config.
	after.Port = intstr.FromInt(statusPort)
	after.Path = newURL
	return after
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
	kubeProbers := &status.KubeAppProbers{}
	for i, c := range podSpec.Containers {
		// Skip sidecar container.
		if c.Name == istioProxyContainerName {
			continue
		}
		portMap := map[string]int32{}
		for _, p := range c.Ports {
			portMap[p.Name] = p.ContainerPort
		}
		if after := calculateRewrite(c.ReadinessProbe, fmt.Sprintf("/app-health/%v/readyz", c.Name),
			statusPort, portMap, kubeProbers); after != nil {
			patch = append(patch, rfc6902PatchOperation{
				Op:    "replace",
				Path:  fmt.Sprintf("/spec/containers/%v/readinessProbe/httpGet", i),
				Value: *after,
			})
		}
		if after := calculateRewrite(c.LivenessProbe, fmt.Sprintf("/app-health/%v/livez", c.Name),
			statusPort, portMap, kubeProbers); after != nil {
			patch = append(patch, rfc6902PatchOperation{
				Op:    "replace",
				Path:  fmt.Sprintf("/spec/containers/%v/livenessProbe/httpGet", i),
				Value: *after,
			})
		}
	}

	// TODO: add a patch to add a args flag, [--kubeAppHTTPProbers=xxx] kubeProbers.
	// patch = append(patch, rfc6902PatchOperation{
	// 	Op:    "append",
	// 	Path:  fmt.Sprintf("/spec/containers/%v/args/-"),
	// 	Value: *kubeProbers,
	// })
	return patch
}

// extractKubeAppProbers returns a pointer to the KubeAppProbers.
// Also update the probers so that all the references to the named port will be resolved to integer.
// TODO: reasoning as following
// - we need to update istio-proxy container args.
// - webhook we append that before we call calculate rewrites.
func extractKubeAppProbers(podspec *corev1.PodSpec) *status.KubeAppProbers {
	probers := status.KubeAppProbers{}
	for _, c := range podspec.Containers {
		if c.ReadinessProbe == nil || c.ReadinessProbe.HTTPGet == nil {
			continue
		}
		// portMap := map[string]int{}
		ha := proto.Clone(c.ReadinessProbe.HTTPGet).(*corev1.HTTPGetAction)
		// TODO here.
		if c.ReadinessProbe.HTTPGet.Port.Type == intstr.String {
		}
		ha.Port = c.ReadinessProbe.HTTPGet.Port
		probers[fmt.Sprintf("/app-health/%v/readyz", c.Name)] = ha
	}
	return &probers
}

func rewriteAppHTTPProbe2(podSpec *corev1.PodSpec, spec *SidecarInjectionSpec) {
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

	appProberInfo := status.KubeAppProbers{}
	for _, c := range podSpec.Containers {
		// Skip sidecar container.
		if c.Name == istioProxyContainerName {
			continue
		}
		portMap := map[string]int32{}
		for _, p := range c.Ports {
			portMap[p.Name] = p.ContainerPort
		}
		newURL := fmt.Sprintf("/app-health/%v/readyz", c.Name)
		after := calculateRewrite(c.ReadinessProbe, newURL, statusPort, portMap, &appProberInfo)
		if after != nil {
			*c.ReadinessProbe.HTTPGet = *after
		}
		newURL = fmt.Sprintf("/app-health/%v/livez", c.Name)
		after = calculateRewrite(c.LivenessProbe, newURL, statusPort, portMap, &appProberInfo)
		if after != nil {
			*c.LivenessProbe.HTTPGet = *after
		}
	}

	// Finally propagate app prober config to `istio-proxy` through command line flag.
	b, err := json.Marshal(appProberInfo)
	if err != nil {
		log.Errorf("failed to serialize the app prober config %v", err)
		return
	}
	// We don't have to escape json encoding here when using golang libraries.
	sidecar.Args = append(sidecar.Args, []string{fmt.Sprintf("--%v", status.KubeAppProberCmdFlagName), string(b)}...)
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
	// Nothing to rewrite.
	if appProberInfo == nil {
		return
	}
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
		newURL := fmt.Sprintf("/app-health/%v/readyz", c.Name)
		after := calculateRewrite2(c.ReadinessProbe, newURL, statusPort, portMap)
		if after != nil {
			*c.ReadinessProbe.HTTPGet = *after
		}
		newURL = fmt.Sprintf("/app-health/%v/livez", c.Name)
		after = calculateRewrite2(c.LivenessProbe, newURL, statusPort, portMap)
		if after != nil {
			*c.LivenessProbe.HTTPGet = *after
		}
	}

}
