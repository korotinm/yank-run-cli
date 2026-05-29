package ids

import "regexp"

var idRegex = regexp.MustCompile(`^[0-9a-f]{64}$`)

// IsValidID reports whether s is a 64-char lowercase hex snippet id.
func IsValidID(s string) bool {
	return idRegex.MatchString(s)
}

// IsNumericIndex reports whether s is a non-negative integer index < 1000.
// Used to disambiguate REF arguments (id vs. last-search index).
func IsNumericIndex(s string) bool {
	if s == "" || len(s) > 3 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// ShortPrefix returns the 8-char id prefix shown in search results.
func ShortPrefix(id string) string {
	if len(id) < 8 {
		return id
	}
	return id[:8]
}
