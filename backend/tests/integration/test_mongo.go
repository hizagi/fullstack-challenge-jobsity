//go:build integration

package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartMongoContainer() (*config.DB, func(context.Context), error) {
	ctx := context.Background()

	// Set up the MongoDB container request
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest", // Pull the latest MongoDB image
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(30 * time.Second),
	}

	// Start the MongoDB container
	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start mongo container: %w", err)
	}

	// Get the host and port for the MongoDB container
	host, err := mongoContainer.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get mongo container host: %w", err)
	}

	port, err := mongoContainer.MappedPort(ctx, "27017")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get mongo container port: %w", err)
	}

	// Construct the MongoDB URI
	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	// Create a DBConfig with the MongoDB URI and default values for testing
	dbConfig := &config.DB{
		URI:         mongoURI,
		User:        "", // No authentication for local test setup
		Password:    "", // No authentication for local test setup
		Database:    "testdb",
		Timeout:     10 * time.Second,
		MaxPoolSize: 100,
	}

	return dbConfig, func(closeCtx context.Context) {
		mongoContainer.Terminate(closeCtx) // Close the container when tests are done
	}, nil
}
