package db

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	ID            uint `gorm:"primaryKey;autoIncrement:true"`
	Title         string
	ReleaseID     uint
	SupportNumber int
	TrackNumber   int
	MusicBrainzID string
}
