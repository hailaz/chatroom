package entity

import (
	"time"
)

// User represents a chat user
type User struct {
	Id        uint      `json:"id"        description:"User ID"`
	Username  string    `json:"username"  description:"Username for login"`
	Password  string    `json:"-"         description:"Password (hashed)"`
	Nickname  string    `json:"nickname"  description:"Display name"`
	Avatar    string    `json:"avatar"    description:"User avatar URL"`
	Status    int       `json:"status"    description:"User status: 0-offline, 1-online"`
	LastLogin time.Time `json:"lastLogin" description:"Last login time"`
	CreatedAt time.Time `json:"createdAt" description:"Created time"`
	UpdatedAt time.Time `json:"updatedAt" description:"Updated time"`
}
