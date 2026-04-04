package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func New(dbPath string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Artist{})
	db.AutoMigrate(&Release{})
	db.AutoMigrate(&Record{})
	db.AutoMigrate(&LocalFile{})

	return &DB{db}, nil
}
