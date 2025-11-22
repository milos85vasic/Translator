package deployment

import (
	"context"
	"testing"
	"time"

	"digital.vasic.translator/internal/config"
)

func TestNetworkDiscoverer_Discovery(t *testing.T) {
	cfg := &config.Config{}
	discoverer := NewNetworkDiscoverer(cfg, nil)
	defer discoverer.Close()

	// Test starting discovery
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := discoverer.StartDiscovery(ctx)
	if err != nil {
		t.Fatalf("Failed to start discovery: %v", err)
	}

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test getting discovered services (should be empty initially)
	services := discoverer.GetDiscoveredServices()
	if len(services) != 0 {
		t.Errorf("Expected 0 discovered services initially, got %d", len(services))
	}
}

func TestNetworkDiscoverer_Broadcasting(t *testing.T) {
	cfg := &config.Config{}
	discoverer := NewNetworkDiscoverer(cfg, nil)
	defer discoverer.Close()

	// Create mock instances
	instances := map[string]*DeployedInstance{
		"test-main": {
			ID:     "test-main",
			Host:   "localhost",
			Port:   8443,
			Status: "healthy",
		},
		"test-worker": {
			ID:     "test-worker",
			Host:   "localhost",
			Port:   8444,
			Status: "healthy",
		},
	}

	// Test starting broadcasting
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := discoverer.StartBroadcasting(ctx, instances)
	if err != nil {
		t.Fatalf("Failed to start broadcasting: %v", err)
	}

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)
}

func TestNetworkDiscoverer_ServiceTypeDetermination(t *testing.T) {
	tests := []struct {
		instanceID string
		expected   string
	}{
		{"translator-main", "coordinator"},
		{"translator-coordinator", "coordinator"},
		{"translator-worker-1", "worker"},
		{"worker-node-2", "worker"},
		{"unknown-service", "unknown"},
		{"gpu-worker", "worker"},
	}

	for _, tt := range tests {
		t.Run(tt.instanceID, func(t *testing.T) {
			cfg := &config.Config{}
			discoverer := NewNetworkDiscoverer(cfg, nil)

			result := discoverer.determineServiceType(tt.instanceID)
			if result != tt.expected {
				t.Errorf("determineServiceType(%s) = %s, want %s", tt.instanceID, result, tt.expected)
			}

			discoverer.Close()
		})
	}
}

func TestNetworkDiscoverer_Contains(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"translator-main", "main", true},
		{"translator-worker-1", "worker", true},
		{"gpu-worker", "worker", true},
		{"coordinator", "main", false},
		{"worker", "main", false},
		{"", "main", false},
		{"main", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.substr, func(t *testing.T) {
			// Test the contains function directly (it's a standalone function)
			result := contains(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("contains(%s, %s) = %v, want %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestNetworkService_TTL(t *testing.T) {
	cfg := &config.Config{}
	discoverer := NewNetworkDiscoverer(cfg, nil)
	defer discoverer.Close()

	// Add a service manually
	service := &NetworkService{
		ID:       "test-service",
		Name:     "Test Service",
		Host:     "localhost",
		Port:     8443,
		Type:     "worker",
		LastSeen: time.Now().Add(-2 * time.Minute), // 2 minutes ago
		TTL:      90 * time.Second,                 // 90 second TTL
	}

	discoverer.mu.Lock()
	discoverer.services["test-service"] = service
	discoverer.mu.Unlock()

	// Get services (should clean up expired ones)
	services := discoverer.GetDiscoveredServices()

	// The service should be cleaned up because it's past its TTL
	if _, exists := services["test-service"]; exists {
		t.Error("Expected expired service to be cleaned up")
	}
}
