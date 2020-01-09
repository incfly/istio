// package envoy bla.
// separate file is on purpose to avoid merging upstream conflict.
package envoy

import (
	"fmt"
	ts "istio.io/istio/pilot/pkg/troubleshooting"
)

func (a *agent) runTroubleShooting() error {
	// todo: later on move it to constructor config, rather than here...
	if a.troubleshootingAgent == nil {
		ta, err := ts.NewAgent(&ts.AgentConfig{
			// Everything is hardcoded for now.
			ID:             "sidecar-proxy.default.httpbin",
			ServiceAddress: "ts-server.istio-system.svc.cluster.local:8000",
		})
		if err != nil {
			return err
		}
		a.troubleshootingAgent = ta
	}
	if err := a.troubleshootingAgent.Start(); err != nil {
		return fmt.Errorf("failed to start trouble shooting agent... %v", err)
	}
	return nil
}
