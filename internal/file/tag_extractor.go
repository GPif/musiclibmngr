package file

import (
	"os"

	"github.com/dhowden/tag"
)

func ExtractTag(path string) (map[string]any, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	tags, err := tag.ReadFrom(fd)
	if err != nil {
		return nil, err
	}

	return tags.Raw(), nil
}
