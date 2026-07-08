// Package search implements fuzzy/partial matching over alcohol references.
package search

import "github.com/modembcc/sake-search/internal/episode"

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

// Search returns entries matching query, ranked best-match first.
//
// TODO(you): hand-write this. Match against Alcohol.JapaneseName,
// Alcohol.EnglishName, and Alcohol.Brand. Start with substring matching, then
// layer in fuzzy tolerance (e.g. edit distance or a subsequence match) so
// typos and partial romanized input still hit.
func Search(entries []Entry, query string) []Entry {
	panic("not implemented")
}
