package database

import (
	"context"
	"time"

	"zhku-oj/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDB 创建MongoDB连接
func NewMongoDB(cfg config.MongoDBConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(30 * time.Minute)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

// GetDatabase 获取数据库实例
func GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}

// GetCollection 获取集合实例
func GetCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}
