package pathmatcher

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var placeholderRegex = regexp.MustCompile(`\{([^}]+)\}`)

type Matcher struct{}

func NewMatcher() *Matcher {
	return &Matcher{}
}

var ErrInvalidTemplate = errors.New("path does not match template")

func (m *Matcher) ExtractData(template, path string) (map[string]string, error) {
	pattern, keyMap := buildRegex(template)

	re, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return nil, fmt.Errorf("compile regex: %w", err)
	}

	matches := re.FindStringSubmatch(path)
	if matches == nil {
		return nil, ErrInvalidTemplate
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		result[keyMap[name]] = matches[i]
	}

	return result, nil
}

func buildRegex(template string) (string, map[string]string) {
	keyMap := make(map[string]string)

	pattern := placeholderRegex.ReplaceAllStringFunc(template, func(match string) string {
		content := placeholderRegex.FindStringSubmatch(match)[1]
		key := sanitizeKey(content)

		keyMap[key] = content
		return fmt.Sprintf("(?P<%s>[^/]+)", key)
	})

	return pattern, keyMap
}

func sanitizeKey(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, ":", "_")

	// keep only safe characters
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
