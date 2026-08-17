package model

import "time"

// User 系统用户。
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Role      string    `gorm:"size:32;index;not null" json:"role"`
	OwnerID   *uint     `gorm:"index" json:"owner_id"`
	Status    string    `gorm:"size:16;default:Active" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
