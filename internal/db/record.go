package db

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	ID            uint `gorm:"primaryKey;autoIncrement:true"`
	Title         string
	ReleaseID     int
	SupportNumber *int
	TrackNumber   *int
	MusicBrainzID *string
}
