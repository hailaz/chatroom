package entity

import (
	"time"
)

// ChatRoom represents a chat room where users can join and send messages
type ChatRoom struct {
	Id          uint      `json:"id"          description:"Room ID"`
	Name        string    `json:"name"        description:"Room name"`
	Description string    `json:"description" description:"Room description"`
	CreatorId   uint      `json:"creatorId"   description:"ID of user who created the room"`
	IsPrivate   bool      `json:"isPrivate"   description:"Whether the room is private (invite only)"`
	CreatedAt   time.Time `json:"createdAt"   description:"Created time"`
	UpdatedAt   time.Time `json:"updatedAt"   description:"Updated time"`
}
