package file

import (
	// "encoding/json"
	"fmt"
	"testing"

)

func TestExtractTag(t *testing.T) {
	r, err := ExtractTag("../../testdata/09 Parties for Prostitutes.mp3")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Format : ", r.Format())
	fmt.Println("FileType : ", r.FileType())
	fmt.Println("Title : ", r.Title())
	fmt.Println("Album : ", r.Album())
	fmt.Println("Artist : ", r.Artist())
	fmt.Println("AlbumArtist : ", r.AlbumArtist())
	fmt.Println("Composer : ", r.Composer())
	fmt.Println("Genre : ", r.Genre())
	fmt.Println("Year : ", r.Year())
	t1, t2 := r.Track()
	fmt.Println("Track : ", t1, t2)
	d1, d2 := r.Disc()
	fmt.Println("Disc : ", d1, d2)
}
