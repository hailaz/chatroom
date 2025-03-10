package chatroom

import (
	"chatroom/api/chatroom"
	"chatroom/internal/consts"
	"chatroom/internal/model/entity"
	"chatroom/internal/service"
	"context"
)

// Controller handles chat room-related HTTP requests
type Controller struct {
	roomService *service.ChatRoomService
}

// NewController creates a new chat room controller
func NewController() *Controller {
	return &Controller{
		roomService: service.NewChatRoomService(),
	}
}

// Create handles chat room creation
func (c *Controller) Create(ctx context.Context, req *chatroom.CreateReq) (res *chatroom.CreateRes, err error) {
	ctxUser := ctx.Value(consts.ContextKeyUser).(*entity.User)
	return c.roomService.Create(ctx, ctxUser.Id, req)
}

// List returns a list of available chat rooms
func (c *Controller) List(ctx context.Context, req *chatroom.ListReq) (res *chatroom.ListRes, err error) {
	return c.roomService.List(ctx, req)
}

// Detail returns details of a specific chat room
func (c *Controller) Detail(ctx context.Context, req *chatroom.DetailReq) (res *chatroom.DetailRes, err error) {
	return c.roomService.Detail(ctx, req)
}

// Join handles a user joining a chat room
func (c *Controller) Join(ctx context.Context, req *chatroom.JoinReq) (res *chatroom.JoinRes, err error) {
	ctxUser := ctx.Value(consts.ContextKeyUser).(*entity.User)
	return c.roomService.Join(ctx, ctxUser.Id, req)
}

// Leave handles a user leaving a chat room
func (c *Controller) Leave(ctx context.Context, req *chatroom.LeaveReq) (res *chatroom.LeaveRes, err error) {
	ctxUser := ctx.Value(consts.ContextKeyUser).(*entity.User)
	return c.roomService.Leave(ctx, ctxUser.Id, req)
}

// Delete handles chat room deletion
func (c *Controller) Delete(ctx context.Context, req *chatroom.DeleteReq) (res *chatroom.DeleteRes, err error) {
	ctxUser := ctx.Value(consts.ContextKeyUser).(*entity.User)
	return c.roomService.Delete(ctx, ctxUser.Id, req)
}
