// Package search implements fuzzy/partial matching over alcohol references.
package search

import (
	"sort"
	"strings"

	"github.com/modembcc/kamiina-lookup/internal/episode"
)

// Entry is a searchable alcohol reference tied to the episode it appeared in.
type Entry struct {
	EpisodeNumber int
	Alcohol       episode.Alcohol
}

// BuildEntries flattens every episode's alcohol references into one searchable list.
func BuildEntries(episodes []episode.Episode) []Entry {
	var entries []Entry
	for _, ep := range episodes {
		for _, a := range ep.Alcohols {
			entries = append(entries, Entry{EpisodeNumber: ep.Number, Alcohol: a})
		}
	}
	return entries
}

// Search returns entries matching query, ranked best-match first. It matches
// against Alcohol.JapaneseName, Alcohol.EnglishName, Alcohol.Brand, and
// Alcohol.Type, using exact/prefix/substring matching first and falling back
// to a fuzzy (edit-distance) match so typos and partial romanized input
// still hit.
func Search(entries []Entry, query string) []Entry {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return nil
	}

	type scoredEntry struct {
		entry Entry
		score int
	}

	var matches []scoredEntry
	for _, e := range entries {
		score := max(
			fieldScore(e.Alcohol.JapaneseName, query),
			fieldScore(e.Alcohol.EnglishName, query),
			fieldScore(e.Alcohol.Brand, query),
			fieldScore(e.Alcohol.Type, query),
		)
		if score > 0 {
			matches = append(matches, scoredEntry{entry: e, score: score})
		}
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].score > matches[j].score
	})

	results := make([]Entry, len(matches))
	for i, m := range matches {
		results[i] = m.entry
	}
	return results
}

// fieldScore rates how well query matches field: higher is better, 0 means
// no match at all (not even a fuzzy one). query must already be lowercased.
func fieldScore(field, query string) int {
	if field == "" {
		return 0
	}
	fieldLower := strings.ToLower(field)

	switch {
	case fieldLower == query:
		return 100
	case strings.HasPrefix(fieldLower, query):
		return 80
	case strings.Contains(fieldLower, query):
		return 60
	}

	if dist, ok := fuzzyDistance(fieldLower, query); ok {
		return max(50-dist*10, 10)
	}
	return 0
}

// fuzzyDistance reports the edit distance between query and the closest
// whitespace-separated token in field (or field as a whole), and whether
// that distance is within the tolerance for query's length.
func fuzzyDistance(field, query string) (int, bool) {
	tokens := strings.Fields(field)
	tokens = append(tokens, field)

	best := -1
	for _, tok := range tokens {
		d := levenshtein(tok, query)
		if best == -1 || d < best {
			best = d
		}
	}

	threshold := 1
	if len([]rune(query)) > 5 {
		threshold = 2
	}
	return best, best <= threshold
}

// levenshtein computes the edit distance between two strings, operating on
// runes so multi-byte characters (e.g. Japanese) count as single units.
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	if len(ra) == 0 {
		return len(rb)
	}
	if len(rb) == 0 {
		return len(ra)
	}

	prev := make([]int, len(rb)+1)
	curr := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(ra); i++ {
		curr[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			curr[j] = min(
				prev[j]+1,      // deletion
				curr[j-1]+1,    // insertion
				prev[j-1]+cost, // substitution
			)
		}
		prev, curr = curr, prev
	}
	return prev[len(rb)]
}
