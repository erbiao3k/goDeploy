package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goDeploy/config"
	"time"
)

const userCollection = "user"

// UserDetails 定义用户详细信息结构体
type UserDetails struct {
	UserID       string `bson:"userid" json:"userid"`
	Name         string `bson:"name" json:"name"`
	Mobile       string `bson:"mobile" json:"mobile"`
	Email        string `bson:"email" json:"email"`
	BizMail      string `bson:"biz_mail" json:"biz_mail"`
	DirectLeader string `bson:"direct_leader" json:"direct_leader"`
	// 其他字段根据需要添加
}

// UserManager 用户信息管理器
type UserManager struct {
	BaseDBManager
}

// NewUserManager 初始化用户信息管理器
func NewUserManager(uri string) *UserManager {
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

	return &UserManager{
		BaseDBManager: BaseDBManager{Client: db.Client, DB: db.DB},
	}
}

// AddOrUpdateUser 添加或更新用户信息
func (m *UserManager) AddOrUpdateUser(user UserDetails) error {
	coll := m.DB.Collection(userCollection)
	// 使用UserID作为查询条件
	filter := bson.M{"userid": user.UserID}
	// 使用UserID和所有字段作为更新内容
	update := bson.M{"$set": user}
	// 插入或更新用户信息
	result, err := coll.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	if result.UpsertedID != nil {
		return nil
	}
	return nil
}

// ListUsers 列出所有用户信息
func (m *UserManager) ListUsers() ([]UserDetails, error) {
	coll := m.DB.Collection(userCollection)
	var users []UserDetails
	cur, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var user UserDetails
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID 根据用户ID查询用户信息
func (m *UserManager) GetUserByID(userID string) (*UserDetails, error) {
	coll := m.DB.Collection(userCollection)
	var user UserDetails
	err := coll.FindOne(context.TODO(), bson.M{"userid": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 如果没有找到文档，返回nil和nil错误
			return nil, mongo.ErrNoDocuments
		}
		// 如果有其他错误，返回错误
		return nil, err
	}
	return &user, nil
}

// GetUserByName 根据用户名称查询用户信息
func (m *UserManager) GetUserByName(Name string) (*UserDetails, error) {
	coll := m.DB.Collection(userCollection)
	var user UserDetails
	err := coll.FindOne(context.TODO(), bson.M{"name": Name}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 如果没有找到文档，返回nil和nil错误
			return nil, mongo.ErrNoDocuments
		}
		// 如果有其他错误，返回错误
		return nil, err
	}
	return &user, nil
}

// EnsureIndexes 确保所需的索引被创建
func (m *UserManager) EnsureIndexes() error {
	SkipDuplicateKeyError := func(err error) error {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		} else {
			return err
		}
	}

	// 为userid字段创建唯一索引
	_, err := m.DB.Collection(userCollection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{"userid", 1}},
		Options: options.Index().SetUnique(true),
	})

	err = SkipDuplicateKeyError(err)
	if err != nil {
		return err
	}

	return nil
}

// Close 关闭MongoDB连接
func (m *UserManager) Close() error {
	return m.Client.Disconnect(context.TODO())
}
