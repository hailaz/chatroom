package service

import (
	"context"
	"goframechat/api/chatroom"
	"goframechat/internal/dao"
	"goframechat/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
)

// ChatRoomService handles chat room business logic
type ChatRoomService struct {
	roomDao *dao.ChatRoomDao
}

// NewChatRoomService creates a new ChatRoomService instance
func NewChatRoomService() *ChatRoomService {
	return &ChatRoomService{
		roomDao: dao.NewChatRoomDao(),
	}
}

// Create creates a new chat room
func (s *ChatRoomService) Create(ctx context.Context, userId uint, req *chatroom.CreateReq) (*chatroom.CreateRes, error) {
	room := &entity.ChatRoom{
		Name:        req.Name,
		Description: req.Description,
		CreatorId:   userId,
		IsPrivate:   req.IsPrivate,
	}

	roomId, err := s.roomDao.Create(ctx, room)
	if err != nil {
		return nil, err
	}

	return &chatroom.CreateRes{
		Id:          roomId,
		Name:        req.Name,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
	}, nil
}

// List returns a paginated list of chat rooms
func (s *ChatRoomService) List(ctx context.Context, req *chatroom.ListReq) (*chatroom.ListRes, error) {
	rooms, total, err := s.roomDao.List(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// Convert to response format and get user count for each room
	list := make([]chatroom.Room, 0, len(rooms))
	for _, r := range rooms {
		userCount, err := s.roomDao.GetUserCount(ctx, r.Id)
		if err != nil {
			return nil, err
		}

		list = append(list, chatroom.Room{
			Id:          r.Id,
			Name:        r.Name,
			Description: r.Description,
			IsPrivate:   r.IsPrivate,
			CreatorId:   r.CreatorId,
			UserCount:   userCount,
		})
	}

	return &chatroom.ListRes{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// Detail returns details of a chat room
func (s *ChatRoomService) Detail(ctx context.Context, req *chatroom.DetailReq) (*chatroom.DetailRes, error) {
	room, err := s.roomDao.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, gerror.New("Chat room not found")
	}

	userCount, err := s.roomDao.GetUserCount(ctx, room.Id)
	if err != nil {
		return nil, err
	}

	return &chatroom.DetailRes{
		Id:          room.Id,
		Name:        room.Name,
		Description: room.Description,
		IsPrivate:   room.IsPrivate,
		CreatorId:   room.CreatorId,
		UserCount:   userCount,
		CreatedAt:   room.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// Join lets a user join a chat room
func (s *ChatRoomService) Join(ctx context.Context, userId uint, req *chatroom.JoinReq) (*chatroom.JoinRes, error) {
	// Check if room exists
	room, err := s.roomDao.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, gerror.New("Chat room not found")
	}

	// Check if user is already in room
	isInRoom, err := s.roomDao.IsUserInRoom(ctx, req.Id, userId)
	if err != nil {
		return nil, err
	}
	if isInRoom {
		return &chatroom.JoinRes{Success: true}, nil
	}

	// Add user to room
	err = s.roomDao.AddUser(ctx, req.Id, userId)
	if err != nil {
		return nil, err
	}

	return &chatroom.JoinRes{Success: true}, nil
}

// Leave lets a user leave a chat room
func (s *ChatRoomService) Leave(ctx context.Context, userId uint, req *chatroom.LeaveReq) (*chatroom.LeaveRes, error) {
	// Check if room exists
	room, err := s.roomDao.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, gerror.New("Chat room not found")
	}

	// Check if user is in room
	isInRoom, err := s.roomDao.IsUserInRoom(ctx, req.Id, userId)
	if err != nil {
		return nil, err
	}
	if !isInRoom {
		return &chatroom.LeaveRes{Success: true}, nil
	}

	// Remove user from room
	err = s.roomDao.RemoveUser(ctx, req.Id, userId)
	if err != nil {
		return nil, err
	}

	return &chatroom.LeaveRes{Success: true}, nil
}

// Delete deletes a chat room if the user is the creator
func (s *ChatRoomService) Delete(ctx context.Context, userId uint, req *chatroom.DeleteReq) (*chatroom.DeleteRes, error) {
	// Check if room exists and user is the creator
	room, err := s.roomDao.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, gerror.New("Chat room not found")
	}
	if room.CreatorId != userId {
		return nil, gerror.New("Only the room creator can delete the room")
	}

	// Delete the room and all associated data
	messageDao := dao.NewMessageDao()
	if err := messageDao.DeleteRoomMessages(ctx, req.Id); err != nil {
		return nil, err
	}

	if err := s.roomDao.Delete(ctx, req.Id); err != nil {
		return nil, err
	}

	return &chatroom.DeleteRes{Success: true}, nil
}
