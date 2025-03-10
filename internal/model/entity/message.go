package entity

import (
	"time"
)

// Message represents a chat message sent in a room
type Message struct {
	Id        uint      `json:"id"        description:"Message ID"`
	RoomId    uint      `json:"roomId"    description:"Room the message belongs to"`
	UserId    uint      `json:"userId"    description:"User who sent the message"`
	Content   string    `json:"content"   description:"Message content"`
	Type      int       `json:"type"      description:"Message type: 0-text, 1-image, 2-file, 3-system"`
	CreatedAt time.Time `json:"createdAt" description:"Created time"`
}
