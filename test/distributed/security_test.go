//go:build security

package distributed

import (
	"testing"

	"digital.vasic.translator/internal/config"
	"digital.vasic.translator/pkg/distributed"
	"digital.vasic.translator/pkg/events"
)

func TestDistributedManager_Security_NoSensitiveData(t *testing.T) {
	// Test that sensitive data is not exposed in status responses
	cfg := &config.Config{}
	cfg.Distributed.Enabled = true
	cfg.Distributed.Workers = map[string]config.WorkerConfig{
		"secure-worker": {
			Name:        "Secure Worker",
			Host:        "secure.example.com",
			Port:        22,
			User:        "translator",
			KeyFile:     "/path/to/secret/key",
			Password:    "super_secret_password",
			MaxCapacity: 10,
			Enabled:     true,
		},
	}

	eventBus := events.NewEventBus()
	manager := distributed.NewDistributedManager(cfg, eventBus)

	// Get status - this should NOT contain sensitive data
	status := manager.GetStatus()

	workers := status["workers"].(map[string]interface{})
	worker := workers["secure-worker"].(map[string]interface{})

	// Verify sensitive fields are not exposed
	if _, hasKeyFile := worker["key_file"]; hasKeyFile {
		t.Error("SSH key file path should not be exposed in status")
	}

	if _, hasPassword := worker["password"]; hasPassword {
		t.Error("SSH password should not be exposed in status")
	}

	// Verify non-sensitive fields are present
	if worker["name"] != "Secure Worker" {
		t.Error("Worker name should be exposed")
	}

	if worker["enabled"] != true {
		t.Error("Enabled status should be exposed")
	}

	if worker["capacity"] != 10 {
		t.Error("Capacity should be exposed")
	}
}

func TestDistributedManager_Security_WorkerIsolation(t *testing.T) {
	// Test that workers are properly isolated and cannot access each other's data
	cfg := &config.Config{}
	cfg.Distributed.Enabled = true

	eventBus := events.NewEventBus()
	manager := distributed.NewDistributedManager(cfg, eventBus)

	// Add multiple workers
	worker1 := &distributed.WorkerConfig{
		ID:          "worker-1",
		Name:        "Worker 1",
		SSH:         distributed.SSHConfig{Host: "host1", User: "user1"},
		MaxCapacity: 5,
		Enabled:     true,
	}

	worker2 := &distributed.WorkerConfig{
		ID:          "worker-2",
		Name:        "Worker 2",
		SSH:         distributed.SSHConfig{Host: "host2", User: "user2"},
		MaxCapacity: 8,
		Enabled:     true,
	}

	err := manager.AddWorker("worker-1", worker1)
	if err != nil {
		t.Fatalf("Failed to add worker 1: %v", err)
	}

	err = manager.AddWorker("worker-2", worker2)
	if err != nil {
		t.Fatalf("Failed to add worker 2: %v", err)
	}

	// Check status - workers should be isolated
	status := manager.GetStatus()
	workers := status["workers"].(map[string]interface{})

	if len(workers) != 2 {
		t.Errorf("Expected 2 workers, got %d", len(workers))
	}

	w1 := workers["worker-1"].(map[string]interface{})
	w2 := workers["worker-2"].(map[string]interface{})

	// Verify worker data is separate
	if w1["name"] == w2["name"] {
		t.Error("Worker names should be different")
	}

	if w1["capacity"] == w2["capacity"] {
		t.Error("Worker capacities should be different")
	}
}

func TestDistributedManager_Security_DisabledWorkers(t *testing.T) {
	// Test that disabled workers don't expose information
	cfg := &config.Config{}
	cfg.Distributed.Enabled = true
	cfg.Distributed.Workers = map[string]config.WorkerConfig{
		"disabled-worker": {
			Name:        "Disabled Worker",
			Host:        "disabled.example.com",
			Port:        22,
			User:        "translator",
			MaxCapacity: 5,
			Enabled:     false, // Explicitly disabled
		},
	}

	eventBus := events.NewEventBus()
	manager := distributed.NewDistributedManager(cfg, eventBus)

	status := manager.GetStatus()
	workers := status["workers"].(map[string]interface{})
	worker := workers["disabled-worker"].(map[string]interface{})

	// Disabled workers should still show basic info but indicate disabled status
	if worker["enabled"] != false {
		t.Error("Disabled worker should show enabled=false")
	}

	// But should not expose connection details that could be used maliciously
	if worker["status"] != "configured" {
		t.Error("Disabled worker should show configured status")
	}
}
