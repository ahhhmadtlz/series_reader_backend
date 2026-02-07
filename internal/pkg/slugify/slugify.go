package slugify

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9-]+`)
	multipleHyphensRegex=regexp.MustCompile(`-+`)
)

func Make(s string)string{
	s=strings.ToLower(s)

	t:=transform.Chain(norm.NFD,runes.Remove(runes.In(unicode.Mn)),norm.NFC)
	s,_,_=transform.String(t,s)

	s=strings.ReplaceAll(s," ","-")

  s = nonAlphanumericRegex.ReplaceAllString(s, "")

	// Replace multiple consecutive hyphens with a single hyphen
	s = multipleHyphensRegex.ReplaceAllString(s, "-")

	// Trim hyphens from both ends
	s = strings.Trim(s, "-")

	return s

}

func MakeUnique(s string, existsFunc func(string) bool) string {
	slug := Make(s)
	if !existsFunc(slug) {
		return slug
	}

	// If slug exists, append a number
	counter := 2
	for {
		uniqueSlug := slug + "-" + strconv.Itoa(counter)
		if !existsFunc(uniqueSlug) {
			return uniqueSlug
		}
		counter++
		if counter > 9999 {
			// Fallback to avoid infinite loop
			return slug + "-" + strconv.Itoa(counter)
		}
	}
}