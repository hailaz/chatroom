# 待办事项

## 已实现功能

### WebSocket 核心功能
- [x] WebSocket 连接管理
  - [x] 自动重连机制（最多5次重试）
  - [x] 心跳检测
  - [x] 异常处理和错误恢复
  - [x] 实时状态监控
- [x] 消息广播系统
  - [x] 群聊消息转发
  - [x] 系统通知广播
  - [x] 用户状态变更通知
  - [x] 在线用户列表同步

### 聊天室功能
- [x] 用户系统
  - [x] 用户注册
  - [x] 用户登录
  - [x] JWT 认证
  - [x] 个人信息管理
- [x] 聊天室管理
  - [x] 创建聊天室
  - [x] 加入/离开聊天室
  - [x] 删除聊天室（仅创建者可删除）
  - [x] 聊天室列表管理
  - [x] 聊天室成员管理
- [x] 消息系统
  - [x] 文本消息发送/接收
  - [x] 图片消息支持
  - [x] 文件传输功能
  - [x] 系统消息通知
  - [x] 消息历史记录
- [x] 用户状态管理
  - [x] 在线/离线状态
  - [x] 实时在线用户列表
  - [x] 用户加入/离开提醒

### 数据存储
- [x] SQLite 数据库集成
  - [x] 用户信息存储
  - [x] 聊天室信息存储
  - [x] 消息历史记录持久化
  - [x] 用户关系存储

### 前端实现
- [x] 用户界面
  - [x] 响应式布局
  - [x] 聊天室列表展示
  - [x] 在线用户列表
  - [x] 消息区域展示
- [x] 消息展示
  - [x] 文本消息渲染
  - [x] 图片预览功能
  - [x] 文件下载功能
  - [x] 系统消息特殊样式
- [x] 实时更新
  - [x] 消息实时显示
  - [x] 用户列表实时更新
  - [x] 聊天室状态实时同步
- [x] 操作交互
  - [x] 创建聊天室
  - [x] 上传文件/图片
  - [x] 发送消息
  - [x] 退出登录

## 待实现功能
- [ ] 私聊功能
- [ ] 消息提醒（未读消息通知）
- [ ] 表情包支持
- [ ] 性能优化
- [ ] 项目打包与部署优化