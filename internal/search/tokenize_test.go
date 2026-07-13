package search

import "testing"

func TestTokenizeSplitsOnSlashAndSpace(t *testing.T) {
	cases := map[string][]string{
		"/usr/share/sddm/theme": {"usr", "share", "sddm", "theme"},
		"cael notes":            {"cael", "notes"},
		"cael/notes drawer":     {"cael", "notes", "drawer"},
		"  spaced   out  ":      {"spaced", "out"},
		"":                      {},
		"///":                   {},
	}
	for input, want := range cases {
		got := Tokenize(input)
		if len(got) != len(want) {
			t.Errorf("Tokenize(%q) = %v, want %v", input, got, want)
			continue
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("Tokenize(%q) = %v, want %v", input, got, want)
				break
			}
		}
	}
}
