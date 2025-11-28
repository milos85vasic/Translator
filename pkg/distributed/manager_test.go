package distributed

import (
	"testing"

	"digital.vasic.translator/internal/config"
	"digital.vasic.translator/pkg/deployment"
	"digital.vasic.translator/pkg/events"
)

func TestDefaultLogger_Log(t *testing.T) {
	t.Run("BasicLogging", func(t *testing.T) {
		logger := &defaultLogger{}
		
		// Test that logging doesn't panic
		logger.Log("info", "test message", map[string]interface{}{
			"key": "value",
		})
		
		logger.Log("error", "error message", nil)
		
		// No assertions needed - just verify it doesn't panic
	})
}

func TestNewDistributedManager(t *testing.T) {
	t.Run("Constructor", func(t *testing.T) {
		cfg := &config.Config{
			Distributed: config.DistributedConfig{
				Workers: make(map[string]config.WorkerConfig),
			},
		}
		
		eventBus := events.NewEventBus()
		apiLogger, _ := deployment.NewAPICommunicationLogger("/tmp/test-api.log")
		
		dm := NewDistributedManager(cfg, eventBus, apiLogger)
		
		if dm == nil {
			t.Error("Expected non-nil DistributedManager")
		}
		
		if dm.config != cfg {
			t.Error("Expected config to be set correctly")
		}
		
		if dm.initialized {
			t.Error("Expected initialized to be false initially")
		}
	})
}

func TestDistributedManager_GetWorkerByID(t *testing.T) {
	t.Run("NotInitialized", func(t *testing.T) {
		cfg := &config.Config{
			Distributed: config.DistributedConfig{
				Workers: make(map[string]config.WorkerConfig),
			},
		}
		
		eventBus := events.NewEventBus()
		apiLogger, _ := deployment.NewAPICommunicationLogger("/tmp/test-api.log")
		
		dm := NewDistributedManager(cfg, eventBus, apiLogger)
		
		// Test when not initialized
		worker := dm.GetWorkerByID("test-worker")
		if worker != nil {
			t.Error("Expected nil worker when not initialized")
		}
	})
	
	t.Run("Initialized", func(t *testing.T) {
		cfg := &config.Config{
			Distributed: config.DistributedConfig{
				Workers: make(map[string]config.WorkerConfig),
			},
		}
		
		eventBus := events.NewEventBus()
		apiLogger, _ := deployment.NewAPICommunicationLogger("/tmp/test-api.log")
		
		dm := NewDistributedManager(cfg, eventBus, apiLogger)
		
		// Initialize
		err := dm.Initialize(nil)
		if err != nil {
			t.Errorf("Failed to initialize: %v", err)
		}
		
		// Test with non-existent worker
		worker := dm.GetWorkerByID("non-existent-worker")
		if worker != nil {
			t.Error("Expected nil worker for non-existent worker")
		}
	})
}

func TestDistributedManager_GetStatus(t *testing.T) {
	t.Run("NotInitialized", func(t *testing.T) {
		cfg := &config.Config{
			Distributed: config.DistributedConfig{
				Workers: make(map[string]config.WorkerConfig),
			},
		}
		
		eventBus := events.NewEventBus()
		apiLogger, _ := deployment.NewAPICommunicationLogger("/tmp/test-api.log")
		
		dm := NewDistributedManager(cfg, eventBus, apiLogger)
		
		// Test when not initialized
		status := dm.GetStatus()
		if status == nil {
			t.Error("Expected non-nil status map")
		}
		
		// Should have basic fields
		if _, exists := status["initialized"]; !exists {
			t.Error("Expected 'initialized' field in status")
		}
		
		if status["initialized"].(bool) {
			t.Error("Expected initialized to be false")
		}
	})
	
	t.Run("Initialized", func(t *testing.T) {
		cfg := &config.Config{
			Distributed: config.DistributedConfig{
				Workers: make(map[string]config.WorkerConfig),
			},
		}
		
		eventBus := events.NewEventBus()
		apiLogger, _ := deployment.NewAPICommunicationLogger("/tmp/test-api.log")
		
		dm := NewDistributedManager(cfg, eventBus, apiLogger)
		
		// Initialize
		err := dm.Initialize(nil)
		if err != nil {
			t.Errorf("Failed to initialize: %v", err)
		}
		
		// Test when initialized
		status := dm.GetStatus()
		if status == nil {
			t.Error("Expected non-nil status map")
		}
		
		// Should have basic fields
		if _, exists := status["initialized"]; !exists {
			t.Error("Expected 'initialized' field in status")
		}
		
		if !status["initialized"].(bool) {
			t.Error("Expected initialized to be true")
		}
		
		// Should have workers field
		if _, exists := status["workers"]; !exists {
			t.Error("Expected 'workers' field in status")
		}
	})
}

func TestDistributedManager_GetPairedServices(t *testing.T) {
	t.Run("EmptyServices", func(t *testing.T) {
		cfg := &config.Config{
			Distributed: config.DistributedConfig{
				Workers: make(map[string]config.WorkerConfig),
			},
		}
		
		eventBus := events.NewEventBus()
		apiLogger, _ := deployment.NewAPICommunicationLogger("/tmp/test-api.log")
		
		dm := NewDistributedManager(cfg, eventBus, apiLogger)
		
		// Test without initialization
		services := dm.GetPairedServices()
		if services == nil {
			t.Error("Expected non-nil services map")
		}
		
		if len(services) != 0 {
			t.Error("Expected empty services map")
		}
	})
}

func TestDistributedManager_Close(t *testing.T) {
	t.Run("CloseGracefully", func(t *testing.T) {
		cfg := &config.Config{
			Distributed: config.DistributedConfig{
				Workers: make(map[string]config.WorkerConfig),
			},
		}
		
		eventBus := events.NewEventBus()
		apiLogger, _ := deployment.NewAPICommunicationLogger("/tmp/test-api.log")
		
		dm := NewDistributedManager(cfg, eventBus, apiLogger)
		
		// Should not panic when closing
		dm.Close()
	})
}