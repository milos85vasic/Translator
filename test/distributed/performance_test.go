//go:build performance

package distributed

import (
	"context"
	"testing"
	"time"

	"digital.vasic.translator/internal/config"
	"digital.vasic.translator/pkg/coordination"
	"digital.vasic.translator/pkg/events"
)

func BenchmarkDistributedManager_GetStatus(b *testing.B) {
	cfg := &config.Config{}
	cfg.Distributed.Enabled = true

	// Add multiple workers for realistic benchmark
	for i := 0; i < 10; i++ {
		cfg.Distributed.Workers = map[string]config.WorkerConfig{
			"worker-1": {Name: "Worker 1", Host: "host1", Port: 22, User: "user", MaxCapacity: 5, Enabled: true},
			"worker-2": {Name: "Worker 2", Host: "host2", Port: 22, User: "user", MaxCapacity: 5, Enabled: true},
			"worker-3": {Name: "Worker 3", Host: "host3", Port: 22, User: "user", MaxCapacity: 5, Enabled: true},
		}
	}

	eventBus := events.NewEventBus()
	localCoordinator := coordination.NewMultiLLMCoordinator(coordination.CoordinatorConfig{
		EventBus: eventBus,
	})

	manager := NewDistributedManager(cfg, eventBus)
	manager.Initialize(localCoordinator)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = manager.GetStatus()
	}
}

func BenchmarkDistributedManager_AddWorker(b *testing.B) {
	cfg := &config.Config{}
	cfg.Distributed.Enabled = true

	eventBus := events.NewEventBus()
	localCoordinator := coordination.NewMultiLLMCoordinator(coordination.CoordinatorConfig{
		EventBus: eventBus,
	})

	manager := NewDistributedManager(cfg, eventBus)
	manager.Initialize(localCoordinator)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		worker := &WorkerConfig{
			ID:          "bench-worker",
			Name:        "Benchmark Worker",
			SSH:         SSHConfig{Host: "bench-host", User: "bench-user"},
			MaxCapacity: 5,
			Enabled:     true,
		}

		manager.AddWorker("bench-worker", worker)
		manager.RemoveWorker("bench-worker") // Clean up
	}
}

func BenchmarkDistributedCoordinator_Translate(b *testing.B) {
	cfg := &config.Config{}
	cfg.Distributed.Enabled = true

	eventBus := events.NewEventBus()
	localCoordinator := coordination.NewMultiLLMCoordinator(coordination.CoordinatorConfig{
		EventBus: eventBus,
	})

	manager := NewDistributedManager(cfg, eventBus)
	manager.Initialize(localCoordinator)

	ctx := context.Background()
	text := "This is a test sentence for benchmarking distributed translation performance."

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// This will fail but we want to measure the overhead
		_, _ = manager.TranslateDistributed(ctx, text, "benchmark")
	}
}

func BenchmarkSSHPool_GetConnection(b *testing.B) {
	pool := NewSSHPool()

	// Add a mock worker config
	config := &WorkerConfig{
		ID:          "bench-ssh",
		Name:        "SSH Benchmark",
		SSH:         SSHConfig{Host: "localhost", User: "test", Timeout: time.Second},
		MaxCapacity: 5,
		Enabled:     false, // Don't actually try to connect
	}

	pool.AddWorker(config)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// This will return nil since worker is disabled, but measures lookup overhead
		_, _ = pool.GetConnection("bench-ssh")
	}
}
