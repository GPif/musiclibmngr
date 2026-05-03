package importer

type ImportTask struct {
	Paths []string
	Records []any
	ReleaseCandidate []any
	BestMatch []any
}
