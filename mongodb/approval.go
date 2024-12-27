package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goDeploy/config"
	"time"
)

const approvalCollection = "approval"

// Approval 审批表结构体
type Approval struct {
	SpName          string   `bson:"sp_name"`
	SpNo            string   `bson:"sp_no"`
	ApplyTime       int64    `bson:"apply_time"`
	Applyer         string   `bson:"applyer"`
	UserId          string   `bson:"user_id"`
	Tel             string   `bson:"tel"`
	DeployTime      int64    `bson:"deploy_time"`
	AppName         string   `bson:"app_name"`
	DeployContent   string   `bson:"deploy_content"`
	SqlContent      [][]byte `bson:"sql_content"`
	DeployStatus    string   `bson:"deploy_status"`
	FailedLog       string   `bson:"failed_log"`
	DeployTotalTime int      `bson:"deploy_total_time"`
}

// ApprovalManager 审批表管理器
type ApprovalManager struct {
	BaseDBManager
}

// NewApprovalManager 初始化审批表管理器
func NewApprovalManager(uri string) *ApprovalManager {
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

	return &ApprovalManager{
		BaseDBManager: BaseDBManager{Client: db.Client, DB: db.DB},
	}
}

// AddApproval 添加审批数据
func (am *ApprovalManager) AddApproval(approval Approval) (primitive.ObjectID, error) {
	coll := am.DB.Collection(approvalCollection)
	result, err := coll.InsertOne(context.TODO(), approval)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// DeleteApproval 删除审批数据
func (am *ApprovalManager) DeleteApproval(spNo string) (int64, error) {
	coll := am.DB.Collection(approvalCollection)
	result, err := coll.DeleteOne(context.TODO(), bson.M{"sp_no": spNo})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// UpdateApprovalMany 修改审批数据
func (am *ApprovalManager) UpdateApprovalMany(filter bson.M, update bson.M) (int64, error) {
	coll := am.DB.Collection(approvalCollection)
	result, err := coll.UpdateMany(context.TODO(), filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// GetApprovals 查询指定状态的审批数据
func (am *ApprovalManager) GetApprovals(filter bson.M) ([]Approval, error) {
	coll := am.DB.Collection(approvalCollection)
	var approvals []Approval
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var approval Approval
		err := cur.Decode(&approval)
		if err != nil {
			return nil, err
		}
		approvals = append(approvals, approval)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return approvals, nil
}

// EnsureIndexes 确保所需的索引被创建
func (am *ApprovalManager) EnsureIndexes() error {
	// 为sp_name, sp_no和app_name字段创建复合唯一索引
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"sp_name", 1}, {"sp_no", 1}, {"app_name", 1}},
		Options: options.Index().SetUnique(true),
	}
	coll := am.DB.Collection(approvalCollection)
	_, err := coll.Indexes().CreateOne(context.TODO(), indexModel)

	if err != nil {
		return err
	}
	return nil
}

// Close 关闭MongoDB连接
func (am *ApprovalManager) Close() error {
	return am.Client.Disconnect(context.TODO())
}
