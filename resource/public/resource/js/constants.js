/**
 * 用户状态常量定义
 * 与后端 consts.UserStatus 对应
 */
export const UserStatus = {
    /** 用户离线状态 */
    OFFLINE: 0,
    /** 用户在线状态 */
    ONLINE: 1
};

/**
 * 消息类型常量定义
 * 与后端 consts.MessageType 对应
 */
export const MessageType = {
    /** 文本消息 */
    TEXT: 0,
    /** 图片消息 */
    IMAGE: 1,
    /** 文件消息 */
    FILE: 2,
    /** 系统消息 */
    SYSTEM: 3
};

/**
 * WebSocket消息类型常量定义
 * 与后端 consts.WsMsgType 对应
 */
export const WsMessageType = {
    /** 文本消息 */
    TEXT: 1,
    /** 用户加入房间 */
    JOIN: 2,
    /** 用户离开房间 */
    LEAVE: 3,
    /** 用户列表更新 */
    USER_LIST: 4,
    /** 错误消息 */
    ERROR: 5,
    /** 系统通知 */
    NOTIFICATION: 6
};

/**
 * 默认值常量定义
 */
export const Defaults = {
    /** 默认头像路径 */
    AVATAR: '/resource/image/avatar/default.png'
};

/**
 * 错误消息常量定义
 * 与后端错误消息对应
 */
export const Errors = {
    /** 用户未加入聊天室错误 */
    NOT_IN_ROOM: 'User is not in the chat room'
};