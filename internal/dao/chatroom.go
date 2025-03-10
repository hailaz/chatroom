package dao

import (
	"context"
	"goframechat/internal/model/entity"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// ChatRoomDao handles database operations for chat rooms
type ChatRoomDao struct{}

// ChatRoomTable is the name of the chat room table
const ChatRoomTable = "chatrooms"

// RoomUserTable is the name of the room-user relationship table
const RoomUserTable = "room_users"

// NewChatRoomDao returns a new ChatRoomDao instance
func NewChatRoomDao() *ChatRoomDao {
	return &ChatRoomDao{}
}

// GetByID retrieves a chat room by ID
func (dao *ChatRoomDao) GetByID(ctx context.Context, id uint) (*entity.ChatRoom, error) {
	var chatRoom *entity.ChatRoom
	err := Model(ctx, ChatRoomTable).Where("id", id).Scan(&chatRoom)
	return chatRoom, err
}

// Create creates a new chat room
func (dao *ChatRoomDao) Create(ctx context.Context, chatRoom *entity.ChatRoom) (uint, error) {
	result, err := Model(ctx, ChatRoomTable).Data(g.Map{
		"name":        chatRoom.Name,
		"description": chatRoom.Description,
		"creator_id":  chatRoom.CreatorId,
		"is_private":  chatRoom.IsPrivate,
	}).Insert()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Add creator as a member of the room
	_, err = Model(ctx, RoomUserTable).Data(g.Map{
		"room_id":   id,
		"user_id":   chatRoom.CreatorId,
		"joined_at": time.Now(),
	}).Insert()

	return uint(id), err
}

// Update updates a chat room
func (dao *ChatRoomDao) Update(ctx context.Context, id uint, data g.Map) error {
	// Add update timestamp
	data["updated_at"] = time.Now()

	_, err := Model(ctx, ChatRoomTable).Where("id", id).Data(data).Update()
	return err
}

// Delete deletes a chat room
func (dao *ChatRoomDao) Delete(ctx context.Context, id uint) error {
	// First delete all room-user relationships
	_, err := Model(ctx, RoomUserTable).Where("room_id", id).Delete()
	if err != nil {
		return err
	}

	// Then delete the room
	_, err = Model(ctx, ChatRoomTable).Where("id", id).Delete()
	return err
}

// List returns a paginated list of chat rooms
func (dao *ChatRoomDao) List(ctx context.Context, page, size int) (rooms []entity.ChatRoom, total int, err error) {
	model := Model(ctx, ChatRoomTable)

	// Get total count
	total, err = model.Count()
	if err != nil {
		return nil, 0, err
	}

	// Get paginated rooms
	err = model.Page(page, size).Order("id DESC").Scan(&rooms)
	return rooms, total, err
}

// GetUserCount returns the number of users in a specific chat room
func (dao *ChatRoomDao) GetUserCount(ctx context.Context, roomId uint) (int, error) {
	return Model(ctx, RoomUserTable).Where("room_id", roomId).Count()
}

// AddUser adds a user to a chat room
func (dao *ChatRoomDao) AddUser(ctx context.Context, roomId, userId uint) error {
	// Check if the user is already in the room
	count, err := Model(ctx, RoomUserTable).
		Where("room_id", roomId).
		Where("user_id", userId).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		// User already in the room
		return nil
	}

	// Add user to room
	_, err = Model(ctx, RoomUserTable).Data(g.Map{
		"room_id":   roomId,
		"user_id":   userId,
		"joined_at": time.Now(),
	}).Insert()
	return err
}

// RemoveUser removes a user from a chat room
func (dao *ChatRoomDao) RemoveUser(ctx context.Context, roomId, userId uint) error {
	_, err := Model(ctx, RoomUserTable).
		Where("room_id", roomId).
		Where("user_id", userId).
		Delete()
	return err
}

// ListRoomUsers returns all users in a specific chat room
func (dao *ChatRoomDao) ListRoomUsers(ctx context.Context, roomId uint) (users []entity.User, err error) {
	err = Model(ctx, UserTable).
		As("u").
		InnerJoin("room_users ru", "ru.user_id = u.id").
		Where("ru.room_id", roomId).
		Fields("u.*").
		Scan(&users)
	return
}

// IsUserInRoom checks if a user is in a specific chat room
func (dao *ChatRoomDao) IsUserInRoom(ctx context.Context, roomId, userId uint) (bool, error) {
	count, err := Model(ctx, RoomUserTable).
		Where("room_id", roomId).
		Where("user_id", userId).
		Count()
	return count > 0, err
}

// GetUserRooms returns all chat rooms that a user is in
func (dao *ChatRoomDao) GetUserRooms(ctx context.Context, userId uint) (rooms []entity.ChatRoom, err error) {
	err = Model(ctx, ChatRoomTable).
		As("cr").
		InnerJoin("room_users ru", "ru.room_id = cr.id").
		Where("ru.user_id", userId).
		Fields("cr.*").
		Scan(&rooms)
	return
}
