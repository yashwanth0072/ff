package search

import "strings"

// Tokenize splits a query into AND-ed search tokens on whitespace *and*
// forward slashes. Splitting on '/' matters because typing a real path
// like "/usr/share/sddm/theme" has no whitespace at all — without this,
// it would become a single token containing slashes, which can never
// match a plocate pattern token cleanly and can never match inside a
// single path component in Rank (see fuzzy.go), since no component
// itself contains a '/'. Splitting it into "usr", "share", "sddm",
// "theme" makes typing a path behave the same as typing separate words:
// every piece must independently match some component, in any order.
func Tokenize(query string) []string {
	return strings.FieldsFunc(query, func(r rune) bool {
		return r == '/' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
}
