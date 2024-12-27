package router

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"goDeploy/mongodb"
	p "goDeploy/public"
	"goDeploy/run"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// 定义和X-AppStack-Action的常量
const (
	// X-AppStack-Event 常量
	XAppStackEventReleaseStageExecution = "ReleaseStageExecution"
	XAppStackEventChangeRequest         = "ChangeRequest"
	XAppStackEventVariableGroup         = "VariableGroup"
	XAppStackEventAppOrchestration      = "AppOrchestration"
	XAppStackEventChangeOrder           = "ChangeOrder"
	XAppStackEventEnv                   = "Env"
	XAppStackEventApp                   = "App"

	// X-AppStack-Action 常量
	XAppStackActionStatusUpdate = "StatusUpdate"
	XAppStackActionCreate       = "Create"

	ReleaseStageRunning  = "RUNNING"
	ReleaseStageCanceled = "CANCELED"
	ReleaseStageSuccess  = "SUCCESS"
	ReleaseStageFailed   = "FAILED"
)

// WebhookHeader 映射Webhook请求的头部信息
type WebhookHeader struct {
	XAppStackEvent  string `header:"X-AppStack-Event"`
	XAppStackAction string `header:"X-AppStack-Action"`
	XAppStackToken  string `header:"X-AppStack-Token"`
	XAppStackApp    string `header:"X-AppStack-App"`
}

// WebhookBody 映射Webhook请求的主体信息
type WebhookBody struct {
	ID               string           `json:"id"`
	User             User             `json:"user"`
	OrgID            string           `json:"orgId"`
	Time             int64            `json:"time"`
	ObjectKind       string           `json:"objectKind"`
	ObjectAttributes ObjectAttributes `json:"objectAttributes"`
}

// User 映射Webhook请求中的用户信息
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"userName"`
	StaffId  string `json:"staffId"`
}

// ObjectAttributes 映射Webhook请求中的对象属性
type ObjectAttributes struct {
	SN               string  `json:"sn"`
	ReleaseStageSN   string  `json:"releaseStageSn"`
	EngineType       string  `json:"engineType"`
	EngineInstanceID string  `json:"engineInstanceId"`
	State            string  `json:"state"`
	Context          Context `json:"context"`
	ID               string  `json:"id"`
}

// Context 映射Webhook请求中的上下文信息
type Context struct {
	Version            string              `json:"version"`
	ReleaseStageLabels []ReleaseStageLabel `json:"releaseStageLabels"`
	TriggerMode        string              `json:"triggerMode"`
	StartTime          int64               `json:"startTime"`
	EndTime            int64               `json:"endTime"`
}

// ReleaseStageLabel 映射Webhook请求中的发布阶段标签
type ReleaseStageLabel struct {
	Namespace    string            `json:"namespace"`
	Name         string            `json:"name"`
	Value        string            `json:"value"`
	DisplayName  string            `json:"displayName"`
	DisplayValue string            `json:"displayValue"`
	ExtraMap     map[string]string `json:"extraMap"`
}

type Notification struct {
	TestName       string `json:"testName"`
	TestResult     string `json:"testResult"`
	ErrorMessage   string `json:"errorMessage"`
	ManualTestLink string `json:"manualTestLink"`
	TestDetails    string `json:"testDetails"`
}

var appStacksTemplate = template.Must(template.New("appStacks").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>App Stacks</title>
    <style>
        .container {
            display: flex;
            max-width: 1200px;
            margin: 0 auto;
        }
        .form-container, .list-container {
            padding: 20px;
        }
        .form-container {
            border-right: 1px solid #ccc;
        }
        @media (max-width: 768px) {
            .container {
                flex-direction: column;
            }
            .form-container {
                border-right: none;
                border-bottom: 1px solid #ccc;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="form-container">
            <h2>Add App Stack</h2>
            <form action="/appstacks" method="post">
                <label for="app_name">App Name:</label><br>
                <input type="text" id="app_name" name="app_name" required><br>
                <label for="app_desc">App Desc:</label><br>
                <input type="text" id="app_desc" name="app_desc" required><br>
                <label for="deploy_webhook">Deploy Webhook:</label><br>
                <input type="text" id="deploy_webhook" name="deploy_webhook" required><br>
                <label for="event_webhook">Event Webhook:</label><br>
                <input type="text" id="event_webhook" name="event_webhook" required><br>
                <label for="event_token">Event Token:</label><br>
                <input type="text" id="event_token" name="event_token" required><br>
                <button type="submit">Add App Stack</button>
            </form>
        </div>
        <div class="list-container">
            <h2>App Stacks List</h2>
            <ul>
                {{range .}}
                    <li>
                        <strong>App Name:</strong> {{.AppName}}<br>
                        <strong>App Desc:</strong> {{.AppDesc}}<br>
                        <strong>Deploy Webhook:</strong> {{.DeployWebhook}}<br>
                        <strong>Event Webhook:</strong> {{.EventWebhook}}<br>
                        <strong>Event Token:</strong> {{.EventToken}}<br><br>

                    </li>
                {{end}}
            </ul>
        </div>
    </div>
</body>
</html>
`))

// handleAppStacksGET 处理GET请求，显示AppStack列表
func handleAppStacksGET(c *gin.Context) {
	appStacks, err := mongodb.MyAppStackManager.ListAppStacks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list app stack err" + err.Error()})
		return
	}
	err = appStacksTemplate.Execute(c.Writer, appStacks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get app stack web err：" + err.Error()})
		return
	}
}

func handleAppStacksPOST(c *gin.Context) {
	var form mongodb.AppStack

	// 使用 ShouldBind 绑定表单数据
	if err := c.ShouldBind(&form); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := mongodb.MyAppStackManager.AddAppStack(form)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusSeeOther, "/appstacks")
}

// handleEventWebhookPOST 处理POST请求，接收研发流程阶段运行状态更新信息
func handleEventWebhookPOST(c *gin.Context) {
	var header WebhookHeader
	var body WebhookBody

	appName := c.Param("appname")

	appStack, err := mongodb.MyAppStackManager.GetAppStackByName(appName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no such appstack,err: " + err.Error()})
		return
	}

	// 绑定header
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println("bind app stack webhook header err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if header.XAppStackToken != appStack.EventToken {
		log.Println("Invalid app stack webhook XAppStackToken")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid header"})
		return
	}

	log.Printf("Headers: %+v\n", header)
	// 打印原始请求数据
	requestBody, _ := c.GetRawData()
	log.Printf("Request Body: %s\n", requestBody)

	// 重置请求体，以便 ShouldBindJSON 可以再次读取
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	// 绑定body
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("bind app stack webhook body err：", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if header.XAppStackEvent == XAppStackEventReleaseStageExecution {
		var appChangeStats, autoTestingStatus, str string
		switch body.ObjectAttributes.State {
		case ReleaseStageRunning:
			appChangeStats = run.DeployStatusDeploying
			autoTestingStatus = run.DeployStatusTesting
			str = "正在部署！！！"
		case ReleaseStageCanceled:
			appChangeStats = run.DeployStatusDeployCanceled
			autoTestingStatus = run.DeployStatusTestingCanceled
			str = "部署已取消！！！"
		case ReleaseStageSuccess:
			appChangeStats = run.DeployStatusDeploySuccess
			autoTestingStatus = run.DeployStatusTestingSuccess
			str = "部署成功！！！"
		case ReleaseStageFailed:
			appChangeStats = run.DeployStatusDeployFailed
			autoTestingStatus = run.DeployStatusTestingFailed
			str = "部署失败！！！"
		default:
			log.Println("unknow body.ObjectAttributes.State：", body.ObjectAttributes.State)
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknow body.ObjectAttributes.State"})
			return
		}
		if appName != "xod-autotest" {
			filter := bson.M{"app_name": appName, "deploy_status": run.DeployStatusDeploying}
			updateStatus := bson.M{"deploy_status": appChangeStats}
			_, err := mongodb.MyApprovalManager.UpdateApprovalMany(filter, updateStatus)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "update many approval err: " + err.Error()})
				return
			}

			approvals, err := mongodb.MyApprovalManager.GetApprovals(filter)
			if err != nil {
				return
			}
			if len(approvals) == 0 {
				data := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "【应用未审批部署】\n应用名称：%s\n部署状态：%s","mentioned_list":["%s"]}}`, appName, str, "@all")
				err = p.Send(data)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "post msg err: " + err.Error()})
					return
				}
			}
		}
		_, err := mongodb.MyApprovalManager.UpdateApprovalMany(bson.M{"deploy_status": run.DeployStatusTesting}, bson.M{"deploy_status": autoTestingStatus})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update many approval err: " + err.Error()})
			return
		}
	}
}
