package file

import (
	"os"

	"github.com/dhowden/tag"
)

func ExtractTag(path string) (tag.Metadata, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	tags, err := tag.ReadFrom(fd)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
