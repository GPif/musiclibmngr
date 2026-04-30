package db

import "gorm.io/gorm"

type Release struct {
	gorm.Model
	ID            uint `gorm:"primaryKey;autoIncrement:true"`
	Title         string
	ArtistID      uint
	MusicBrainzID *string
}
