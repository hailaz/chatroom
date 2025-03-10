package user

import (
	"chatroom/api/user"
	"chatroom/internal/consts"
	"chatroom/internal/model/entity"
	"chatroom/internal/service"
	"context"
)

// Controller handles user-related HTTP requests
type Controller struct {
	userService *service.UserService
}

// NewController creates a new user controller
func NewController() *Controller {
	return &Controller{
		userService: service.NewUserService(),
	}
}

// Register handles user registration
func (c *Controller) Register(ctx context.Context, req *user.RegisterReq) (res *user.RegisterRes, err error) {
	return c.userService.Register(ctx, req)
}

// Login handles user login
func (c *Controller) Login(ctx context.Context, req *user.LoginReq) (res *user.LoginRes, err error) {
	return c.userService.Login(ctx, req)
}

// Profile gets user profile information
func (c *Controller) Profile(ctx context.Context, req *user.ProfileReq) (res *user.ProfileRes, err error) {
	// Get user from context (set by auth middleware)
	ctxUser := ctx.Value(consts.ContextKeyUser).(*entity.User)
	return c.userService.GetProfile(ctx, ctxUser.Id)
}

// UpdateProfile updates user profile
func (c *Controller) UpdateProfile(ctx context.Context, req *user.UpdateProfileReq) (res *user.UpdateProfileRes, err error) {
	// Get user from context (set by auth middleware)
	ctxUser := ctx.Value(consts.ContextKeyUser).(*entity.User)
	return c.userService.UpdateProfile(ctx, ctxUser.Id, req)
}
