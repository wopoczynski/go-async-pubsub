package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PersistingStruct struct {
	gorm.Model

	ID   uuid.UUID `json:"id"          gorm:"type:char(36);primary_key;"`
	Data string    `json:"data"`
}
