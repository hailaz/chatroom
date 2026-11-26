# 消息类型支持清单

本文档列出聊天室当前版本实际支持的消息类型、WebSocket 消息格式与相关限制，便于前端、后端和测试对齐。

- 版本基线：`master` 分支 `e1e1b2f`（2026-11-26）
- 相关代码：
  - 类型常量：[internal/consts/consts.go](../internal/consts/consts.go)
  - 前端常量：[resource/public/resource/js/constants.js](../resource/public/resource/js/constants.js)
  - WebSocket 收发：[internal/service/websocket.go](../internal/service/websocket.go)
  - 消息渲染：[resource/public/resource/js/ui.js](../resource/public/resource/js/ui.js)

## 快速结论

| 分类 | 当前版本内容 |
| --- | --- |
| 用户可发送的聊天消息 | 文本、图片、文件 |
| 服务端自动推送 | 系统提示（加入/离开/删除聊天室）、加入通知、用户列表 |
| 未支持 | 语音、视频、贴纸/表情包、位置、名片、链接卡片、Markdown/富文本、引用、转发、撤回、已读回执等 |

> 文本内容本身支持 Unicode，可以直接输入 emoji；但 emoji 只是文本的一部分，不是独立消息类型。

## 两类 `type` 字段

聊天室有**两套**消息类型常量，但都通过 JSON 中同一个 `type` 字段传输，数值空间存在重叠，阅读和扩展时需要区分：

| 分类 | 用途 | 后端 | 前端 |
| --- | --- | --- | --- |
| 业务消息类型 | 聊天内容，通常会写入数据库并出现在历史接口中（系统消息除外） | `consts.MessageType*` | `MessageType` |
| WebSocket 控制消息类型 | 连接状态、用户列表等服务端推送 | `consts.WsMsgType*` | `WsMessageType` |

## 业务消息类型

对应用户在聊天室中看到的消息，定义于 `internal/consts/consts.go`。

| `type` | 名称 | 常量 | 发送方 | `content` 格式 | 写入数据库 | 前端展示 |
| --- | --- | --- | --- | --- | --- | --- |
| 0 | 文本 | `MessageTypeText` | 客户端 | 纯文本 | 是 | 消息气泡，HTML 转义后显示 |
| 1 | 图片 | `MessageTypeImage` | 客户端 | Base64 Data URL：`data:image/...;base64,...` | 是 | 消息气泡内显示图片 |
| 2 | 文件 | `MessageTypeFile` | 客户端 | JSON 字符串，见下文 | 是 | 文件名 + 文件大小卡片 |
| 3 | 系统 | `MessageTypeSystem` | 服务端 | 纯文本 | 否 | 居中的灰色提示文字 |

### 0 - 文本消息

文本消息由客户端通过 WebSocket 发送，`content` 为原始文本。服务端补齐用户信息和时间戳后广播给房间内所有连接，并写入 `messages` 表。

```json
{
  "type": 0,
  "content": "你好",
  "roomId": 1
}
```

前端渲染前会调用 `escapeHtml`，文本中的 HTML 标签不会生效。

### 1 - 图片消息

图片先由浏览器读取为 Base64 Data URL，再作为 `content` 发送。服务端不解析图片，只把整段字符串持久化并广播。

```json
{
  "type": 1,
  "content": "data:image/png;base64,iVBORw0KGgo...",
  "roomId": 1
}
```

图片选择框只接受 `image/*`，但这只是浏览器侧提示，服务端没有 MIME 类型或文件大小校验。

### 2 - 文件消息

文件消息的 `content` 是一个 **JSON 字符串**（注意是字符串而不是嵌套对象），包含文件名、大小、MIME 类型和 Base64 Data URL：

```json
{
  "type": 2,
  "content": "{\"name\":\"report.pdf\",\"size\":12345,\"type\":\"application/pdf\",\"data\":\"data:application/pdf;base64,JVBERi0...\"}",
  "roomId": 1
}
```

前端解析后展示文件名和格式化后的大小。文件同样以 Base64 文本形式存入 SQLite，不落在独立文件存储中。

### 3 - 系统消息

正常流程中，系统消息只由服务端产生，`content` 为纯文本。当前版本会在以下事件广播系统消息：

| 触发事件 | 文案 | 产生位置 |
| --- | --- | --- |
| 加入聊天室 | `<昵称> 加入了聊天室` | `ChatRoomService.Join` |
| 离开聊天室 | `<昵称> 离开了聊天室` | `ChatRoomService.Leave` |
| 删除聊天室 | `聊天室已被管理员删除` | `ChatRoomService.Delete` |

系统消息：

- 只发送给事件发生时仍在线的房间连接；
- 不会写入数据库，因此不会出现在历史消息接口中；
- 广播时不带 `userId`、`username`、`nickname`、`avatar` 字段。

```json
{
  "type": 3,
  "content": "Alice 加入了聊天室",
  "timestamp": "2026-11-26T10:00:00Z"
}
```

## WebSocket 控制消息类型

定义于 `internal/consts/consts.go` 的 `WsMsgType*`，用于连接状态和用户列表同步，不属于聊天内容。

| `type` | 常量 | 方向 | 当前版本是否发送 | 说明 |
| --- | --- | --- | --- | --- |
| 1 | `WsMsgTypeText` | 服务端 → 客户端 | 否 | 仅定义，未使用；文本聊天消息使用业务类型 0 |
| 2 | `WsMsgTypeJoin` | 服务端 → 客户端 | 是 | WebSocket 连接建立时广播 `<昵称> joined the room` |
| 3 | `WsMsgTypeLeave` | 服务端 → 客户端 | 否 | 仅定义；离开提示实际使用业务类型 3 |
| 4 | `WsMsgTypeUserList` | 服务端 → 客户端 | 是 | 房间用户列表更新，数据放在 `data` 字段 |
| 5 | `WsMsgTypeError` | 服务端 → 客户端 | 否 | 仅定义，未使用 |
| 6 | `WsMsgTypeNotification` | 服务端 → 客户端 | 否 | 仅定义，未使用 |

### 加入通知（type = 2）

用户建立 WebSocket 连接后，服务端会向房间内所有连接广播：

```json
{
  "type": 2,
  "content": "Alice joined the room",
  "timestamp": "2026-11-26T10:00:00Z",
  "userId": 1,
  "username": "alice",
  "nickname": "Alice",
  "avatar": "/resource/image/avatar/default.png"
}
```

该通知不写入数据库，且文案为英文；中文的“加入了聊天室”来自 HTTP 加入接口产生的系统消息（type = 3）。

### 用户列表（type = 4）

```json
{
  "type": 4,
  "data": [
    {
      "id": 1,
      "username": "alice",
      "nickname": "Alice",
      "avatar": "/resource/image/avatar/default.png",
      "status": 1
    }
  ],
  "timestamp": "2026-11-26T10:00:00Z"
}
```

`status` 取值：`0` 离线，`1` 在线。前端收到后会刷新用户列表，并重新加载聊天室列表以更新在线人数。

## 收发格式

### 客户端发送

客户端通过 WebSocket 发送聊天消息：

```json
{
  "type": 0,
  "content": "hello",
  "roomId": 1
}
```

- `roomId` 仅由客户端携带，服务端不按消息体中的 `roomId` 路由，而是使用 WebSocket 连接建立时绑定的房间。
- 服务端会覆盖客户端传入的 `userId`、`username`、`nickname`、`avatar`、`timestamp`；`data` 字段会原样保留。

### 服务端广播

服务端会在原消息上补齐发送者信息与 RFC3339 时间戳，然后广播给房间内所有连接：

```json
{
  "type": 0,
  "content": "hello",
  "timestamp": "2026-11-26T10:00:00Z",
  "userId": 1,
  "username": "alice",
  "nickname": "Alice",
  "avatar": "/resource/image/avatar/default.png"
}
```

### 历史消息接口

- 接口：`GET /api/chat/history/{roomId}?page=1&size=20`
- 需要登录，且当前用户必须在房间中。
- 只返回写入数据库的消息；正常客户端产生的消息类型为 `0`、`1`、`2`，系统消息 `3` 不会出现在历史记录中。
- 默认每页 20 条，最大 100 条；前端加入房间时请求 50 条。
- 消息按时间倒序返回，前端展示前会反转为正序。

## 限制与已知问题

1. **两套类型编号重叠**：业务消息和控制消息共用 `type` 字段，数值含义冲突。当前实际影响是：服务端在 WebSocket 连接建立时发送的加入通知 `type = 2`，会被前端当作文件消息解析，`JSON.parse` 失败，因此该加入通知无法正常显示。`type = 1`（图片 / `WsMsgTypeText`）和 `type = 3`（系统消息 / `WsMsgTypeLeave`）也存在同样风险，只是后两者当前未作为控制消息发送。
2. **服务端不校验类型白名单和消息来源**：客户端发送除 `3` 以外的任意 `type` 都会被写入数据库并广播；发送 `type = 3` 虽然不会入库，但仍会被广播，因此客户端可以伪造系统消息。前端只渲染 `0`、`1`、`2`、`3`，其他类型不可见。
3. **没有服务端文件校验**：图片和文件都只依赖浏览器侧读取，服务端不限制 MIME 类型、文件大小或内容；Base64 会增大约 1/3 体积。
4. **WebSocket 单帧上限 1 MB**：服务端 `SetReadLimit(1024 * 1024)`，超过后读取失败并断开连接。图片或文件编码后必须小于该上限，实际可用原始文件大小约为 750 KB 以内。
5. **消息全部存 SQLite 文本字段**：图片、文件以 Base64 存入 `messages.content`，历史记录接口也会返回完整 Base64，大文件会快速增大数据库和响应体。
6. **图片预览与文件下载点击不可用**：渲染层使用 `onclick="chatUI.xxx(...)"`，但当前页面只暴露了 `window.ChatUI` 类，没有挂载 `window.chatUI` 实例，点击会因 `chatUI` 未定义而失败；消息本身仍可正常收发和展示。
7. **转义范围有限**：文本消息会做 HTML 转义，但文件消息的 `name` 直接拼接到 `innerHTML`，存在被构造内容注入的风险。
