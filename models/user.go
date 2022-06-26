package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name      string    `gorm:"size:255"`
	Email     string    `gorm:"size:255;unique"`
	Password  string    `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
