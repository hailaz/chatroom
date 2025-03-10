package user

import (
	"github.com/gogf/gf/v2/frame/g"
)

// RegisterReq is the request for user registration
type RegisterReq struct {
	g.Meta   `path:"/user/register" method:"post" tags:"User" summary:"Register a new user"`
	Username string `v:"required|length:5,30|regex:[a-zA-Z0-9_]+" dc:"Username (5-30 chars, alphanumeric and underscore only)"`
	Password string `v:"required|length:6,30" dc:"Password (6-30 chars)"`
	Nickname string `v:"required|length:2,20" dc:"Display name (2-20 chars)"`
}

// RegisterRes is the response for user registration
type RegisterRes struct {
	Id       uint   `json:"id" dc:"User ID"`
	Username string `json:"username" dc:"Username"`
	Nickname string `json:"nickname" dc:"Display name"`
}

// LoginReq is the request for user login
type LoginReq struct {
	g.Meta   `path:"/user/login" method:"post" tags:"User" summary:"User login"`
	Username string `v:"required" dc:"Username"`
	Password string `v:"required" dc:"Password"`
}

// LoginRes is the response for user login
type LoginRes struct {
	Token    string `json:"token" dc:"JWT token for authentication"`
	Id       uint   `json:"id" dc:"User ID"`
	Username string `json:"username" dc:"Username"`
	Nickname string `json:"nickname" dc:"Display name"`
	Avatar   string `json:"avatar" dc:"User avatar URL"`
}

// ProfileReq is the request for getting user profile
type ProfileReq struct {
	g.Meta `path:"/user/profile" method:"get" tags:"User" summary:"Get user profile" auth:"true"`
}

// ProfileRes is the response for user profile
type ProfileRes struct {
	Id       uint   `json:"id" dc:"User ID"`
	Username string `json:"username" dc:"Username"`
	Nickname string `json:"nickname" dc:"Display name"`
	Avatar   string `json:"avatar" dc:"User avatar URL"`
	Status   int    `json:"status" dc:"User status"`
}

// UpdateProfileReq is the request for updating user profile
type UpdateProfileReq struct {
	g.Meta   `path:"/user/profile" method:"put" tags:"User" summary:"Update user profile" auth:"true"`
	Nickname string `v:"required|length:2,20" dc:"Display name (2-20 chars)"`
	Avatar   string `dc:"Avatar URL"`
}

// UpdateProfileRes is the response for updating user profile
type UpdateProfileRes struct {
	Id       uint   `json:"id" dc:"User ID"`
	Username string `json:"username" dc:"Username"`
	Nickname string `json:"nickname" dc:"Display name"`
	Avatar   string `json:"avatar" dc:"User avatar URL"`
}
