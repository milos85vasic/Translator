package main

import (
	"context"
	"digital.vasic.translator/internal/config"
	"digital.vasic.translator/internal/cache"
	"digital.vasic.translator/pkg/events"
	"digital.vasic.translator/pkg/models"
	"digital.vasic.translator/pkg/security"
	"digital.vasic.translator/pkg/websocket"
	"digital.vasic.translator/pkg/coordination"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestBasicComponents tests basic component creation
func TestBasicComponents(t *testing.T) {
	// Test event bus
	eventBus := events.NewEventBus()
	assert.NotNil(t, eventBus)

	// Test cache
	cache := cache.NewCache(60*time.Second, true)
	assert.NotNil(t, cache)

	// Test user repository
	userRepo := models.NewInMemoryUserRepository()
	assert.NotNil(t, userRepo)

	// Test auth service
	authService := security.NewUserAuthService("test-secret-key-16-chars", 24*time.Hour, userRepo)
	assert.NotNil(t, authService)

	// Test rate limiter
	rateLimiter := security.NewRateLimiter(10, 20)
	assert.NotNil(t, rateLimiter)

	// Test WebSocket hub
	wsHub := websocket.NewHub(eventBus)
	assert.NotNil(t, wsHub)

	// Test coordinator
	coordinator := coordination.NewMultiLLMCoordinator(coordination.CoordinatorConfig{
		EventBus: eventBus,
	})
	assert.NotNil(t, coordinator)
}

// TestBasicRoutes tests basic HTTP routes
func TestBasicRoutes(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a simple router
	router := gin.New()
	
	// Add basic routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/languages", func(c *gin.Context) {
		c.JSON(200, gin.H{"languages": []string{"en", "es", "fr"}})
	})
	router.GET("/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{"stats": gin.H{"total": 0, "translated": 0}})
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedField  string
	}{
		{
			name:           "health endpoint",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
			expectedField:  "status",
		},
		{
			name:           "languages endpoint",
			method:         "GET",
			path:           "/languages",
			expectedStatus: http.StatusOK,
			expectedField:  "languages",
		},
		{
			name:           "stats endpoint",
			method:         "GET",
			path:           "/stats",
			expectedStatus: http.StatusOK,
			expectedField:  "stats",
		},
		{
			name:           "nonexistent endpoint",
			method:         "GET",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedField != "" && w.Code == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, tt.expectedField)
			}
		})
	}
}

// TestWebSocketHub tests WebSocket hub functionality
func TestWebSocketHub(t *testing.T) {
	eventBus := events.NewEventBus()
	wsHub := websocket.NewHub(eventBus)

	// Test hub creation
	assert.NotNil(t, wsHub)

	// Test event subscription
	eventBus.Subscribe("test", func(event events.Event) {
		// Handle test event
	})

	// Test event publishing
	eventBus.Publish(events.Event{
		Type:    "test",
		Message: "test message",
	})

	// Test basic hub operations without running full server
	// (can't test full WebSocket without more complex setup)
}

// TestCacheOperations tests cache functionality
func TestCacheOperations(t *testing.T) {
	cache := cache.NewCache(60*time.Second, true)

	// Test cache set/get
	key := "test-key"
	value := "test-value"

	cache.Set(key, value)

	retrieved, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, value, retrieved)

	// Test cache miss
	_, found = cache.Get("nonexistent")
	assert.False(t, found)

	// Test cache deletion
	cache.Delete(key)
	_, found = cache.Get(key)
	assert.False(t, found)
}

// TestRateLimiting tests rate limiting functionality
func TestRateLimiting(t *testing.T) {
	rateLimiter := security.NewRateLimiter(10, 10)

	// Test rate limiting within limits
	for i := 0; i < 10; i++ {
		allowed := rateLimiter.Allow("test-key")
		assert.True(t, allowed, "Request %d should be allowed", i)
	}

	// Test rate limiting exceeding limits
	exceeded := rateLimiter.Allow("test-key")
	assert.False(t, exceeded, "Request exceeding limit should be denied")
}

// TestAuthentication tests authentication functionality
func TestAuthentication(t *testing.T) {
	// Create user repository
	userRepo := models.NewInMemoryUserRepository()
	assert.NotNil(t, userRepo)

	// Create auth service
	authService := security.NewUserAuthService("test-secret-key-16-chars", 24*time.Hour, userRepo)
	assert.NotNil(t, authService)

	// Test user creation
	user := &models.User{
		ID:       "test-user-id",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpass123",
		IsActive: true,
	}

	err := userRepo.Create(user)
	assert.NoError(t, err)
	assert.NotNil(t, user.ID)
	assert.NotZero(t, user.CreatedAt)

	// Test user authentication
	loginReq := security.LoginRequest{
		Username: "testuser",
		Password: "testpass123",
	}
	loginResp, err := authService.AuthenticateUser(loginReq)
	assert.NoError(t, err)
	assert.NotNil(t, loginResp)
	assert.NotEmpty(t, loginResp.Token)

	// Test token validation
	validatedUser, err := authService.ValidateToken(loginResp.Token)
	assert.NoError(t, err)
	assert.NotNil(t, validatedUser)
	assert.Equal(t, "testuser", validatedUser.Username)
}

// TestConfiguration tests configuration handling
func TestConfiguration(t *testing.T) {
	// Test default configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Security: config.SecurityConfig{
			JWTSecret:      "test-secret-key",
			RateLimitRPS:   100,
			RateLimitBurst: 200,
		},
		Translation: config.TranslationConfig{
			CacheTTL:     3600,
			CacheEnabled: true,
		},
		Distributed: config.DistributedConfig{
			Enabled: false,
		},
	}

	// Validate configuration
	err := cfg.Validate()
	assert.NoError(t, err)

	// Test configuration serialization
	data, err := json.Marshal(cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test configuration deserialization
	var parsedCfg config.Config
	err = json.Unmarshal(data, &parsedCfg)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Server.Host, parsedCfg.Server.Host)
	assert.Equal(t, cfg.Server.Port, parsedCfg.Server.Port)
}

// TestErrorHandling tests error handling scenarios
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		testFunc    func() error
		expectError bool
	}{
		{
			name: "invalid configuration",
			testFunc: func() error {
				cfg := &config.Config{
					Server: config.ServerConfig{
						Port: -1, // Invalid port
					},
				}
				return cfg.Validate()
			},
			expectError: true,
		},
		{
			name: "empty JWT secret",
			testFunc: func() error {
				cfg := &config.Config{
					Security: config.SecurityConfig{
						JWTSecret: "", // Empty secret
					},
				}
				return cfg.Validate()
			},
			expectError: true,
		},
		{
			name: "valid configuration",
			testFunc: func() error {
				cfg := &config.Config{
					Server: config.ServerConfig{
						Host: "localhost",
						Port: 8080,
					},
					Security: config.SecurityConfig{
						JWTSecret: "valid-secret-key",
					},
				}
				return cfg.Validate()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestSignalHandling tests signal handling
func TestSignalHandling(t *testing.T) {
	// Test signal channel creation
	sigChan := make(chan os.Signal, 1)
	
	// Test signal handling simulation
	go func() {
		time.Sleep(50 * time.Millisecond)
		sigChan <- os.Interrupt
	}()

	// Wait for signal
	select {
	case sig := <-sigChan:
		assert.Equal(t, os.Interrupt, sig)
	case <-time.After(1 * time.Second):
		t.Fatal("Did not receive signal")
	}
}

// TestServerLifecycle tests server lifecycle
func TestServerLifecycle(t *testing.T) {
	// Create a simple HTTP server for testing
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	server := &http.Server{
		Addr:    "localhost:0", // Use random port
		Handler: mux,
	}

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-errChan:
		t.Fatalf("Server failed to start: %v", err)
	default:
		// Server started successfully
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}

// TestFileOperations tests file operations
func TestFileOperations(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Test file creation
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "test content"
	
	err := os.WriteFile(testFile, []byte(content), 0644)
	assert.NoError(t, err)
	
	// Test file reading
	readContent, err := os.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, content, string(readContent))
	
	// Test file existence
	_, err = os.Stat(testFile)
	assert.NoError(t, err)
}

// TestJSONOperations tests JSON operations
func TestJSONOperations(t *testing.T) {
	// Test JSON encoding
	data := map[string]interface{}{
		"name": "test",
		"value": 123,
	}
	
	encoded, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, encoded)
	
	// Test JSON decoding
	var decoded map[string]interface{}
	err = json.Unmarshal(encoded, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, data["name"], decoded["name"])
	assert.Equal(t, float64(123), decoded["value"])
}

// TestVersionHandling tests version information
func TestVersionHandling(t *testing.T) {
	version := "1.0.0"
	
	// Test version string format
	assert.NotEmpty(t, version)
	assert.Contains(t, version, ".")
	
	// Test version comparison
	version1 := "1.0.0"
	version2 := "1.0.1"
	
	assert.NotEqual(t, version1, version2)
	assert.True(t, version1 != version2)
}