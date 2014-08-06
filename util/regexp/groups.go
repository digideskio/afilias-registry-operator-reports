package regexp

import (
	goregexp "regexp"
)

// Maps named groups defined in r to the matches found in s
func MapRegexpGroups(r *goregexp.Regexp, s string) map[string]string {
	captures := make(map[string]string)
	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}
	for i, name := range r.SubexpNames() {
		if i == 0 {
			continue
		}
		captures[name] = match[i]

	}
	return captures
}
