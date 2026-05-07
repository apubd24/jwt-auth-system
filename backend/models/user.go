package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Fullname  string    `gorm:"not null" json:"fullname"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `json:"-"`                            // never sent to client
	Role      string    `gorm:"default:readonly" json:"role"` // admin or readonly
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// func (User) TableName() string {
// 	return "user" // singular table name
// }
