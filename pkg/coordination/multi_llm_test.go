package coordination

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNewMultiLLMCoordinator_NoAPIKeys(t *testing.T) {
	// Clear all API keys
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ZHIPU_API_KEY")
	os.Unsetenv("DEEPSEEK_API_KEY")
	os.Unsetenv("QWEN_API_KEY")
	os.Unsetenv("OLLAMA_ENABLED")

	// Disable Qwen OAuth discovery in tests
	os.Setenv("SKIP_QWEN_OAUTH", "1")
	defer os.Unsetenv("SKIP_QWEN_OAUTH")

	coordinator := NewMultiLLMCoordinator(CoordinatorConfig{
		SessionID: "test-session",
	})

	if coordinator.GetInstanceCount() != 0 {
		t.Errorf("Expected 0 instances with no API keys, got %d", coordinator.GetInstanceCount())
	}
}

func TestNewMultiLLMCoordinator_WithAPIKeys(t *testing.T) {
	// Set some API keys
	os.Setenv("OPENAI_API_KEY", "sk-test-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	os.Setenv("DEEPSEEK_API_KEY", "sk-deepseek-test")
	defer os.Unsetenv("DEEPSEEK_API_KEY")

	coordinator := NewMultiLLMCoordinator(CoordinatorConfig{
		SessionID: "test-session",
	})

	// Should create instances for the API keys
	if coordinator.GetInstanceCount() < 2 {
		t.Errorf("Expected at least 2 instances with API keys, got %d", coordinator.GetInstanceCount())
	}
}

func TestMultiLLMCoordinator_TranslateWithRetry_NoInstances(t *testing.T) {
	coordinator := &MultiLLMCoordinator{
		instances: make([]*LLMInstance, 0),
	}

	_, err := coordinator.TranslateWithRetry(context.Background(), "test text", "")
	if err == nil {
		t.Error("Expected error when no instances available")
	}

	if err.Error() != "no LLM instances available" {
		t.Errorf("Expected 'no LLM instances available' error, got: %v", err)
	}
}

func TestMultiLLMCoordinator_GetNextInstance(t *testing.T) {
	coordinator := &MultiLLMCoordinator{
		instances: []*LLMInstance{
			{ID: "instance1", Available: true},
			{ID: "instance2", Available: false},
			{ID: "instance3", Available: true},
		},
		currentIndex: 0,
	}

	// Should return first available instance
	instance := coordinator.getNextInstance()
	if instance == nil {
		t.Fatal("Expected instance, got nil")
	}
	if instance.ID != "instance1" {
		t.Errorf("Expected instance1, got %s", instance.ID)
	}

	// Should skip unavailable instance and return next available
	instance = coordinator.getNextInstance()
	if instance == nil {
		t.Fatal("Expected instance, got nil")
	}
	if instance.ID != "instance3" {
		t.Errorf("Expected instance3, got %s", instance.ID)
	}
}

func TestLLMInstance_Availability(t *testing.T) {
	instance := &LLMInstance{
		ID:        "test-instance",
		Available: true,
	}

	if !instance.Available {
		t.Error("Instance should be available initially")
	}

	// Mark as unavailable
	instance.Available = false
	if instance.Available {
		t.Error("Instance should be unavailable")
	}
}

func TestDiscoverProviders(t *testing.T) {
	// Clear environment
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ZHIPU_API_KEY")
	os.Unsetenv("DEEPSEEK_API_KEY")
	os.Unsetenv("OLLAMA_ENABLED")
	os.Setenv("SKIP_QWEN_OAUTH", "1")
	defer os.Unsetenv("SKIP_QWEN_OAUTH")

	coordinator := &MultiLLMCoordinator{}

	providers := coordinator.discoverProviders()

	// May have some providers from OAuth or other sources, but should be minimal
	initialCount := len(providers)

	// Set some keys
	os.Setenv("OPENAI_API_KEY", "test-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	os.Setenv("DEEPSEEK_API_KEY", "deepseek-key")
	defer os.Unsetenv("DEEPSEEK_API_KEY")

	providers = coordinator.discoverProviders()

	// Should find at least the providers we set
	if len(providers) < initialCount+2 {
		t.Errorf("Expected at least %d providers with API keys, got %d", initialCount+2, len(providers))
	}

	if _, exists := providers["openai"]; !exists {
		t.Error("Expected openai provider to be discovered")
	}

	if _, exists := providers["deepseek"]; !exists {
		t.Error("Expected deepseek provider to be discovered")
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_VAR", "test-value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnvOrDefault("TEST_VAR", "default")
	if result != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", result)
	}

	// Test with non-existing env var
	result = getEnvOrDefault("NON_EXISTING_VAR", "default-value")
	if result != "default-value" {
		t.Errorf("Expected 'default-value', got '%s'", result)
	}
}

func TestCoordinatorConfig_Defaults(t *testing.T) {
	config := CoordinatorConfig{}

	if config.MaxRetries != 0 {
		t.Errorf("Expected MaxRetries to be 0 initially, got %d", config.MaxRetries)
	}

	if config.RetryDelay != 0 {
		t.Errorf("Expected RetryDelay to be 0 initially, got %v", config.RetryDelay)
	}

	// Test that NewMultiLLMCoordinator sets defaults
	coordinator := NewMultiLLMCoordinator(config)

	if coordinator.maxRetries != 3 {
		t.Errorf("Expected maxRetries to be 3, got %d", coordinator.maxRetries)
	}

	if coordinator.retryDelay != 2*time.Second {
		t.Errorf("Expected retryDelay to be 2s, got %v", coordinator.retryDelay)
	}
}
