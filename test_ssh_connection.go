package main

import (
	"context"
	"fmt"
	"time"

	"digital.vasic.translator/pkg/sshworker"
	"digital.vasic.translator/pkg/logger"
)

func main() {
	// Create logger
	log := logger.NewLogger(logger.LoggerConfig{
		Level:  logger.DEBUG,
		Format: logger.FORMAT_TEXT,
	})

	// Create worker config
	config := sshworker.SSHWorkerConfig{
		Host:              "thinker.local",
		Port:              22,
		Username:          "milosvasic",
		Password:          "WhiteSnake8587",
		RemoteDir:         "/tmp/translate-ssh",
		ConnectionTimeout: 60 * time.Second,
		CommandTimeout:    30 * time.Second,
	}

	// Create worker
	worker, err := sshworker.NewSSHWorker(config, log)
	if err != nil {
		fmt.Printf("Failed to create worker: %v\n", err)
		return
	}
	defer worker.Close()

	// Test connection
	ctx := context.Background()
	fmt.Println("Testing connection...")
	if err := worker.TestConnection(ctx); err != nil {
		fmt.Printf("Connection test failed: %v\n", err)
		return
	}
	fmt.Println("Connection test passed!")

	// Test hash command
	fmt.Println("Testing hash command...")
	hash, err := worker.GetRemoteCodebaseHash(ctx)
	if err != nil {
		fmt.Printf("Hash test failed: %v\n", err)
		return
	}
	fmt.Printf("Remote hash: %s\n", hash)
}