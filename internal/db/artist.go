package db

import (
	"gorm.io/gorm"
)

type Artist struct {
	gorm.Model
	ID            uint `gorm:"primaryKey;autoIncrement:true"`
	Name          string
	MusicBrainzID *string
}

// FindOrCreateArtist finds an artist by name or creates it if it doesn't exist
func (db *Artist) FindOrCreateArtist(tx *gorm.DB, name string) (*Artist, error) {
	var artist Artist
	err := tx.Where("name = ?", name).First(&artist).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Artist doesn't exist, create it
			artist.Name = name
			err = tx.Create(&artist).Error
			if err != nil {
				return nil, err
			}
			return &artist, nil
		}
		return nil, err
	}
	return &artist, nil
}
