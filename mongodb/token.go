package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"goDeploy/config"
	"log"
	"time"
)

const tokenCollection = "token"

// Token Token表结构体
type Token struct {
	AK string `bson:"ak"`
}

// TokenManager Token表管理器
type TokenManager struct {
	BaseDBManager
}

// NewTokenManager 初始化Token表管理器
func NewTokenManager(uri string) *TokenManager {
	dbConfig := &DBConfig{
		URI:                    uri,
		ConnectTimeout:         10 * time.Second,
		ServerSelectionTimeout: 5 * time.Second,
		HeartbeatInterval:      10 * time.Second,
		MaxPoolSize:            100,
		MinPoolSize:            5,
		RetryWrites:            true,
		DatabaseName:           config.ApprovalDb,
	}

	db, err := NewDBManager(dbConfig)
	if err != nil {
		panic(err)
	}

	return &TokenManager{
		BaseDBManager: BaseDBManager{Client: db.Client, DB: db.DB},
	}
}

// UpsertToken 插入或更新Token数据，确保表中只有一条记录
func (tm *TokenManager) UpsertToken(token *Token) error {
	collection := tm.DB.Collection(tokenCollection)
	// 尝试获取当前的Token数量
	count, err := collection.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		return err
	}
	if count > 0 {
		// 如果记录数大于0，则更新
		_, err = collection.UpdateOne(context.TODO(), bson.M{}, bson.M{"$set": token})
		return err
	} else {
		// 如果记录数为0，则插入
		_, err = collection.InsertOne(context.TODO(), token)
		return err
	}
}

// GetToken 查询Token数据
func (tm *TokenManager) GetToken() (*Token, error) {
	collection := tm.DB.Collection(tokenCollection)
	var token Token
	// 由于只需要一条记录，使用FindOne进行查询
	result := collection.FindOne(context.TODO(), bson.D{})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			// 没有找到文档，可以选择返回nil或错误
			log.Println("No token document found")
			return nil, nil
		}
		// 查询过程中出现其他错误
		return nil, result.Err()
	}
	err := result.Decode(&token)
	return &token, err
}
