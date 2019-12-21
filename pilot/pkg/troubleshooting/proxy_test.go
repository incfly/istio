package troubleshooting

import (
	"testing"
)

func TestProxyServer(t *testing.T) {
	_, err := NewProxyServer()
	if err != nil {
		t.Errorf("failed to create server %v", err)
	}
}

// pilot agent behavior.
func TestProxyClient(t *testing.T) {
}
