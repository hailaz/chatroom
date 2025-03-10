package dao

import (
	"context"
	"goframechat/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

// MessageDao handles database operations for messages
type MessageDao struct{}

// MessageTable is the name of the message table
const MessageTable = "messages"

// NewMessageDao creates a new MessageDao instance
func NewMessageDao() *MessageDao {
	return &MessageDao{}
}

// Create stores a new message in the database
func (dao *MessageDao) Create(ctx context.Context, message *entity.Message) (uint, error) {
	result, err := Model(ctx, MessageTable).Data(message).Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint(id), err
}

// GetRoomMessagesWithUser retrieves messages with user information for a room
func (dao *MessageDao) GetRoomMessagesWithUser(ctx context.Context, roomId uint, page, size int) ([]g.Map, int, error) {
	model := Model(ctx, MessageTable).
		As("m").
		LeftJoin("users u", "u.id = m.user_id").
		Where("m.room_id", roomId)

	// Get total count
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	// Get messages with user information
	var messages []g.Map
	err = model.Fields(`
		m.*,
		u.username,
		u.nickname,
		u.avatar
	`).
		Order("m.created_at DESC").
		Page(page, size).
		Scan(&messages)

	return messages, total, err
}

// GetUserMessages retrieves all messages sent by a user
func (dao *MessageDao) GetUserMessages(ctx context.Context, userId uint, page, size int) ([]entity.Message, int, error) {
	model := Model(ctx, MessageTable).Where("user_id", userId)

	// Get total count
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	// Get messages
	var messages []entity.Message
	err = model.Order("created_at DESC").Page(page, size).Scan(&messages)

	return messages, total, err
}

// DeleteRoomMessages deletes all messages in a room
func (dao *MessageDao) DeleteRoomMessages(ctx context.Context, roomId uint) error {
	_, err := Model(ctx, MessageTable).Where("room_id", roomId).Delete()
	return err
}

// GetMessageById retrieves a message by its ID
func (dao *MessageDao) GetMessageById(ctx context.Context, id uint) (*entity.Message, error) {
	var message *entity.Message
	err := Model(ctx, MessageTable).Where("id", id).Scan(&message)
	return message, err
}
