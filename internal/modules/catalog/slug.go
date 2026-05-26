package catalog

import (
	"regexp"
	"strings"
)

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)
var multiDash = regexp.MustCompile(`-+`)

func makeSlug(input string) string {
	s := strings.ToLower(strings.TrimSpace(input))
	s = nonSlugChars.ReplaceAllString(s, "-")
	s = multiDash.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	if s == "" {
		return "item"
	}

	return s
}
