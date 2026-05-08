package importer

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"musiclibmngr/internal/utils"
	"os"
	"path/filepath"
	"sync"

	"github.com/dhowden/tag"
)

type Result[T any] struct {
	Value T
	Err   error
}

func Stage[I any, O any](
	ctx context.Context,
	in <-chan Result[I],
	workers int,
	buffer int,
	fn func(context.Context, I) (O, error),
) <-chan Result[O] {

	out := make(chan Result[O], buffer)

	var wg sync.WaitGroup
	wg.Add(workers)

	for range workers {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case item, ok := <-in:
					if !ok {
						return
					}

					if item.Err != nil {
						return
					}

					val, err := fn(ctx, item.Value)

					select {
					case out <- Result[O]{Value: val, Err: err}:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func GroupTaskByFolder(ctx context.Context, baseDir string) <-chan Result[ImportTask] {
	out := make(chan Result[ImportTask], 10)

	go func() {
		defer close(out)

		fileMap := make(map[string][]string)
		err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			isAudio, err := utils.IsAudio(path)
			if err != nil || !isAudio {
				return nil
			}

			rel, err := filepath.Rel(baseDir, path)
			if err != nil {
				return err
			}

			dir := filepath.Dir(rel)
			fileMap[dir] = append(fileMap[dir], path)

			return nil
		})

		if err != nil {
			out <- Result[ImportTask]{Err: err}
			return
		}

		for _, paths := range fileMap {
			task := ImportTask{Paths: paths}

			select {
			case out <- Result[ImportTask]{Value: task}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func extractTag(path string) (LocalMetadata, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	tags, err := tag.ReadFrom(fd)
	if err != nil {
		return nil, err
	}

	res := &localMetadata{
		Metadata:  tags,
		localPath: path,
	}
	return res, nil
}

func extractAllTags(ctx context.Context, task ImportTask) (ImportTask, error) {
	tracks := make([]LocalMetadata, 0, len(task.Paths))
	for _, p := range task.Paths {
		m, err := extractTag(p)
		if err != nil {
			log.Println(err)
		}
		tracks = append(tracks, m)
	}

	task.Records = tracks
	return task, nil
}

func SplitByTag(ctx context.Context, in <-chan Result[ImportTask], workers int, buffer int) <-chan Result[ImportTask] {
	out := make(chan Result[ImportTask], 10)
	var wg sync.WaitGroup
	wg.Add(workers)

	go func() {
		defer close(out)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return

			case item, ok := <-in:
				if !ok {
					return
				}

				taskMap := make(map[string][]LocalMetadata)
				for _, record := range item.Value.Records {
					recordKey := RecordKey(record)
					taskMap[recordKey] = append(taskMap[recordKey], record)
				}
				for _, tasks := range taskMap {
					paths := make([]string, 0, len(tasks))
					for _, task := range tasks {
						paths = append(paths, task.GetLocalPath())
					}
					select {
					case out <- Result[ImportTask]{
						Value: ImportTask{
							Paths:   paths,
							Records: tasks,
						}}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return out
}

func RecordKey(record LocalMetadata) string {
	artist := utils.NormalizeString(record.AlbumArtist())
	recordName := utils.NormalizeString(record.Album())
	return fmt.Sprintf("%s - %s", artist, recordName)
}

func GetRecordInfo(ctx context.Context, task ImportTask) (ImportTask, error) {
	if len(task.Records) > 1 {
		r := task.Records[0]
		_, tn := r.Track()
		task.ReleaseInfo = ReleaseInfo{
			Title:   r.Album(),
			Artist:  r.AlbumArtist(),
			TrackNb: tn,
			Year:    r.Year(),
		}
	} else {
		return task, fmt.Errorf("No album candidat")
	}
	return task, nil
}

func Run(ctx context.Context, baseDir string) {
	source := GroupTaskByFolder(ctx, baseDir)
	extracted := Stage(ctx, source, 4, 10, extractAllTags)
	split := SplitByTag(ctx, extracted, 4, 10)
	albums := Stage(ctx, split, 4, 10, GetRecordInfo)

	for res := range albums {
		if res.Err != nil {
			fmt.Println("error:", res.Err)
			continue
		}

		fmt.Println(res.Value)
	}
}
