package chat

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ConnectReq is the request for connecting to WebSocket chat
type ConnectReq struct {
	g.Meta `path:"/ws/chat" method:"get" tags:"Chat" summary:"Connect to WebSocket chat" auth:"true"`
	RoomId uint `v:"required|min:1" dc:"Room ID to join"`
}

// MessageReq represents a message sent from the client
type MessageReq struct {
	Type    int    `json:"type"    description:"Message type: 0-text, 1-image, 2-file"`
	Content string `json:"content" description:"Message content"`
	RoomId  uint   `json:"roomId"  description:"Room ID"`
}

// MessageRes represents a message sent to the client
type MessageRes struct {
	Id        uint   `json:"id"        description:"Message ID"`
	Type      int    `json:"type"      description:"Message type: 0-text, 1-image, 2-file, 3-system"`
	Content   string `json:"content"   description:"Message content"`
	RoomId    uint   `json:"roomId"    description:"Room ID"`
	UserId    uint   `json:"userId"    description:"User ID who sent the message"`
	Username  string `json:"username"  description:"Username who sent the message"`
	Nickname  string `json:"nickname"  description:"Nickname who sent the message"`
	Avatar    string `json:"avatar"    description:"Avatar URL of the sender"`
	Timestamp string `json:"timestamp" description:"Message timestamp"`
}

// RoomMembersReq is the request for getting room members
type RoomMembersReq struct {
	g.Meta `path:"/chat/room/{id}/members" method:"get" tags:"Chat" summary:"Get room members" auth:"true"`
	Id     uint `v:"required|min:1" in:"path" dc:"Room ID"`
}

// RoomMembersRes is the response for room members
type RoomMembersRes struct {
	Members []Member `json:"members" dc:"List of room members"`
	Total   int      `json:"total"   dc:"Total number of members in the room"`
}

// Member represents a member in a chat room
type Member struct {
	Id       uint   `json:"id"       dc:"User ID"`
	Username string `json:"username" dc:"Username"`
	Nickname string `json:"nickname" dc:"Nickname"`
	Avatar   string `json:"avatar"   dc:"Avatar URL"`
	Status   int    `json:"status"   dc:"User status: 0-offline, 1-online"`
}

// HistoryReq is the request for getting chat history
type HistoryReq struct {
	g.Meta `path:"/chat/history/{roomId}" method:"get" tags:"Chat" summary:"Get chat history" auth:"true"`
	RoomId uint `v:"required|min:1" in:"path" dc:"Room ID"`
	Page   int  `d:"1"  v:"min:1"    dc:"Page number, starting from 1"`
	Size   int  `d:"20" v:"max:100"  dc:"Page size, maximum 100"`
}

// HistoryRes is the response for chat history
type HistoryRes struct {
	Messages []MessageRes `json:"messages" dc:"List of messages"`
	Total    int          `json:"total"    dc:"Total number of messages"`
	Page     int          `json:"page"     dc:"Current page number"`
	Size     int          `json:"size"     dc:"Page size"`
}
