# GoFrame 聊天室

**实验性的，用copilot生成的代码，bug较多。**

基于 GoFrame 框架开发的 WebSocket 聊天室应用。

## 功能特点

- 💬 实时聊天

  - 支持文本、图片、文件等多种消息类型
  - 消息历史记录查看
  - 在线用户状态实时更新
  - 系统消息通知（加入/离开等事件）

- 👥 聊天室管理

  - 创建公开/私密聊天室
  - 聊天室列表显示在线人数
  - 支持加入/离开聊天室
  - 房主可删除聊天室

- 👤 用户系统

  - 用户注册/登录
  - JWT 身份认证
  - 个人信息管理
  - 在线状态显示

- 🛡️ 安全特性

  - WebSocket 连接鉴权
  - 私密聊天室访问控制
  - 聊天室创建者权限管理

- 💻 技术特点
  - 基于 GoFrame v2 框架开发
  - SQLite 数据存储
  - WebSocket 实时通信
  - 前端采用 Bootstrap 5 UI 框架
  - 支持断线重连机制

## 项目开发步骤

### 1. 环境准备

- 安装 Go 环境 (1.16+)
- 安装 GoFrame CLI 工具
  ```
  go install github.com/gogf/gf/v2/cmd/gf@latest
  ```

### 2. 项目配置

- 配置文件位于 `manifest/config/config.yaml`
- 默认使用 SQLite 数据库，无需额外配置
- 支持自定义服务端口、日志级别等

### 3. 项目运行

```bash
# 启动服务
gf run main.go

# 或直接运行编译后的程序
go build -o main.exe
./main.exe
```

### 4. 访问服务

- 打开浏览器访问: http://localhost:8000
- 注册新用户并登录
- 创建或加入聊天室开始聊天

## 项目结构

```
.
├── api           # API 接口定义
├── internal      # 内部逻辑实现
│   ├── cmd      # 启动入口
│   ├── consts   # 常量定义
│   ├── controller# 控制器
│   ├── dao      # 数据访问
│   ├── model    # 数据模型
│   └── service  # 业务逻辑
├── manifest     # 配置文件
└── resource     # 静态资源
    └── public   # Web 前端文件
```

## 技术栈

### 后端

- GoFrame v2：基础框架
- SQLite：数据存储
- JWT：用户认证
- WebSocket：实时通信

### 前端

- Bootstrap 5：UI 框架
- Font Awesome：图标库
- WebSocket：实时通信
- 原生 JavaScript：业务逻辑

## 主要功能实现

### WebSocket 通信

- 基于 gorilla/websocket 实现的实时通信
- 支持自动重连（最多 5 次尝试，间隔 3 秒）
- 支持房间广播和用户列表同步
- 在线状态实时更新

### 消息处理

- 支持多种消息类型：文本、图片、文件
- 文件上传使用 Base64 编码传输
- 消息持久化存储，支持历史记录查询
- 系统消息自动通知（用户加入/退出等）

### 用户认证

- 基于 JWT 的用户认证机制
- Token 过期自动处理
- WebSocket 连接携带 Token 认证
- 用户信息安全存储

### 权限控制

- 基于角色的权限控制：房主/普通用户
- 私密房间访问控制
- 房主特殊权限：删除房间
- 用户操作验证

## 部署说明

### Docker 部署

```bash
# 构建镜像
docker build -t goframe-chat .

# 运行容器
docker run -d -p 8000:8000 goframe-chat
```

### K8s 部署

- 部署文件位于 `manifest/deploy/kustomize/`
- 支持不同环境配置（base/overlays）
- 包含 Deployment、Service、ConfigMap 等资源

## 贡献指南

1. Fork 本项目
2. 创建新的功能分支
3. 提交你的更改
4. 发起 Pull Request

## 开源许可

MIT License
