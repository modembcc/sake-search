package search

import (
	"testing"

	"github.com/modembcc/kamiina-lookup/internal/episode"
)

func testEntries() []Entry {
	return []Entry{
		{EpisodeNumber: 1, Alcohol: episode.Alcohol{JapaneseName: "獺祭", EnglishName: "Dassai", Brand: "Asahi Shuzo", Type: "Sake"}},
		{EpisodeNumber: 2, Alcohol: episode.Alcohol{JapaneseName: "十四代", EnglishName: "Juyondai", Brand: "Takagi Shuzo", Type: "Sake"}},
		{EpisodeNumber: 3, Alcohol: episode.Alcohol{JapaneseName: "而今", EnglishName: "Jikon", Brand: "Kiyamasho Shuzo", Type: "Sake"}},
		{EpisodeNumber: 5, Alcohol: episode.Alcohol{EnglishName: "Salty Dog", Type: "Cocktail"}},
	}
}

func containsEnglishName(entries []Entry, name string) bool {
	for _, e := range entries {
		if e.Alcohol.EnglishName == name {
			return true
		}
	}
	return false
}

func TestSearch_ExactMatch(t *testing.T) {
	got := Search(testEntries(), "獺祭")
	if !containsEnglishName(got, "Dassai") {
		t.Errorf("expected exact match on 獺祭 to return Dassai, got %+v", got)
	}
}

func TestSearch_PartialMatch(t *testing.T) {
	got := Search(testEntries(), "十四")
	if !containsEnglishName(got, "Juyondai") {
		t.Errorf("expected partial match on 十四 to return Juyondai, got %+v", got)
	}
}

func TestSearch_MatchesOnBrand(t *testing.T) {
	got := Search(testEntries(), "Takagi Shuzo")
	if !containsEnglishName(got, "Juyondai") {
		t.Errorf("expected brand match on Takagi Shuzo to return Juyondai, got %+v", got)
	}
}

func TestSearch_NoMatchReturnsEmpty(t *testing.T) {
	got := Search(testEntries(), "zzzznonexistent")
	if len(got) != 0 {
		t.Errorf("expected no matches, got %+v", got)
	}
}

func TestSearch_FuzzyTypoMatch(t *testing.T) {
	got := Search(testEntries(), "Dasai")
	if !containsEnglishName(got, "Dassai") {
		t.Errorf("expected fuzzy match on typo 'Dasai' to return Dassai, got %+v", got)
	}
}

func TestSearch_RanksExactMatchAboveSubstringMatch(t *testing.T) {
	entries := []Entry{
		{EpisodeNumber: 1, Alcohol: episode.Alcohol{EnglishName: "Dassai Blue"}},
		{EpisodeNumber: 2, Alcohol: episode.Alcohol{EnglishName: "Dassai"}},
	}
	got := Search(entries, "Dassai")
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d: %+v", len(got), got)
	}
	if got[0].Alcohol.EnglishName != "Dassai" {
		t.Errorf("expected exact match 'Dassai' ranked first, got %+v", got)
	}
}

func TestSearch_IsCaseInsensitive(t *testing.T) {
	got := Search(testEntries(), "dassai")
	if !containsEnglishName(got, "Dassai") {
		t.Errorf("expected case-insensitive match on 'dassai' to return Dassai, got %+v", got)
	}
}

func TestSearch_MatchesOnType(t *testing.T) {
	got := Search(testEntries(), "cocktail")
	if !containsEnglishName(got, "Salty Dog") {
		t.Errorf("expected type match on 'cocktail' to return Salty Dog, got %+v", got)
	}
	if containsEnglishName(got, "Dassai") {
		t.Errorf("expected type match on 'cocktail' to exclude sake entries, got %+v", got)
	}
}
