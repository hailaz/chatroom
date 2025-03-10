package chat

import (
	"context"
	"goframechat/api/chat"
	"goframechat/internal/consts"
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

	// Get user from context (set by auth middleware)
	ctxUser := r.Context().Value(consts.ContextKeyUser).(*entity.User)

	// Handle WebSocket connection
	c.wsManager.HandleWebSocket(r, ctxUser, req.RoomId)
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
