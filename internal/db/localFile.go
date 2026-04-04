package db

import "gorm.io/gorm"

type LocalFile struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;autoIncrement:true"`
	Path     string
	RecordID int
}
