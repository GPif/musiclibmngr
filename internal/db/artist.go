package db

import (
	"gorm.io/gorm"
)

type Artist struct {
	gorm.Model
	ID            uint `gorm:"primaryKey;autoIncrement:true"`
	Name          string
	MusicBrainzID string
}
