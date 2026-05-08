package importer

import (
	"fmt"
)

type ReleaseInfo struct {
	Title   string
	Artist  string
	TrackNb int
	Year    int
}

type ImportTask struct {
	Paths            []string
	Records          []LocalMetadata
	ReleaseInfo      ReleaseInfo
	ReleaseCandidate []any
	BestMatch        []any
}

func (t ImportTask) String() string {
	return fmt.Sprintf(`ImportTask{
Paths: %v
Records: %v
ReleaseCandidate: %v
ReleaseInfo: %v
BestMatch: %v
`, t.Paths, t.Records, t.ReleaseCandidate, t.ReleaseInfo, t.BestMatch)
}
