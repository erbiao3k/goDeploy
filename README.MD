# goDeploy 说明文档

## 简介
goDeploy 是一个基于企业微信审批流的自动化部署工具。

## 功能特性
- **自动化部署**：通过企业微信审批流程实现应用的自动化部署。
- **集成企业微信**：与企业微信深度集成，确保部署流程的安全性和可控性。

## 环境变量配置说明

### 1. WECOM_ROBOT_ADDR
- **环境变量名称**: `WECOM_ROBOT_ADDR`
- **示例**: `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx,https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=yyy`
- **用法**: 用于存储企业微信机器人地址，支持多个地址以逗号分隔。

### 2. MONGODB_URI
- **环境变量名称**: `MONGODB_URI`
- **示例**: `mongodb://localhost:27017`
- **用法**: 用于连接 MongoDB 数据库的 URI 地址。

### 3. PROVIDER_SECRET
- **环境变量名称**: `PROVIDER_SECRET`
- **示例**: `mysecretkey123`
- **用法**: 提供者密钥，用于身份验证或其他安全相关的配置。

### 4. DEPLOY_TEMPLATE_ID
- **环境变量名称**: `DEPLOY_TEMPLATE_ID`
- **示例**: `template-001`
- **用法**: 部署模板 ID，用于指定部署时使用的模板。

### 5. CORP_ID
- **环境变量名称**: `CORP_ID`
- **示例**: `ww1234567890abcde3`
- **用法**: 企业微信的企业 ID。

### 6. AUTO_NODE_LIST
- **环境变量名称**: `AUTO_NODE_LIST`
- **示例**: `node1,node2,node3`
- **用法**: 自动节点列表，支持多个节点以逗号分隔。此变量还用于生成 `AutoNodePeople` 变量。

以上环境变量需要在运行应用程序之前设置好，以便程序能够正确读取并使用这些配置。
