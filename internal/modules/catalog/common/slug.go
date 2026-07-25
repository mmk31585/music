package common

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
)

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// persianToLatin maps Persian/Arabic characters to their Latin equivalents
// for slug generation on a Persian music platform.
var persianToLatin = strings.NewReplacer(
	// Persian-specific
	"پ", "p",
	"چ", "ch",
	"ژ", "zh",
	"گ", "g",
	"ک", "k",
	"ی", "y",
	// Arabic-origin letters used in Persian
	"ا", "a",
	"آ", "a",
	"أ", "a",
	"إ", "a",
	"ب", "b",
	"ت", "t",
	"ث", "s",
	"ج", "j",
	"ح", "h",
	"خ", "kh",
	"د", "d",
	"ذ", "z",
	"ر", "r",
	"ز", "z",
	"س", "s",
	"ش", "sh",
	"ص", "s",
	"ض", "z",
	"ط", "t",
	"ظ", "z",
	"ع", "a",
	"غ", "gh",
	"ف", "f",
	"ق", "gh",
	"ل", "l",
	"م", "m",
	"ن", "n",
	"و", "v",
	"ه", "h",
	"ة", "h",
	"ء", "",
	"ـ", "",
	"ئ", "y",
	"ؤ", "v",
	"ۀ", "h",
	"ھ", "h",
	// Also handle lowercase Persian (which are the same in Persian script)
	"پ", "p",
	"چ", "ch",
	"ژ", "zh",
	"گ", "g",
	"ک", "k",
	"ی", "y",
	"ا", "a",
	"آ", "a",
	"ب", "b",
	"ت", "t",
	"ث", "s",
	"ج", "j",
	"ح", "h",
	"خ", "kh",
	"د", "d",
	"ذ", "z",
	"ر", "r",
	"ز", "z",
	"س", "s",
	"ش", "sh",
	"ص", "s",
	"ض", "z",
	"ط", "t",
	"ظ", "z",
	"ع", "a",
	"غ", "gh",
	"ف", "f",
	"ق", "gh",
	"ل", "l",
	"م", "m",
	"ن", "n",
	"و", "v",
	"ه", "h",
	"ة", "h",
	"ء", "",
	"ـ", "",
	"ئ", "y",
	"ؤ", "v",
)

// Slugify converts a string to a URL-friendly slug.
// It transliterates Persian characters to Latin, then strips non-alphanumeric
// characters. If the result is empty, it falls back to a short hash of the
// original input.
func Slugify(s string) string {
	original := strings.TrimSpace(s)
	if original == "" {
		return "item"
	}

	// Transliterate Persian/Arabic characters to Latin
	result := persianToLatin.Replace(original)

	result = strings.ToLower(result)
	result = nonAlnum.ReplaceAllString(result, "-")
	result = strings.Trim(result, "-")

	// If the slug is still empty or too short (e.g., all special characters),
	// use a short hash of the original input for uniqueness.
	if result == "" || len(result) < 2 {
		h := sha256.Sum256([]byte(original))
		return fmt.Sprintf("item-%x", h[:4])
	}

	return result
}
