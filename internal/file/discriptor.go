package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Descriptor struct {
	path string
}

func NewDescriptor(path string) *Descriptor {
	return &Descriptor{path: path}
}

func (d *Descriptor) IsAudio() (bool, error) {
	t, err := detectType(d.path)
	if err != nil {
		return false, fmt.Errorf("failed to detect type: %w", err)
	}
	if t == "audio/mpeg" {
		return true, nil
	}
	if strings.HasSuffix(d.path, ".flac") && t == "application/octet-stream" {
		return true, nil
	}
	return false, nil
}

func detectType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := io.ReadFull(f, buf)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buf[:n]), nil
}
