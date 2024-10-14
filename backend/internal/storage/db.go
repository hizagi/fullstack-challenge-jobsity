package storage

import (
	"context"
	"fmt"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const database = "taskDatabase"

type MongoStorage struct {
	client *mongo.Client
}

func NewMongoStorage(ctx context.Context, dbConfig config.DB) (*MongoStorage, error) {
	// Set client options
	clientOptions := options.Client().
		SetConnectTimeout(dbConfig.Timeout).
		SetMaxPoolSize(dbConfig.MaxPoolSize).
		ApplyURI(dbConfig.URI)

	if dbConfig.User != "" && dbConfig.Password != "" {
		clientOptions.SetAuth(options.Credential{Username: dbConfig.User, Password: dbConfig.Password})
	}

	if err := clientOptions.Validate(); err != nil {
		return nil, fmt.Errorf("invalid mongo configuration: %w", err)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	return &MongoStorage{
		client: client,
	}, nil
}

func (m *MongoStorage) GetDatabase() *mongo.Database {
	return m.client.Database(database)
}

func (m *MongoStorage) Disconnect(ctx context.Context) error {
	if err := m.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect mongo client: %w", err)
	}
	return nil
}
