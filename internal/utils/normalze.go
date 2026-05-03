package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var stopWords = map[string]struct{}{
	"the": {},
	"a":   {},
	"an":  {},
}

var (
	reFeat       = regexp.MustCompile(`(?i)\b(feat|ft)\.?\b`)
	reNonAlnum   = regexp.MustCompile(`[^a-z0-9\s]`)
	reMultiSpace = regexp.MustCompile(`\s+`)
)

func NormalizeString(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "&", " and ")
	s = removeAccents(s)
	s = reFeat.ReplaceAllString(s, " ")
	s = reNonAlnum.ReplaceAllString(s, " ")
	words := strings.Fields(s)
	filtered := make([]string, 0, len(words))
	for _, w := range words {
		if _, ok := stopWords[w]; !ok {
			filtered = append(filtered, w)
		}
	}
	s = strings.Join(filtered, " ")
	s = reMultiSpace.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

func removeAccents(s string) string {
	t := norm.NFD.String(s)
	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue // skip accent marks
		}
		b.WriteRune(r)
	}
	return b.String()
}
