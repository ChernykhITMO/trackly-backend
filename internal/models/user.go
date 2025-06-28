package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          int       `gorm:"column:id"`
	Username    string    `json:"username"`
	Email       string    `json:email`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Country     string
	City        string
	AvatarId    *uuid.UUID
	Password    string `json:"password"`
}
