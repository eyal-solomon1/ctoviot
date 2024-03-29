// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID        int64        `json:"id"`
	Username  string       `json:"username"`
	VideoName string       `json:"video_name"`
	Amount    int64        `json:"amount"`
	CreatedAt sql.NullTime `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Balance           int64     `json:"balance"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type Video struct {
	ID              int64     `json:"id"`
	Owner           string    `json:"owner"`
	VideoName       string    `json:"video_name"`
	VideoIdentifier string    `json:"video_identifier"`
	VideoLength     int64     `json:"video_length"`
	VideoRemotePath string    `json:"video_remote_path"`
	VideoDecs       string    `json:"video_decs"`
	CreatedAt       time.Time `json:"created_at"`
}
