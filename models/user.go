package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

/*
User table structure
*/
type User struct {
	UID       uint32 `json:"uid" gorm:"primary_key"`
	Uname     string `json:"uname" gorm:"unique;not null"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"createdAt"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", makeTimestamp())
	return nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
