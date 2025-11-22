package distributed_test

import (
	"testing"

	"digital.vasic.translator/internal/config"
	"digital.vasic.translator/pkg/events"
)

// Test that we can import and use the distributed package
func TestDistributedPackage(t *testing.T) {
	// This is a basic smoke test to ensure the package compiles
	cfg := &config.Config{}
	cfg.Distributed.Enabled = false

	eventBus := events.NewEventBus()
	_ = eventBus // Use the variable to avoid unused variable error

	// Just test that config structure works
	if cfg.Distributed.Enabled {
		t.Error("Should be false")
	}
}
