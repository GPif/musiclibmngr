package db

import "gorm.io/gorm"

type LocalFile struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;autoIncrement:true"`
	Path     string
	RecordID uint
}

func AddFile(path string, track int, title string, album string, artist string) {

}
