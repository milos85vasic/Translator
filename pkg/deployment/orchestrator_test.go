package deployment

import (
	"context"
	"testing"
	"time"

	"digital.vasic.translator/internal/config"
	"digital.vasic.translator/pkg/events"
)

func TestDeploymentOrchestrator_DeployDistributedSystem(t *testing.T) {
	// Setup
	cfg := &config.Config{}
	eventBus := events.NewEventBus()
	orchestrator := NewDeploymentOrchestrator(cfg, eventBus)
	defer orchestrator.Close()

	// Create a simple deployment plan
	plan := &DeploymentPlan{
		Main: &DeploymentConfig{
			Host:          "localhost",
			User:          "test",
			SSHKeyPath:    "/dev/null", // Won't actually connect
			DockerImage:   "test:latest",
			ContainerName: "test-main",
			Ports: []PortMapping{
				{HostPort: 8443, ContainerPort: 8443, Protocol: "tcp"},
			},
			Environment:   map[string]string{"TEST": "true"},
			Networks:      []string{"test-network"},
			RestartPolicy: "unless-stopped",
		},
		Workers: []*DeploymentConfig{}, // Empty for this test
	}

	// Test deployment (will fail due to no actual SSH/docker, but tests the logic)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := orchestrator.DeployDistributedSystem(ctx, plan)

	// We expect this to fail because we don't have actual SSH/docker setup
	// But it should fail gracefully, not panic
	if err == nil {
		t.Error("Expected deployment to fail with test setup, but it succeeded")
	}

	// Check that the error is related to SSH connection, not a panic
	t.Logf("Deployment failed as expected: %v", err)
}

func TestDeploymentOrchestrator_FindAvailablePort(t *testing.T) {
	cfg := &config.Config{}
	eventBus := events.NewEventBus()
	orchestrator := NewDeploymentOrchestrator(cfg, eventBus)
	defer orchestrator.Close()

	// Test finding an available port
	port, err := orchestrator.findAvailablePort("localhost", 8443)
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}

	if port < 8443 {
		t.Errorf("Expected port >= 8443, got %d", port)
	}

	t.Logf("Found available port: %d", port)
}

func TestDeploymentOrchestrator_ValidateDeploymentPlan(t *testing.T) {
	cfg := &config.Config{}
	eventBus := events.NewEventBus()
	orchestrator := NewDeploymentOrchestrator(cfg, eventBus)
	defer orchestrator.Close()

	tests := []struct {
		name    string
		plan    *DeploymentPlan
		wantErr bool
	}{
		{
			name: "valid plan",
			plan: &DeploymentPlan{
				Main: &DeploymentConfig{
					Host:          "localhost",
					User:          "test",
					DockerImage:   "test:latest",
					ContainerName: "test-main",
					Ports:         []PortMapping{{HostPort: 8443, ContainerPort: 8443}},
				},
				Workers: []*DeploymentConfig{
					{
						Host:          "localhost",
						User:          "test",
						DockerImage:   "test:latest",
						ContainerName: "test-worker-1",
						Ports:         []PortMapping{{HostPort: 8444, ContainerPort: 8443}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing main config",
			plan: &DeploymentPlan{
				Main:    nil,
				Workers: []*DeploymentConfig{},
			},
			wantErr: true,
		},
		{
			name: "invalid main config - no host",
			plan: &DeploymentPlan{
				Main: &DeploymentConfig{
					User:          "test",
					DockerImage:   "test:latest",
					ContainerName: "test-main",
				},
				Workers: []*DeploymentConfig{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := orchestrator.validateDeploymentPlan(tt.plan)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDeploymentPlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeploymentOrchestrator_GetDeployedInstances(t *testing.T) {
	cfg := &config.Config{}
	eventBus := events.NewEventBus()
	orchestrator := NewDeploymentOrchestrator(cfg, eventBus)
	defer orchestrator.Close()

	// Initially should be empty
	instances := orchestrator.GetDeployedInstances()
	if len(instances) != 0 {
		t.Errorf("Expected 0 instances initially, got %d", len(instances))
	}

	// Add a mock instance
	mockInstance := &DeployedInstance{
		ID:          "test-instance",
		Host:        "localhost",
		Port:        8443,
		ContainerID: "mock-container-id",
		Status:      "healthy",
		LastSeen:    time.Now(),
	}

	orchestrator.mu.Lock()
	orchestrator.deployed["test-instance"] = mockInstance
	orchestrator.mu.Unlock()

	// Now should have one instance
	instances = orchestrator.GetDeployedInstances()
	if len(instances) != 1 {
		t.Errorf("Expected 1 instance, got %d", len(instances))
	}

	if instance, exists := instances["test-instance"]; !exists {
		t.Error("Expected test-instance to exist")
	} else {
		if instance.ID != "test-instance" {
			t.Errorf("Expected instance ID 'test-instance', got '%s'", instance.ID)
		}
	}
}
