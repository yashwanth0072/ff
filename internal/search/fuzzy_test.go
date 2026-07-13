package search

import "testing"

func TestRankPrefersBasenameMatch(t *testing.T) {
	candidates := []string{
		"/home/yashwanth/.config/caelestia/modules/notes/NotesDrawer.qml",
		"/home/yashwanth/random/deeply/nested/path/that/happens/to/contain/n/o/t/e/s",
		"/usr/share/caelestia/config.json",
	}

	results := Rank(candidates, "notes qml")
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if results[0].Path != candidates[0] {
		t.Fatalf("expected NotesDrawer.qml to rank first, got %s", results[0].Path)
	}
}

func TestRankEmptyQuery(t *testing.T) {
	if got := Rank([]string{"a", "b"}, ""); got != nil {
		t.Fatalf("expected nil for empty query, got %v", got)
	}
}

func TestRankAllTokensMustMatch(t *testing.T) {
	candidates := []string{
		"/home/yashwanth/Projects/caelestia-shell/README.md",
		"/home/yashwanth/Documents/notes.txt",
	}
	results := Rank(candidates, "cael notes")
	if len(results) != 0 {
		t.Fatalf("expected no results when tokens don't jointly match a single candidate, got %v", results)
	}
}

// Regression test: a real match sitting deep in a long absolute path
// used to get thrown away because sahilm/fuzzy's raw score goes negative
// the further a match sits from the start of the string. Per-component
// matching keeps the positional penalty bounded to the component itself.
func TestRankDeepMatchNotDropped(t *testing.T) {
	candidates := []string{
		"/home/yashwanth/.config/caelestia/modules/notes/NotesDrawer.qml",
	}
	results := Rank(candidates, "qml")
	if len(results) != 1 {
		t.Fatalf("expected the deep .qml match to survive, got %d results", len(results))
	}
}

// Regression test: fuzzy-matching a token against the whole concatenated
// path let letters from unrelated segments combine into a false
// positive (e.g. "notes" scavenging n/o/t/e/s across "yashwanth",
// "Projects", "config.json"). Per-component matching requires the full
// token to land inside a single segment.
func TestRankRejectsCrossComponentScatter(t *testing.T) {
	candidates := []string{
		"/home/yashwanth/Projects/caelestia-shell/README.md",
		"/usr/share/caelestia/config.json",
	}
	results := Rank(candidates, "notes")
	if len(results) != 0 {
		t.Fatalf("expected no matches for a token scattered across unrelated segments, got %v", results)
	}
}

// Regression test: typing an actual path (no spaces, only slashes) has
// to work the same as typing its components as separate words, since
// that's the natural way to search once you already half-remember
// where something lives.
func TestRankMatchesTypedPath(t *testing.T) {
	candidates := []string{
		"/usr/share/sddm/theme/candy-machine",
		"/usr/share/fonts/noto",
	}
	results := Rank(candidates, "/usr/share/sddm/theme")
	if len(results) != 1 {
		t.Fatalf("expected typing a path to match, got %d results", len(results))
	}
	if results[0].Path != candidates[0] {
		t.Fatalf("wrong match: %s", results[0].Path)
	}
}

func TestRankSingleToken(t *testing.T) {
	candidates := []string{
		"/home/yashwanth/.config/caelestia/modules/notes/NotesDrawer.qml",
		"/etc/hostname",
	}
	results := Rank(candidates, "qml")
	if len(results) != 1 {
		t.Fatalf("expected exactly 1 match, got %d", len(results))
	}
	if results[0].Path != candidates[0] {
		t.Fatalf("wrong match: %s", results[0].Path)
	}
}
