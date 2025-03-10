package consts

const (
	// User status constants
	UserStatusOffline = 0
	UserStatusOnline  = 1

	// Message type constants
	MessageTypeText   = 0
	MessageTypeImage  = 1
	MessageTypeFile   = 2
	MessageTypeSystem = 3

	// WebSocket message types
	WsMsgTypeText         = 1 // Text message
	WsMsgTypeJoin         = 2 // User joined
	WsMsgTypeLeave        = 3 // User left
	WsMsgTypeUserList     = 4 // User list update
	WsMsgTypeError        = 5 // Error message
	WsMsgTypeNotification = 6 // System notification

	// JWT related constants
	JwtExpireTime = 86400 // JWT token expire time in seconds (24 hours)
	JwtIssuer     = "GoFrameChat"

	// Default values
	DefaultAvatar = "/resource/image/avatar/default.png"

	// Error messages
	ErrNotInRoom = "User is not in the chat room"
)

// ContextKey is the key type for context values
type ContextKey string

const (
	ContextKeyUser ContextKey = "user" // Context key for storing user information
)
