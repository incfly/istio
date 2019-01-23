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

// // extractStatusPort returns the port value of the pilot agent sidecar container statusPort.
// // Return -1 if not found.
// func extractStatusPort(spec *SidecarInjectionSpec) int {
// 	if spec == nil {
// 		return -1
// 	}
// 	statusPort := -1
// 	for _, c := range spec.Containers {
// 		if c.Name != istioProxyContainerName {
// 			continue
// 		}
// 		for i, arg := range c.Args {
// 			// Skip for unrelated args.
// 			match := statusPortPattern.FindAllStringSubmatch(strings.TrimSpace(arg), -1)
// 			if len(match) != 1 {
// 				continue
// 			}
// 			groups := statusPortPattern.SubexpNames()
// 			portStr := ""
// 			for ind, s := range match[0] {
// 				if groups[ind] == "port" {
// 					portStr = s
// 					break
// 				}
// 			}
// 			// Port not found from current arg, extract from next arg.
// 			if portStr == "" {
// 				// Matches the regex pattern, but without actual values provided.
// 				if len(c.Args) <= i+1 {
// 					log.Errorf("No statusPort value provided, skip app probe rewriting")
// 					return -1
// 				}
// 				portStr = c.Args[i+1]
// 			}
// 			p, err := strconv.Atoi(portStr)
// 			if err != nil {
// 				log.Errorf("Failed to convert statusPort to int %v, err %v", portStr, err)
// 				return -1
// 			}
// 			statusPort = p
// 			break
// 		}
// 	}
// 	return statusPort
// }

// createProbeRewritePatch generates the patch for webhook.
func createProbeRewritePatch(podSpec *corev1.PodSpec, spec *SidecarInjectionSpec) []rfc6902PatchOperation {
	patch := []rfc6902PatchOperation{}
	if spec == nil || podSpec == nil || !spec.RewriteAppHTTPProbe {
		return patch
	}
	var sidecar *corev1.Container
	for i := range podSpec.Containers {
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
	// Change the application containers' probe to point to sidecar's status port.
	// TODO: here rewrite the rewriteProbe function in new approach.
	rewriteProbe := func(probe *corev1.Probe, portMap map[string]int32, path string) *rfc6902PatchOperation {
		return nil
		// if probe == nil || probe.HTTPGet == nil {
		// 	return nil
		// }
		// httpGet := proto.Clone(probe.HTTPGet).(*corev1.HTTPGetAction)
		// // note(incfly): workaround... proto.Clone can't copy corev1.IntOrStr somehow.
		// httpGet.Port = probe.HTTPGet.Port
		// header := corev1.HTTPHeader{
		// 	Name:  status.IstioAppPortHeader,
		// 	Value: httpGet.Port.String(),
		// }
		// // A named port, resolve by looking at port map.
		// if httpGet.Port.Type == intstr.String {
		// 	port, exists := portMap[httpGet.Port.StrVal]
		// 	if !exists {
		// 		log.Errorf("named port not found in the map skip rewriting probing %v", *probe)
		// 		return nil
		// 	}
		// 	header.Value = strconv.Itoa(int(port))
		// }
		// httpGet.HTTPHeaders = append(httpGet.HTTPHeaders, header)
		// httpGet.Port = intstr.FromInt(statusPort)
		// return &rfc6902PatchOperation{
		// 	Op:    "replace",
		// 	Path:  path,
		// 	Value: *httpGet,
		// }
	}
	for i, c := range podSpec.Containers {
		// Skip sidecar container.
		if c.Name == istioProxyContainerName {
			continue
		}
		portMap := map[string]int32{}
		for _, p := range c.Ports {
			portMap[p.Name] = p.ContainerPort
		}
		if p := rewriteProbe(c.ReadinessProbe, portMap, fmt.Sprintf("/spec/containers/%v/readinessProbe/httpGet", i)); p != nil {
			patch = append(patch, *p)
		}
		if p := rewriteProbe(c.LivenessProbe, portMap, fmt.Sprintf("/spec/containers/%v/livenessProbe/httpGet", i)); p != nil {
			patch = append(patch, *p)
		}
	}
	return patch
}

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

// rewriteProbe changes application containers' probe to point to sidecar's status port.
func rewriteProbe(probe *corev1.Probe, appProbers *status.KubeAppProbers,
	newURL string, statusPort int, portMap map[string]int32) {
	if probe == nil || probe.HTTPGet == nil {
		return
	}
	httpGet := probe.HTTPGet

	// Save app probe config to pass to pilot agent later.
	savedProbe := proto.Clone(probe.HTTPGet).(*corev1.HTTPGetAction)
	(*appProbers)[newURL] = savedProbe
	// A named port, resolve by looking at port map.
	if httpGet.Port.Type == intstr.String {
		port, exists := portMap[httpGet.Port.StrVal]
		if !exists {
			log.Errorf("named port not found in the map skip rewriting probing %v", *probe)
			return
		}
		savedProbe.Port = intstr.FromInt(int(port))
	}
	// Change the application csince ontainer prober config.
	httpGet.Port = intstr.FromInt(statusPort)
	httpGet.Path = newURL
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
		rewriteProbe(c.ReadinessProbe, &appProberInfo, fmt.Sprintf("/app-health/%v/readyz", c.Name), statusPort, portMap)
		rewriteProbe(c.LivenessProbe, &appProberInfo, fmt.Sprintf("/app-health/%v/livez", c.Name), statusPort, portMap)
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
