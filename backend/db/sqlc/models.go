// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"database/sql"
	"time"

	"github.com/sqlc-dev/pqtype"
)

type Driver struct {
	ID           int64          `json:"id"`
	Year         sql.NullInt32  `json:"year"`
	Make         sql.NullString `json:"make"`
	Model        sql.NullString `json:"model"`
	Color        sql.NullString `json:"color"`
	LicensePlate sql.NullString `json:"license_plate"`
	Status       int32          `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
}

type Trip struct {
	ID              int64                 `json:"id"`
	UserID          int64                 `json:"user_id"`
	DriverID        int64                 `json:"driver_id"`
	IsStarted       bool                  `json:"is_started"`
	IsCompleted     bool                  `json:"is_completed"`
	Origin          pqtype.NullRawMessage `json:"origin"`
	Destination     pqtype.NullRawMessage `json:"destination"`
	DestinationName sql.NullString        `json:"destination_name"`
	DriverLocation  pqtype.NullRawMessage `json:"driver_location"`
	CreatedAt       time.Time             `json:"created_at"`
}

type User struct {
	ID        int64          `json:"id"`
	Name      sql.NullString `json:"name"`
	Phone     string         `json:"phone"`
	LoginCode sql.NullString `json:"login_code"`
	CreatedAt time.Time      `json:"created_at"`
}
