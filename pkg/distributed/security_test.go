package distributed

import (
	"crypto/rand"
	"crypto/rsa"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSecurityConfig_createHostKeyCallback(t *testing.T) {
	// Create a temporary known hosts file
	tmpDir := t.TempDir()
	knownHostsFile := filepath.Join(tmpDir, "known_hosts")

	// Generate a test SSH key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Failed to create public key: %v", err)
	}

	// Write a known hosts entry
	knownHostsContent := "example.com " + string(ssh.MarshalAuthorizedKey(publicKey))
	err = os.WriteFile(knownHostsFile, []byte(knownHostsContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write known hosts file: %v", err)
	}

	// Test with valid known hosts file
	config := &SecurityConfig{
		SSHHostKeyVerification: true,
		KnownHostsFile:         knownHostsFile,
	}

	callback, err := config.createHostKeyCallback()
	if err != nil {
		t.Fatalf("Failed to create host key callback: %v", err)
	}

	// Test valid key
	err = callback("example.com", &testAddr{}, publicKey)
	if err != nil {
		t.Errorf("Valid key should be accepted: %v", err)
	}

	// Generate a different key for testing rejection
	privateKey2, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate second test key: %v", err)
	}

	publicKey2, err := ssh.NewPublicKey(&privateKey2.PublicKey)
	if err != nil {
		t.Fatalf("Failed to create second public key: %v", err)
	}

	// Test invalid key
	err = callback("example.com", &testAddr{}, publicKey2)
	if err == nil {
		t.Error("Invalid key should be rejected")
	}
}

func TestSecurityConfig_loadKnownHosts(t *testing.T) {
	tmpDir := t.TempDir()
	knownHostsFile := filepath.Join(tmpDir, "known_hosts")

	// Test with non-existent file
	config := &SecurityConfig{KnownHostsFile: knownHostsFile}
	callback2, err := config.loadKnownHosts(knownHostsFile)
	if err != nil {
		t.Fatalf("Should not fail with non-existent known hosts file: %v", err)
	}

	// Should reject any key for non-existent file
	err = callback2("example.com", &testAddr{}, nil)
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Error("Non-existent known hosts file should reject all keys")
	}

	// Create empty file
	err = os.WriteFile(knownHostsFile, []byte(""), 0600)
	if err != nil {
		t.Fatalf("Failed to create empty known hosts file: %v", err)
	}

	// Test with empty file
	callback, err := config.loadKnownHosts(knownHostsFile)
	if err != nil {
		t.Fatalf("Failed to load empty known hosts file: %v", err)
	}

	// Should reject any key
	err = callback("example.com", &testAddr{}, nil)
	if err == nil || !strings.Contains(err.Error(), "no matching key found") {
		t.Error("Empty known hosts file should reject all keys")
	}
}

func TestKeysEqual(t *testing.T) {
	// Generate two identical keys
	privateKey1, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate first key: %v", err)
	}

	publicKey1, err := ssh.NewPublicKey(&privateKey1.PublicKey)
	if err != nil {
		t.Fatalf("Failed to create first public key: %v", err)
	}

	publicKey2, err := ssh.NewPublicKey(&privateKey1.PublicKey)
	if err != nil {
		t.Fatalf("Failed to create second public key: %v", err)
	}

	// Same key should be equal
	if !keysEqual(publicKey1, publicKey2) {
		t.Error("Identical keys should be equal")
	}

	// Generate different key
	privateKey3, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate third key: %v", err)
	}

	publicKey3, err := ssh.NewPublicKey(&privateKey3.PublicKey)
	if err != nil {
		t.Fatalf("Failed to create third public key: %v", err)
	}

	// Different keys should not be equal
	if keysEqual(publicKey1, publicKey3) {
		t.Error("Different keys should not be equal")
	}
}

// testAddr implements net.Addr for testing
type testAddr struct{}

func (a *testAddr) Network() string { return "tcp" }
func (a *testAddr) String() string  { return "127.0.0.1:22" }

// MockLogger for security testing
type MockSecurityLogger struct{}

func (m *MockSecurityLogger) Log(level, message string, fields map[string]interface{}) {
	// Do nothing in tests
}

func TestDefaultSecurityConfig(t *testing.T) {
	t.Run("DefaultConfiguration", func(t *testing.T) {
		config := DefaultSecurityConfig()
			
		if config == nil {
			t.Error("Expected non-nil security config")
		}
		
		// Check SSH settings
		if !config.SSHHostKeyVerification {
			t.Error("Expected SSH host key verification to be enabled")
		}
		
		if config.KnownHostsFile != "~/.ssh/known_hosts" {
			t.Errorf("Expected known hosts file '~/.ssh/known_hosts', got '%s'", config.KnownHostsFile)
		}
		
		if len(config.SSHCiphers) == 0 {
			t.Error("Expected SSH ciphers to be configured")
		}
		
		if len(config.SSHKexAlgorithms) == 0 {
			t.Error("Expected SSH key exchange algorithms to be configured")
		}
		
		if len(config.SSHMACs) == 0 {
			t.Error("Expected SSH MACs to be configured")
		}
		
		// Check TLS settings
		if !config.TLSCertVerification {
			t.Error("Expected TLS certificate verification to be enabled")
		}
		
		if len(config.TLSCipherSuites) == 0 {
			t.Error("Expected TLS cipher suites to be configured")
		}
		
		// Check connection limits
		if config.MaxConnectionsPerWorker != 5 {
			t.Errorf("Expected max connections per worker to be 5, got %d", config.MaxConnectionsPerWorker)
		}
		
		if config.ConnectionTimeout != 30*time.Second {
			t.Errorf("Expected connection timeout to be 30s, got %v", config.ConnectionTimeout)
		}
	})
}

func TestNewSecurityAuditor(t *testing.T) {
	t.Run("Constructor", func(t *testing.T) {
		logger := &MockSecurityLogger{}
		
		// Test enabled auditor
		auditor := NewSecurityAuditor(true, logger)
		
		if auditor == nil {
			t.Error("Expected non-nil security auditor")
		}
		
		if !auditor.enabled {
			t.Error("Expected auditor to be enabled")
		}
		
		if auditor.logger != logger {
			t.Error("Expected logger to be set correctly")
		}
		
		// Test disabled auditor
		disabledAuditor := NewSecurityAuditor(false, logger)
		
		if disabledAuditor.enabled {
			t.Error("Expected auditor to be disabled")
		}
	})
}

func TestSecurityAuditor_LogSecurityEvent(t *testing.T) {
	t.Run("EnabledAuditor", func(t *testing.T) {
		logger := &MockSecurityLogger{}
		auditor := NewSecurityAuditor(true, logger)
		
		// Should not panic when logging security event
		auditor.LogSecurityEvent("auth_success", "User authenticated successfully", map[string]interface{}{
			"user": "testuser",
			"ip":   "127.0.0.1",
		})
		
		auditor.LogSecurityEvent("auth_failure", "Authentication failed", map[string]interface{}{
			"user":   "baduser",
			"ip":     "127.0.0.1",
			"reason": "invalid_password",
		})
		
		// No assertions needed - just verify it doesn't panic
	})
	
	t.Run("DisabledAuditor", func(t *testing.T) {
		logger := &MockSecurityLogger{}
		auditor := NewSecurityAuditor(false, logger)
		
		// Should return early without logging when disabled
		auditor.LogSecurityEvent("auth_success", "User authenticated successfully", map[string]interface{}{
			"user": "testuser",
		})
		
		// No assertions needed - just verify it doesn't panic
	})
}

func TestSecurityAuditor_LogConnectionAttempt(t *testing.T) {
	t.Run("ConnectionLogging", func(t *testing.T) {
		logger := &MockSecurityLogger{}
		auditor := NewSecurityAuditor(true, logger)
		
		// Should not panic when logging connection attempt
		auditor.LogConnectionAttempt("worker1", "127.0.0.1:22", true, "")
		
		auditor.LogConnectionAttempt("worker2", "127.0.0.1:2222", false, "connection failed")
		
		// No assertions needed - just verify it doesn't panic
	})
	
	t.Run("DisabledAuditor", func(t *testing.T) {
		logger := &MockSecurityLogger{}
		auditor := NewSecurityAuditor(false, logger)
		
		// Should return early without logging when disabled
		auditor.LogConnectionAttempt("worker1", "127.0.0.1:22", true, "")
		
		// No assertions needed - just verify it doesn't panic
	})
}

func TestSecurityAuditor_LogAuthAttempt(t *testing.T) {
	t.Run("AuthLogging", func(t *testing.T) {
		logger := &MockSecurityLogger{}
		auditor := NewSecurityAuditor(true, logger)
		
		// Should not panic when logging auth attempt
		auditor.LogAuthAttempt("worker1", "testuser", "password", true)
		
		auditor.LogAuthAttempt("worker2", "baduser", "key", false)
		
		// No assertions needed - just verify it doesn't panic
	})
	
	t.Run("DisabledAuditor", func(t *testing.T) {
		logger := &MockSecurityLogger{}
		auditor := NewSecurityAuditor(false, logger)
		
		// Should return early without logging when disabled
		auditor.LogAuthAttempt("worker1", "testuser", "password", true)
		
		// No assertions needed - just verify it doesn't panic
	})
}

func TestSecurityAuditor_LogNetworkAccess(t *testing.T) {
	t.Run("NetworkLogging", func(t *testing.T) {
		logger := &MockSecurityLogger{}
		auditor := NewSecurityAuditor(true, logger)
		
		// Should not panic when logging network access
		auditor.LogNetworkAccess("127.0.0.1:22", true)
		
		auditor.LogNetworkAccess("192.168.1.100:8080", false)
		
		// No assertions needed - just verify it doesn't panic
	})
	
	t.Run("DisabledAuditor", func(t *testing.T) {
		logger := &MockSecurityLogger{}
		auditor := NewSecurityAuditor(false, logger)
		
		// Should return early without logging when disabled
		auditor.LogNetworkAccess("127.0.0.1:22", true)
		
		// No assertions needed - just verify it doesn't panic
	})
}
