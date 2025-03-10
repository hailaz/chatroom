package chatroom

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CreateReq is the request for creating a new chat room
type CreateReq struct {
	g.Meta      `path:"/chatroom/create" method:"post" tags:"ChatRoom" summary:"Create a new chat room" auth:"true"`
	Name        string `v:"required|length:2,50" dc:"Room name (2-50 chars)"`
	Description string `v:"max-length:200" dc:"Room description (max 200 chars)"`
	IsPrivate   bool   `dc:"Whether the room is private (invite only)"`
}

// CreateRes is the response for creating a chat room
type CreateRes struct {
	Id          uint   `json:"id" dc:"Room ID"`
	Name        string `json:"name" dc:"Room name"`
	Description string `json:"description" dc:"Room description"`
	IsPrivate   bool   `json:"isPrivate" dc:"Whether the room is private"`
}

// ListReq is the request for listing chat rooms
type ListReq struct {
	g.Meta `path:"/chatroom/list" method:"get" tags:"ChatRoom" summary:"List available chat rooms" auth:"true"`
	Page   int `d:"1" v:"min:1" dc:"Page number, starting from 1"`
	Size   int `d:"10" v:"max:50" dc:"Page size, maximum 50"`
}

// ListRes is the response for listing chat rooms
type ListRes struct {
	List  []Room `json:"list" dc:"List of rooms"`
	Total int    `json:"total" dc:"Total number of rooms"`
	Page  int    `json:"page" dc:"Current page number"`
	Size  int    `json:"size" dc:"Page size"`
}

// Room is a simplified chat room structure for list response
type Room struct {
	Id          uint   `json:"id" dc:"Room ID"`
	Name        string `json:"name" dc:"Room name"`
	Description string `json:"description" dc:"Room description"`
	IsPrivate   bool   `json:"isPrivate" dc:"Whether the room is private"`
	CreatorId   uint   `json:"creatorId" dc:"ID of user who created the room"`
	UserCount   int    `json:"userCount" dc:"Number of users in the room"`
}

// DetailReq is the request for getting room details
type DetailReq struct {
	g.Meta `path:"/chatroom/detail/{id}" method:"get" tags:"ChatRoom" summary:"Get chat room details" auth:"true"`
	Id     uint `v:"required|min:1" dc:"Room ID"`
}

// DetailRes is the response for room details
type DetailRes struct {
	Id          uint   `json:"id" dc:"Room ID"`
	Name        string `json:"name" dc:"Room name"`
	Description string `json:"description" dc:"Room description"`
	IsPrivate   bool   `json:"isPrivate" dc:"Whether the room is private"`
	CreatorId   uint   `json:"creatorId" dc:"ID of user who created the room"`
	UserCount   int    `json:"userCount" dc:"Number of users in the room"`
	CreatedAt   string `json:"createdAt" dc:"Creation time"`
}

// JoinReq is the request for joining a chat room
type JoinReq struct {
	g.Meta `path:"/chatroom/join/{id}" method:"post" tags:"ChatRoom" summary:"Join a chat room" auth:"true"`
	Id     uint `v:"required|min:1" dc:"Room ID"`
}

// JoinRes is the response for joining a chat room
type JoinRes struct {
	Success bool `json:"success" dc:"Whether the operation was successful"`
}

// LeaveReq is the request for leaving a chat room
type LeaveReq struct {
	g.Meta `path:"/chatroom/leave/{id}" method:"post" tags:"ChatRoom" summary:"Leave a chat room" auth:"true"`
	Id     uint `v:"required|min:1" dc:"Room ID"`
}

// LeaveRes is the response for leaving a chat room
type LeaveRes struct {
	Success bool `json:"success" dc:"Whether the operation was successful"`
}
