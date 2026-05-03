package importer

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestImportTask_GroupCommon(t *testing.T) {
	// Use filepath.Glob to extract paths matching a pattern
	pattern := "../../testdata/ACDC/*"
	path, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("failed to glob paths: %v", err)
	}

	task := &ImportTask{
		Paths: path,
	}
	task.ExtractTags()
	grp := task.GroupCommon()
	fmt.Println(grp)
}
