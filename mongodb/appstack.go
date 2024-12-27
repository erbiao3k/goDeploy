package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goDeploy/config"
	"time"
)

const appstackCollection = "appstack"

//// AppStack 应用堆栈信息表结构体
//type AppStack struct {
//	AppName       string `bson:"app_name" json:"app_name" form:"app_name"`
//	AppDesc       string `bson:"app_desc" json:"app_desc" form:"app_desc"`
//	DeployWebhook string `bson:"deploy_webhook" json:"deploy_webhook" form:"deploy_webhook"`
//	EventWebhook  string `bson:"event_webhook" json:"event_webhook" form:"event_webhook"`
//	EventToken    string `bson:"event_token" json:"event_token" form:"event_token"`
//}

// AppStack 应用堆栈信息表结构体
type AppStack struct {
	AppName       string `bson:"app_name" form:"app_name"`
	AppDesc       string `bson:"app_desc" form:"app_desc"`
	DeployWebhook string `bson:"deploy_webhook" form:"deploy_webhook"`
	EventWebhook  string `bson:"event_webhook" form:"event_webhook"`
	EventToken    string `bson:"event_token"  form:"event_token"`
}

// AppStackManager 应用堆栈信息管理器
type AppStackManager struct {
	BaseDBManager
}

// NewAppStackManager 初始化应用堆栈信息管理器
func NewAppStackManager(uri string) *AppStackManager {
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

	return &AppStackManager{
		BaseDBManager: BaseDBManager{Client: db.Client, DB: db.DB},
	}
}

// AddAppStack 添加应用堆栈信息
func (m *AppStackManager) AddAppStack(appStack AppStack) (primitive.ObjectID, error) {
	coll := m.DB.Collection(appstackCollection)
	result, err := coll.InsertOne(context.TODO(), appStack)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// DeleteAppStack 删除应用堆栈信息
func (m *AppStackManager) DeleteAppStack(appName string) (int64, error) {
	coll := m.DB.Collection(appstackCollection)
	result, err := coll.DeleteOne(context.TODO(), bson.M{"app_name": appName})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// UpdateAppStack 修改应用堆栈信息
func (m *AppStackManager) UpdateAppStack(appName string, update bson.M) (int64, error) {
	coll := m.DB.Collection(appstackCollection)
	result, err := coll.UpdateOne(context.TODO(), bson.M{"app_name": appName}, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// ListAppStacks 列出所有应用堆栈信息
func (m *AppStackManager) ListAppStacks() ([]AppStack, error) {
	coll := m.DB.Collection(appstackCollection)
	var appStacks []AppStack
	cur, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var appStack AppStack
		err := cur.Decode(&appStack)
		if err != nil {
			return nil, err
		}
		appStacks = append(appStacks, appStack)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return appStacks, nil
}

// GetAppStackByName 根据应用名称查询应用堆栈信息
func (m *AppStackManager) GetAppStackByName(appName string) (*AppStack, error) {
	coll := m.DB.Collection(appstackCollection)
	var appStack AppStack
	err := coll.FindOne(context.TODO(), bson.M{"app_name": appName}).Decode(&appStack)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 如果没有找到文档，返回nil和nil错误
			return nil, mongo.ErrNoDocuments
		}
		//如果有其他错误，返回错误
		return nil, err
	}
	return &appStack, nil
}

// EnsureIndexes 确保所需的索引被创建
func (m *AppStackManager) EnsureIndexes() error {
	SkipDuplicateKeyError := func(err error) error {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		} else {
			return err
		}
	}

	// 为app_name字段创建唯一索引
	_, err := m.DB.Collection(appstackCollection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{"app_name", 1}},
		Options: options.Index().SetUnique(true),
	})

	err = SkipDuplicateKeyError(err)
	if err != nil {
		return err
	}

	// 为deploy_webhook字段创建唯一索引
	_, err = m.DB.Collection(appstackCollection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{"deploy_webhook", 1}},
		Options: options.Index().SetUnique(true),
	})

	err = SkipDuplicateKeyError(err)
	if err != nil {
		return err
	}

	// 为event_webhook字段创建唯一索引
	_, err = m.DB.Collection(appstackCollection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{"event_webhook", 1}},
		Options: options.Index().SetUnique(true),
	})

	err = SkipDuplicateKeyError(err)
	if err != nil {
		return err
	}

	return nil
}

// Close 关闭MongoDB连接
func (m *AppStackManager) Close() error {
	return m.Client.Disconnect(context.TODO())
}
