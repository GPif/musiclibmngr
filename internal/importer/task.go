package importer

import (
	"fmt"
	"log"
	"musiclibmngr/internal/file"
	"musiclibmngr/internal/utils"

	"github.com/dhowden/tag"
)

type ImportTask struct {
	Paths []string
	Records []file.LocalMetadata
	ReleaseCandidate []any
	BestMatch []any
}


func (t ImportTask) String() string {
 return fmt.Sprintf("ImportTask{\n  Paths: %v,\n  Records: %v,\n  ReleaseCandidate: %v,\n  BestMatch: %v\n}", t.Paths, t.Records, t.ReleaseCandidate, t.BestMatch)
}

// A struct representing a task, in order I want to
// * Split if differents albums
// * Fetch tag to detect actual artist, album and track number and years
// *

func (t *ImportTask) ExtractTags() {
	t.Records = make([]file.LocalMetadata, 0, len(t.Paths))
	for _, path := range t.Paths {
		tag, err := file.ExtractTag(path)
		if err != nil {
			log.Printf("failed to extract tag from %s: %v", path, err)
			continue
		}
		t.Records = append(t.Records, tag)
	}
}

func (t *ImportTask) GroupCommon() []*ImportTask {
	taskMap := make(map[string][]file.LocalMetadata)
	for _, record := range t.Records {
		recordKey := RecordKey(record)
		taskMap[recordKey] = append(taskMap[recordKey], record)
	}
	result := make([]*ImportTask, 0, len(taskMap))
	for _, tasks := range taskMap {
		paths := make([]string, 0, len(tasks))
		for _, task := range tasks {
			paths = append(paths, task.GetLocalPath())
		}
		result = append(result, &ImportTask{
			Paths: paths,
			Records: tasks,
		})
	}
	return result
}

func RecordKey(record tag.Metadata) string {
	artist := utils.NormalizeString(record.Artist())
	recordName := utils.NormalizeString(record.Album())
	return fmt.Sprintf("%s - %s", artist, recordName)
}
