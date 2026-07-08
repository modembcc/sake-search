package search

import (
	"testing"

	"github.com/modembcc/sake-search/internal/episode"
)

func testEntries() []Entry {
	return []Entry{
		{EpisodeNumber: 1, Alcohol: episode.Alcohol{JapaneseName: "獺祭", EnglishName: "Dassai", Brand: "Asahi Shuzo"}},
		{EpisodeNumber: 2, Alcohol: episode.Alcohol{JapaneseName: "十四代", EnglishName: "Juyondai", Brand: "Takagi Shuzo"}},
		{EpisodeNumber: 3, Alcohol: episode.Alcohol{JapaneseName: "而今", EnglishName: "Jikon", Brand: "Kiyamasho Shuzo"}},
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
