package search

import (
	"strings"

	"github.com/sahilm/fuzzy"
)

// Result is a ranked path with the index positions that matched, so the
// TUI can highlight them.
type Result struct {
	Path           string
	MatchedIndexes []int // indexes into Path (deduped, sorted)
	Score          int
}

// component is one slash-separated piece of a path, plus the byte offset
// where it starts in the original path string (so a match found inside
// the component can be translated back to a whole-path index).
type component struct {
	text   string
	offset int
}

func splitComponents(path string) []component {
	var comps []component
	start := 0
	for i := 0; i <= len(path); i++ {
		if i == len(path) || path[i] == '/' {
			if i > start {
				comps = append(comps, component{text: path[start:i], offset: start})
			}
			start = i + 1
		}
	}
	return comps
}

// Rank fuzzy-scores candidate paths against query. The query is split on
// whitespace into tokens (mirroring plocate's AND-of-patterns semantics);
// every token must fuzzy-match at least one path component for the
// candidate to be included.
//
// Each token is matched against individual path *components* (each
// directory name and the filename, matched separately) rather than the
// full concatenated path string. Two reasons:
//
//  1. sahilm/fuzzy penalizes matches by their distance from the start of
//     the string. On an absolute path like
//     /home/yashwanth/.config/caelestia/modules/notes/NotesDrawer.qml,
//     even a perfect match on ".qml" sits ~60 characters in and can score
//     negative despite being a real, contiguous match. Matching against
//     each short component instead keeps that positional penalty bounded
//     to the component's own length, so deep-but-real matches don't get
//     scored as if they were noise.
//  2. Fuzzy subsequence matching against one giant string lets unrelated
//     letters from *different* path segments combine into a false
//     positive -- e.g. token "notes" can find its n-o-t-e-s scattered
//     across "yashwanth", "Projects", "config.json" even though none of
//     those segments actually relate to notes. Requiring the whole token
//     to match within a single component closes that gap.
func Rank(candidates []string, query string) []Result {
	query = strings.TrimSpace(query)
	tokens := Tokenize(query)
	if len(tokens) == 0 || len(candidates) == 0 {
		return nil
	}

	compsByPath := make([][]component, len(candidates))
	textsByPath := make([][]string, len(candidates))
	for i, c := range candidates {
		compsByPath[i] = splitComponents(c)
		texts := make([]string, len(compsByPath[i]))
		for j, comp := range compsByPath[i] {
			texts[j] = comp.text
		}
		textsByPath[i] = texts
	}

	type acc struct {
		score   int
		indexes map[int]bool
		ok      bool
	}
	accs := make([]acc, len(candidates))
	for i := range accs {
		accs[i] = acc{indexes: map[int]bool{}, ok: true}
	}

	for _, tok := range tokens {
		for i, texts := range textsByPath {
			if !accs[i].ok {
				continue // already disqualified by an earlier token
			}
			if len(texts) == 0 {
				accs[i].ok = false
				continue
			}

			matches := fuzzy.Find(tok, texts)
			if len(matches) == 0 {
				accs[i].ok = false
				continue
			}

			// Take the best-scoring component for this token. The last
			// component is the basename (filename) -- weight a basename
			// hit heavier so "notes" ranks NotesDrawer.qml above a path
			// that only matches somewhere in a parent directory name.
			best := matches[0]
			for _, m := range matches[1:] {
				if m.Score > best.Score {
					best = m
				}
			}
			comp := compsByPath[i][best.Index]
			isBasename := best.Index == len(texts)-1

			score := best.Score
			if isBasename {
				score *= 2
			}
			accs[i].score += score
			for _, idx := range best.MatchedIndexes {
				accs[i].indexes[comp.offset+idx] = true
			}
		}
	}

	results := make([]Result, 0, len(candidates))
	for i, a := range accs {
		if !a.ok {
			continue
		}
		idxs := make([]int, 0, len(a.indexes))
		for idx := range a.indexes {
			idxs = append(idxs, idx)
		}
		// simple insertion sort, these slices are tiny
		for x := 1; x < len(idxs); x++ {
			for y := x; y > 0 && idxs[y] < idxs[y-1]; y-- {
				idxs[y], idxs[y-1] = idxs[y-1], idxs[y]
			}
		}
		results = append(results, Result{
			Path:           candidates[i],
			MatchedIndexes: idxs,
			Score:          a.score,
		})
	}

	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Score > results[j-1].Score; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}

	return results
}
