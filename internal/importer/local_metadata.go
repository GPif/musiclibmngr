package importer

import (
	"fmt"

	"github.com/dhowden/tag"
)

type LocalMetadata interface {
	tag.Metadata
	GetLocalPath() string
}

// Private wrapper that implements it
type localMetadata struct {
	tag.Metadata
	localPath string
}

func (m *localMetadata) String() string {
	t, tn := m.Metadata.Track()
	return fmt.Sprintf(`
		localMetadata{
			artist: %s,
			album: %s,
			localPath: %s,
			title: %s,
			track: %d/%d
		}
		`, m.localPath, m.Metadata.Title(), m.Metadata.AlbumArtist(), m.Metadata.Album(), t, tn)
}

func (m *localMetadata) GetLocalPath() string {
	return m.localPath
}
