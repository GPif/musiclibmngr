package file

import (
	"fmt"
	"os"

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
	return fmt.Sprintf(`
		localMetadata{
			localPath: %s,
			title: %s,
			artist: %s,
			album: %s,
		}
		`, m.localPath, m.Metadata.Title(), m.Metadata.Artist(), m.Metadata.Album())
}

func (m *localMetadata) GetLocalPath() string {
    return m.localPath
}

func ExtractTag(path string) (LocalMetadata, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	tags, err := tag.ReadFrom(fd)
	if err != nil {
		return nil, err
	}

	res := &localMetadata{
		Metadata: tags,
		localPath: path,
	}
	return res, nil
}
