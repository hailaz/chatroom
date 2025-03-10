package dao

import (
	"context"
	"goframechat/internal/model/entity"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
)

// UserDao handles database operations for users
type UserDao struct{}

// UserTable is the name of the user table
const UserTable = "users"

// NewUserDao returns a new UserDao instance
func NewUserDao() *UserDao {
	return &UserDao{}
}

// GetByID retrieves a user by ID
func (dao *UserDao) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user *entity.User
	err := Model(ctx, UserTable).Where("id", id).Scan(&user)
	return user, err
}

// GetByUsername retrieves a user by username
func (dao *UserDao) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user *entity.User
	err := Model(ctx, UserTable).Where("username", username).Scan(&user)
	return user, err
}

// Create creates a new user
func (dao *UserDao) Create(ctx context.Context, user *entity.User) (uint, error) {
	if user.Password == "" {
		// Set a default password for admin if none provided
		if user.Username == "admin" {
			pwd, err := gmd5.EncryptString("admin123")
			if err != nil {
				return 0, err
			}
			user.Password = pwd
		}
	} else {
		// Hash the password before storing
		pwd, err := gmd5.EncryptString(user.Password)
		if err != nil {
			return 0, err
		}
		user.Password = pwd
	}

	result, err := Model(ctx, UserTable).Data(user).Insert()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return uint(id), err
}

// Update updates a user
func (dao *UserDao) Update(ctx context.Context, id uint, data g.Map) error {
	// Add update timestamp
	data["updated_at"] = time.Now()

	_, err := Model(ctx, UserTable).Where("id", id).Data(data).Update()
	return err
}

// UpdateStatus updates a user's status
func (dao *UserDao) UpdateStatus(ctx context.Context, id uint, status int) error {
	data := g.Map{
		"status":     status,
		"updated_at": time.Now(),
	}
	if status == 1 {
		data["last_login"] = time.Now()
	}

	_, err := Model(ctx, UserTable).Where("id", id).Data(data).Update()
	return err
}

// VerifyPassword checks if a password matches the user's stored password
func (dao *UserDao) VerifyPassword(ctx context.Context, username, password string) (bool, *entity.User, error) {
	user, err := dao.GetByUsername(ctx, username)
	if err != nil {
		return false, nil, err
	}
	if user == nil {
		return false, nil, nil
	}

	// Hash the password and compare
	hashedPassword, err := gmd5.EncryptString(password)
	if err != nil {
		return false, nil, err
	}

	return hashedPassword == user.Password, user, nil
}
