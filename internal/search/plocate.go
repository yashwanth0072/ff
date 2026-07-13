package search

import (
	"bufio"
	"errors"
	"os/exec"
	"strings"
)

// ErrNoPlocate is returned when the plocate binary isn't on PATH.
var ErrNoPlocate = errors.New("plocate not found on PATH (install it: sudo pacman -S plocate)")

// MaxCandidates caps how many paths we pull from plocate before fuzzy
// ranking. plocate is fast but we don't want to fuzzy-score 500k paths
// on every keystroke.
const MaxCandidates = 4000

// Candidates runs plocate using the query's whitespace-separated tokens
// as AND-ed patterns (plocate/mlocate semantics: multiple patterns must
// all match), which acts as a cheap prefilter. The caller is expected to
// fuzzy-rank the returned paths against the full query afterwards.
func Candidates(query string) ([]string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}

	if _, err := exec.LookPath("plocate"); err != nil {
		return nil, ErrNoPlocate
	}

	tokens := Tokenize(query)
	if len(tokens) == 0 {
		return nil, nil
	}

	args := []string{"-i", "-l", "1500"}
	args = append(args, tokens...)

	cmd := exec.Command("plocate", args...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var paths []string
	scanner := bufio.NewScanner(out)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		paths = append(paths, line)
		if len(paths) >= MaxCandidates {
			break
		}
	}

	// plocate exits with status 1 when there are simply no matches.
	// That's not an application error, so we swallow it as long as we
	// got a clean scan.
	_ = cmd.Wait()

	if scanErr := scanner.Err(); scanErr != nil {
		return paths, scanErr
	}

	return paths, nil
}
