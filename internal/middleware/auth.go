package middleware

import (
	"chatroom/internal/consts"
	"chatroom/internal/dao"
	"chatroom/internal/service"
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Auth creates a JWT authentication middleware
func Auth(r *ghttp.Request) {
	// Get token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		g.Log().Error(r.Context(), "Missing authorization header")
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeNotAuthorized.Code(),
			Message: "Missing authorization header",
		})
		r.Exit()
		return
	}

	// Check token format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		g.Log().Error(r.Context(), "Invalid authorization format:", authHeader)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeNotAuthorized.Code(),
			Message: "Invalid authorization format",
		})
		r.Exit()
		return
	}

	// Parse and validate token
	jwtService := service.NewJwtService()
	claims, err := jwtService.ParseToken(parts[1])
	if err != nil {
		g.Log().Error(r.Context(), "Invalid or expired token:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeNotAuthorized.Code(),
			Message: "Invalid or expired token",
		})
		r.Exit()
		return
	}

	// Get user from database
	userDao := dao.NewUserDao()
	user, err := userDao.GetByID(context.Background(), claims.UserId)
	if err != nil {
		g.Log().Error(r.Context(), "Failed to get user information:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeInternalError.Code(),
			Message: "Failed to get user information",
		})
		r.Exit()
		return
	}
	if user == nil {
		g.Log().Error(r.Context(), "User not found, ID:", claims.UserId)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeNotAuthorized.Code(),
			Message: "User not found",
		})
		r.Exit()
		return
	}

	// Store user in context
	r.SetCtxVar(consts.ContextKeyUser, user)

	r.Middleware.Next()
}
