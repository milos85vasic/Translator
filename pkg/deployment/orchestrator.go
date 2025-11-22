package deployment

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"digital.vasic.translator/internal/config"
	"digital.vasic.translator/pkg/events"
)

// DeploymentOrchestrator manages automated deployment of distributed instances
type DeploymentOrchestrator struct {
	config      *config.Config
	eventBus    *events.EventBus
	deployer    *SSHDeployer
	discoverer  *NetworkDiscoverer
	apiLogger   *APICommunicationLogger
	logger      *log.Logger
	mu          sync.RWMutex
	initialized bool
	deployed    map[string]*DeployedInstance
}

// DeployedInstance represents a deployed instance
type DeployedInstance struct {
	ID          string
	Host        string
	Port        int
	ContainerID string
	Status      string
	Config      *DeploymentConfig
	LastSeen    time.Time
	mu          sync.RWMutex
}

// DeploymentConfig holds deployment configuration for a single instance
type DeploymentConfig struct {
	Host          string
	User          string
	SSHKeyPath    string
	DockerImage   string
	ContainerName string
	Ports         []PortMapping
	Environment   map[string]string
	Volumes       []VolumeMapping
	Networks      []string
	RestartPolicy string
	HealthCheck   *HealthCheckConfig
}

// PortMapping represents port mapping configuration
type PortMapping struct {
	HostPort      int
	ContainerPort int
	Protocol      string
}

// VolumeMapping represents volume mapping configuration
type VolumeMapping struct {
	HostPath      string
	ContainerPath string
	ReadOnly      bool
}

// HealthCheckConfig represents health check configuration
type HealthCheckConfig struct {
	Test        []string
	Interval    time.Duration
	Timeout     time.Duration
	Retries     int
	StartPeriod time.Duration
}

// NewDeploymentOrchestrator creates a new deployment orchestrator
func NewDeploymentOrchestrator(cfg *config.Config, eventBus *events.EventBus) *DeploymentOrchestrator {
	logger := log.New(os.Stdout, "[DEPLOYMENT] ", log.LstdFlags)

	deployer := NewSSHDeployer(logger)
	discoverer := NewNetworkDiscoverer(cfg, logger)

	// Initialize API communication logger
	apiLogger, err := NewAPICommunicationLogger("workers_api_communication.log")
	if err != nil {
		logger.Printf("Warning: failed to initialize API logger: %v", err)
		apiLogger = nil
	}

	return &DeploymentOrchestrator{
		config:     cfg,
		eventBus:   eventBus,
		deployer:   deployer,
		discoverer: discoverer,
		apiLogger:  apiLogger,
		logger:     logger,
		deployed:   make(map[string]*DeployedInstance),
	}
}

// DeployDistributedSystem deploys the complete distributed system
func (do *DeploymentOrchestrator) DeployDistributedSystem(ctx context.Context, deploymentPlan *DeploymentPlan) error {
	do.logger.Println("Starting automated distributed system deployment...")

	// Validate deployment plan
	if err := do.validateDeploymentPlan(deploymentPlan); err != nil {
		return fmt.Errorf("invalid deployment plan: %w", err)
	}

	// Deploy main instance first
	if err := do.deployMainInstance(ctx, deploymentPlan.Main); err != nil {
		return fmt.Errorf("failed to deploy main instance: %w", err)
	}

	// Deploy worker instances
	for i, worker := range deploymentPlan.Workers {
		if err := do.deployWorkerInstance(ctx, worker, i+1); err != nil {
			do.logger.Printf("Failed to deploy worker %d: %v", i+1, err)
			continue
		}
	}

	// Wait for all instances to be healthy
	if err := do.waitForSystemHealth(ctx, deploymentPlan); err != nil {
		return fmt.Errorf("system health check failed: %w", err)
	}

	// Initialize network discovery and broadcasting
	if err := do.initializeNetworkDiscovery(ctx); err != nil {
		return fmt.Errorf("failed to initialize network discovery: %w", err)
	}

	do.emitEvent(events.Event{
		Type:      "deployment_completed",
		SessionID: "system",
		Message:   "Distributed system deployment completed successfully",
		Data: map[string]interface{}{
			"main_instance":    deploymentPlan.Main.ContainerName,
			"worker_instances": len(deploymentPlan.Workers),
			"total_instances":  len(do.deployed),
		},
	})

	return nil
}

// deployMainInstance deploys the main coordinator instance
func (do *DeploymentOrchestrator) deployMainInstance(ctx context.Context, cfg *DeploymentConfig) error {
	do.logger.Printf("Deploying main instance to %s...", cfg.Host)

	// Find available port for main instance
	port, err := do.findAvailablePort(cfg.Host, 8443)
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}

	// Update port mapping
	cfg.Ports[0].HostPort = port

	// Deploy via SSH
	containerID, err := do.deployer.DeployInstance(ctx, cfg)
	if err != nil {
		return fmt.Errorf("SSH deployment failed: %w", err)
	}

	// Register deployed instance
	instance := &DeployedInstance{
		ID:          cfg.ContainerName,
		Host:        cfg.Host,
		Port:        port,
		ContainerID: containerID,
		Status:      "deploying",
		Config:      cfg,
		LastSeen:    time.Now(),
	}

	do.mu.Lock()
	do.deployed[cfg.ContainerName] = instance
	do.mu.Unlock()

	do.logger.Printf("Main instance deployed successfully: %s (port %d)", containerID[:12], port)
	return nil
}

// deployWorkerInstance deploys a worker instance
func (do *DeploymentOrchestrator) deployWorkerInstance(ctx context.Context, cfg *DeploymentConfig, index int) error {
	do.logger.Printf("Deploying worker instance %d to %s...", index, cfg.Host)

	// Find available port for worker
	port, err := do.findAvailablePort(cfg.Host, 8443+index)
	if err != nil {
		return fmt.Errorf("failed to find available port for worker %d: %w", index, err)
	}

	// Update port mapping
	cfg.Ports[0].HostPort = port

	// Set worker-specific environment variables
	if cfg.Environment == nil {
		cfg.Environment = make(map[string]string)
	}
	cfg.Environment["WORKER_INDEX"] = fmt.Sprintf("%d", index)
	cfg.Environment["MAIN_HOST"] = do.getMainInstanceHost()

	// Deploy via SSH
	containerID, err := do.deployer.DeployInstance(ctx, cfg)
	if err != nil {
		return fmt.Errorf("SSH deployment failed for worker %d: %w", index, err)
	}

	// Register deployed instance
	instance := &DeployedInstance{
		ID:          cfg.ContainerName,
		Host:        cfg.Host,
		Port:        port,
		ContainerID: containerID,
		Status:      "deploying",
		Config:      cfg,
		LastSeen:    time.Now(),
	}

	do.mu.Lock()
	do.deployed[cfg.ContainerName] = instance
	do.mu.Unlock()

	do.logger.Printf("Worker instance %d deployed successfully: %s (port %d)", index, containerID[:12], port)
	return nil
}

// findAvailablePort finds the first available port starting from preferredPort
func (do *DeploymentOrchestrator) findAvailablePort(host string, preferredPort int) (int, error) {
	for port := preferredPort; port < preferredPort+100; port++ {
		if do.isPortAvailable(host, port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports found in range %d-%d", preferredPort, preferredPort+99)
}

// isPortAvailable checks if a port is available on the given host
func (do *DeploymentOrchestrator) isPortAvailable(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return true // Port is available if we can't connect
	}
	conn.Close()
	return false // Port is in use
}

// waitForSystemHealth waits for all deployed instances to become healthy
func (do *DeploymentOrchestrator) waitForSystemHealth(ctx context.Context, plan *DeploymentPlan) error {
	do.logger.Println("Waiting for system health checks...")

	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for system health")
		case <-ticker.C:
			allHealthy := true

			do.mu.RLock()
			for id, instance := range do.deployed {
				healthy, err := do.checkInstanceHealth(ctx, instance)
				if err != nil {
					do.logger.Printf("Health check failed for %s: %v", id, err)
					allHealthy = false
					continue
				}

				if !healthy {
					allHealthy = false
					continue
				}

				instance.mu.Lock()
				instance.Status = "healthy"
				instance.LastSeen = time.Now()
				instance.mu.Unlock()
			}
			do.mu.RUnlock()

			if allHealthy {
				do.logger.Println("All instances are healthy!")
				return nil
			}
		}
	}
}

// checkInstanceHealth performs health check on a deployed instance
func (do *DeploymentOrchestrator) checkInstanceHealth(ctx context.Context, instance *DeployedInstance) (bool, error) {
	return do.deployer.CheckInstanceHealth(ctx, instance.Host, instance.Port)
}

// initializeNetworkDiscovery initializes network discovery and broadcasting
func (do *DeploymentOrchestrator) initializeNetworkDiscovery(ctx context.Context) error {
	do.logger.Println("Initializing network discovery and broadcasting...")

	// Start broadcasting service configurations
	if err := do.discoverer.StartBroadcasting(ctx, do.deployed); err != nil {
		return fmt.Errorf("failed to start broadcasting: %w", err)
	}

	// Start discovery listener
	if err := do.discoverer.StartDiscovery(ctx); err != nil {
		return fmt.Errorf("failed to start discovery: %w", err)
	}

	do.logger.Println("Network discovery initialized successfully")
	return nil
}

// validateDeploymentPlan validates the deployment plan
func (do *DeploymentOrchestrator) validateDeploymentPlan(plan *DeploymentPlan) error {
	if plan.Main == nil {
		return fmt.Errorf("main instance configuration is required")
	}

	if len(plan.Workers) == 0 {
		return fmt.Errorf("at least one worker instance is required")
	}

	// Validate main instance config
	if err := do.validateInstanceConfig(plan.Main); err != nil {
		return fmt.Errorf("invalid main instance config: %w", err)
	}

	// Validate worker instance configs
	for i, worker := range plan.Workers {
		if err := do.validateInstanceConfig(worker); err != nil {
			return fmt.Errorf("invalid worker %d config: %w", i+1, err)
		}
	}

	return nil
}

// validateInstanceConfig validates a single instance configuration
func (do *DeploymentOrchestrator) validateInstanceConfig(cfg *DeploymentConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("host is required")
	}

	if cfg.User == "" {
		return fmt.Errorf("user is required")
	}

	if cfg.DockerImage == "" {
		return fmt.Errorf("docker image is required")
	}

	if cfg.ContainerName == "" {
		return fmt.Errorf("container name is required")
	}

	if len(cfg.Ports) == 0 {
		return fmt.Errorf("at least one port mapping is required")
	}

	return nil
}

// getMainInstanceHost returns the host of the main instance
func (do *DeploymentOrchestrator) getMainInstanceHost() string {
	do.mu.RLock()
	defer do.mu.RUnlock()

	for _, instance := range do.deployed {
		if strings.Contains(instance.ID, "main") {
			return instance.Host
		}
	}
	return ""
}

// GetDeployedInstances returns all deployed instances
func (do *DeploymentOrchestrator) GetDeployedInstances() map[string]*DeployedInstance {
	do.mu.RLock()
	defer do.mu.RUnlock()

	result := make(map[string]*DeployedInstance)
	for k, v := range do.deployed {
		result[k] = v
	}
	return result
}

// emitEvent emits an event if event bus is available
func (do *DeploymentOrchestrator) emitEvent(event events.Event) {
	if do.eventBus != nil {
		do.eventBus.Publish(event)
	}
}

// Close shuts down the deployment orchestrator
func (do *DeploymentOrchestrator) Close() error {
	do.logger.Println("Shutting down deployment orchestrator...")

	if do.apiLogger != nil {
		do.apiLogger.Close()
	}

	if do.discoverer != nil {
		do.discoverer.Close()
	}

	if do.deployer != nil {
		do.deployer.Close()
	}

	return nil
}
