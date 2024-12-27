package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// DBConfig 定义数据库配置
type DBConfig struct {
	URI                    string
	ConnectTimeout         time.Duration
	ServerSelectionTimeout time.Duration
	HeartbeatInterval      time.Duration
	MaxPoolSize            uint64
	MinPoolSize            uint64
	RetryWrites            bool
	DatabaseName           string
}

// BaseDBManager 定义通用数据库管理器结构体
type BaseDBManager struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewDBManager 初始化数据库管理器
func NewDBManager(cfg *DBConfig) (*BaseDBManager, error) {
	clientOptions := options.Client().ApplyURI(cfg.URI).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetServerSelectionTimeout(cfg.ServerSelectionTimeout).
		SetHeartbeatInterval(cfg.HeartbeatInterval).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetRetryWrites(cfg.RetryWrites)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// 尝试连接以确认客户端是否能够连接到MongoDB
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.DatabaseName)
	return &BaseDBManager{
		Client: client,
		DB:     db,
	}, nil
}
