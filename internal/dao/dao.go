package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

// InitDatabase initializes the database and creates tables if they don't exist
func InitDatabase(ctx context.Context) error {
	// Create users table
	_, err := g.DB().Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(30) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			nickname VARCHAR(50) NOT NULL,
			avatar VARCHAR(255) DEFAULT '',
			status INTEGER DEFAULT 0,
			last_login DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		glog.Error(ctx, "Create users table failed:", err)
		return err
	}

	// Create chatrooms table
	_, err = g.DB().Exec(ctx, `
		CREATE TABLE IF NOT EXISTS chatrooms (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(50) NOT NULL,
			description VARCHAR(200) DEFAULT '',
			creator_id INTEGER NOT NULL,
			is_private BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (creator_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		glog.Error(ctx, "Create chatrooms table failed:", err)
		return err
	}

	// Create messages table
	_, err = g.DB().Exec(ctx, `
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			room_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			type INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (room_id) REFERENCES chatrooms(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		glog.Error(ctx, "Create messages table failed:", err)
		return err
	}

	// Create room_users table for many-to-many relationship between users and rooms
	_, err = g.DB().Exec(ctx, `
		CREATE TABLE IF NOT EXISTS room_users (
			room_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (room_id, user_id),
			FOREIGN KEY (room_id) REFERENCES chatrooms(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		glog.Error(ctx, "Create room_users table failed:", err)
		return err
	}

	return nil
}

// Model returns a model with transaction support
func Model(ctx context.Context, tableName string) *gdb.Model {
	return g.DB().Model(tableName).Safe().Ctx(ctx)
}
