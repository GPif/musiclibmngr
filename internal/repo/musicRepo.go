package repo

import (
	"musiclibmngr/internal/db"

	"gorm.io/gorm"
)

type MusicRepo struct {
	db *gorm.DB
}

type MusicFile struct {
	Path      string
	Artist    string
	Record    string
	Release   string
	TrackNb   int
	Supportnb int
}

func NewMusicRepo(db *gorm.DB) *MusicRepo {
	return &MusicRepo{db: db}
}

func (r *MusicRepo) CreateFile(musicFile MusicFile) error {
	var artist db.Artist
	var release db.Release
	var record db.Record
	var localFile db.LocalFile
	r.db.Where(
		db.Artist{Name: musicFile.Artist},
	).FirstOrCreate(&artist)

	r.db.Where(
		db.Release{Title: musicFile.Release},
	).Assign(
		db.Release{ArtistID: artist.ID},
	).FirstOrCreate(&release)

	r.db.Where(
		db.Record{Title: musicFile.Record},
	).Assign(
		db.Record{
			ReleaseID:     release.ID,
			SupportNumber: musicFile.Supportnb,
			TrackNumber:   musicFile.TrackNb,
		},
	).FirstOrCreate(&record)

	r.db.Where(
		db.LocalFile{Path: musicFile.Path},
	).Assign(
		db.LocalFile{RecordID: record.ID},
	).FirstOrCreate(&localFile)

	return nil
}
