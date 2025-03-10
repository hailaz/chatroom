package chat

import (
	"context"
	"goframechat/api/chat"
	"goframechat/internal/consts"
	"goframechat/internal/dao"
	"goframechat/internal/model/entity"
	"goframechat/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Controller handles chat-related requests
type Controller struct {
	wsManager      *service.WebSocketManager
	messageService *service.MessageService
}

// NewController creates a new chat controller
func NewController() *Controller {
	return &Controller{
		wsManager:      service.GetWebSocketManager(),
		messageService: service.NewMessageService(),
	}
}

// Connect handles WebSocket connection requests
func (c *Controller) Connect(r *ghttp.Request) {
	var req *chat.ConnectReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    1,
			Message: err.Error(),
		})
		return
	}

	// Get token from URL parameter
	token := r.Get("token").String()
	if token == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    1,
			Message: "Missing token",
		})
		return
	}

	// Parse token
	jwtService := service.NewJwtService()
	claims, err := jwtService.ParseToken(token)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    1,
			Message: "Invalid token: " + err.Error(),
		})
		return
	}

	// Get user
	userDao := dao.NewUserDao()
	user, err := userDao.GetByID(r.Context(), claims.UserId)
	if err != nil || user == nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    1,
			Message: "User not found",
		})
		return
	}

	// Handle WebSocket connection
	c.wsManager.HandleWebSocket(r, user, req.RoomId)
}

// GetHistory returns chat message history
func (c *Controller) GetHistory(ctx context.Context, req *chat.HistoryReq) (res *chat.HistoryRes, err error) {
	ctxUser := ctx.Value(consts.ContextKeyUser).(*entity.User)
	return c.messageService.GetHistory(ctx, ctxUser.Id, req)
}

// GetRoomMembers returns all members in a chat room
func (c *Controller) GetRoomMembers(ctx context.Context, req *chat.RoomMembersReq) (res *chat.RoomMembersRes, err error) {
	ctxUser := ctx.Value(consts.ContextKeyUser).(*entity.User)
	return c.messageService.GetRoomMembers(ctx, ctxUser.Id, req)
}
