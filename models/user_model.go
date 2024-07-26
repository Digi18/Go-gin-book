package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Email     string    `json:"email" gorm:"unique;not null;size:255"`
	Password  string    `json:"password" gorm:"not null;size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
