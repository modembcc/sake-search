package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/modembcc/sake-search/internal/episode"
	"github.com/modembcc/sake-search/internal/search"
)

func main() {
	dataPath := flag.String("data", "data/episodes.json", "path to episodes JSON dataset")
	episodeNum := flag.Int("episode", 0, "episode number to look up")
	query := flag.String("query", "", "search all episodes for a drink by name or brand")
	flag.Parse()

	if *episodeNum == 0 && *query == "" {
		fmt.Println("usage: sake-search --episode N [--data path/to/episodes.json]")
		fmt.Println("       sake-search --query TEXT [--data path/to/episodes.json]")
		os.Exit(1)
	}

	dataset, err := episode.Load(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading dataset:", err)
		os.Exit(1)
	}

	if *query != "" {
		runSearch(dataset, *query)
		return
	}

	ep, ok := dataset.ByNumber(*episodeNum)
	if !ok {
		fmt.Fprintf(os.Stderr, "episode %d not found\n", *episodeNum)
		os.Exit(1)
	}

	printEpisode(ep)
}

func runSearch(dataset *episode.Dataset, query string) {
	entries := search.BuildEntries(dataset.Episodes)
	results := search.Search(entries, query)

	if len(results) == 0 {
		fmt.Printf("No matches for %q.\n", query)
		return
	}

	fmt.Printf("%d match(es) for %q:\n", len(results), query)
	for _, r := range results {
		fmt.Printf("\nEpisode %d\n", r.EpisodeNumber)
		printAlcohol(r.Alcohol)
	}
}

func printEpisode(ep episode.Episode) {
	fmt.Printf("Episode %d", ep.Number)
	if ep.Title != "" {
		fmt.Printf(": %s", ep.Title)
	}
	fmt.Println()
	if ep.AirDate != "" {
		fmt.Printf("Aired: %s\n", ep.AirDate)
	}

	if len(ep.Alcohols) == 0 {
		fmt.Println("No alcohol references recorded for this episode.")
		return
	}

	for _, a := range ep.Alcohols {
		printAlcohol(a)
	}
}

func printAlcohol(a episode.Alcohol) {
	fmt.Printf("%s", a.EnglishName)
	if a.JapaneseName != "" {
		fmt.Printf(" (%s)", a.JapaneseName)
	}
	fmt.Println()
	fmt.Printf("    Type:   %s\n", a.Type)
	if a.Brand != "" {
		fmt.Printf("    Brand:  %s\n", a.Brand)
	}
	if a.ABV != "" {
		fmt.Printf("    ABV:    %s\n", a.ABV)
	}
	if a.Origin != "" {
		fmt.Printf("    Origin: %s\n", a.Origin)
	}
	if a.Image != "" {
		fmt.Printf("    Image:  %s\n", a.Image)
	}
	for _, s := range a.JPStores {
		if s.URL == "" {
			continue
		}
		fmt.Printf("    JP store: %s (%s)\n", s.Name, s.URL)
	}
	for _, s := range a.SGStores {
		if s.URL == "" {
			continue
		}
		fmt.Printf("    SG store: %s (%s)\n", s.Name, s.URL)
	}
}
