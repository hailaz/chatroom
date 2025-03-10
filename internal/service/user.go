package service

import (
	"chatroom/api/user"
	"chatroom/internal/consts"
	"chatroom/internal/dao"
	"chatroom/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

// UserService handles user-related business logic
type UserService struct {
	userDao *dao.UserDao
}

// NewUserService creates a new UserService instance
func NewUserService() *UserService {
	return &UserService{
		userDao: dao.NewUserDao(),
	}
}

// Register handles user registration
func (s *UserService) Register(ctx context.Context, req *user.RegisterReq) (*user.RegisterRes, error) {
	// Check if username exists
	existingUser, err := s.userDao.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, gerror.New("Username already exists")
	}

	// Create new user
	newUser := &entity.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   consts.DefaultAvatar,
		Status:   consts.UserStatusOffline,
	}

	// Insert user into database
	userId, err := s.userDao.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &user.RegisterRes{
		Id:       userId,
		Username: req.Username,
		Nickname: req.Nickname,
	}, nil
}

// Login handles user login
func (s *UserService) Login(ctx context.Context, req *user.LoginReq) (*user.LoginRes, error) {
	// Verify password
	ok, u, err := s.userDao.VerifyPassword(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	if !ok || u == nil {
		return nil, gerror.New("Invalid username or password")
	}

	// Generate JWT token
	jwtService := NewJwtService()
	token, err := jwtService.GenerateToken(u)
	if err != nil {
		return nil, err
	}

	// Update user status
	err = s.userDao.UpdateStatus(ctx, u.Id, consts.UserStatusOnline)
	if err != nil {
		return nil, err
	}

	return &user.LoginRes{
		Token:    token,
		Id:       u.Id,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
	}, nil
}

// GetProfile retrieves user profile
func (s *UserService) GetProfile(ctx context.Context, userId uint) (*user.ProfileRes, error) {
	u, err := s.userDao.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, gerror.New("User not found")
	}

	return &user.ProfileRes{
		Id:       u.Id,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Status:   u.Status,
	}, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(ctx context.Context, userId uint, req *user.UpdateProfileReq) (*user.UpdateProfileRes, error) {
	u, err := s.userDao.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, gerror.New("User not found")
	}

	// Update user profile
	err = s.userDao.Update(ctx, userId, map[string]interface{}{
		"nickname": req.Nickname,
		"avatar":   req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return &user.UpdateProfileRes{
		Id:       u.Id,
		Username: u.Username,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}, nil
}
