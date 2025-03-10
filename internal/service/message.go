package service

import (
	"chatroom/api/chat"
	"chatroom/internal/consts"
	"chatroom/internal/dao"
	"chatroom/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

// MessageService handles message-related business logic
type MessageService struct {
	messageDao *dao.MessageDao
	roomDao    *dao.ChatRoomDao
}

// NewMessageService creates a new MessageService instance
func NewMessageService() *MessageService {
	return &MessageService{
		messageDao: dao.NewMessageDao(),
		roomDao:    dao.NewChatRoomDao(),
	}
}

// CreateMessage creates a new chat message
func (s *MessageService) CreateMessage(ctx context.Context, userId uint, msg *chat.MessageReq) (uint, error) {
	// Check if user is in the room
	inRoom, err := s.roomDao.IsUserInRoom(ctx, msg.RoomId, userId)
	if err != nil {
		return 0, err
	}
	if !inRoom {
		return 0, gerror.New(consts.ErrNotInRoom)
	}

	// Create message
	message := &entity.Message{
		RoomId:  msg.RoomId,
		UserId:  userId,
		Content: msg.Content,
		Type:    msg.Type,
	}

	return s.messageDao.Create(ctx, message)
}

// GetHistory retrieves chat message history
func (s *MessageService) GetHistory(ctx context.Context, userId uint, req *chat.HistoryReq) (*chat.HistoryRes, error) {
	// Check if user is in the room
	inRoom, err := s.roomDao.IsUserInRoom(ctx, req.RoomId, userId)
	if err != nil {
		return nil, err
	}
	if !inRoom {
		return nil, gerror.New(consts.ErrNotInRoom)
	}

	// Get messages with user information
	messages, total, err := s.messageDao.GetRoomMessagesWithUser(ctx, req.RoomId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	messageList := make([]chat.MessageRes, 0, len(messages))
	for _, m := range messages {
		msg := chat.MessageRes{}
		gconv.Struct(m, &msg, map[string]string{
			"created_at": "timestamp",
		})
		messageList = append(messageList, msg)
	}

	// glog.Debug(ctx, "History Response: ", messageList)

	return &chat.HistoryRes{
		Messages: messageList,
		Total:    total,
		Page:     req.Page,
		Size:     req.Size,
	}, nil
}

// GetRoomMembers retrieves all members in a chat room
func (s *MessageService) GetRoomMembers(ctx context.Context, userId uint, req *chat.RoomMembersReq) (*chat.RoomMembersRes, error) {
	// Check if user is in the room
	inRoom, err := s.roomDao.IsUserInRoom(ctx, req.Id, userId)
	if err != nil {
		return nil, err
	}
	if !inRoom {
		return nil, gerror.New(consts.ErrNotInRoom)
	}

	// Get room members
	users, err := s.roomDao.ListRoomUsers(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	members := make([]chat.Member, 0, len(users))
	for _, u := range users {
		members = append(members, chat.Member{
			Id:       u.Id,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Status:   u.Status,
		})
	}

	return &chat.RoomMembersRes{
		Members: members,
		Total:   len(members),
	}, nil
}
