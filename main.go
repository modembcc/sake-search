package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/modembcc/kamiina-lookup/internal/episode"
	"github.com/modembcc/kamiina-lookup/internal/search"
	"github.com/modembcc/kamiina-lookup/internal/style"
)

func main() {
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = printUsage

	dataPath := flag.String("data", "data/episodes.json", "path to episodes JSON dataset")
	episodeNum := flag.Int("episode", 0, "episode number to look up")
	query := flag.String("query", "", "search all episodes for a drink by name or brand")
	flag.Parse()

	if *episodeNum == 0 && *query == "" {
		printUsage()
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

func printUsage() {
	fmt.Println(style.Banner("KAMIINA LOOKUP"))
	fmt.Println(style.Muted("  Alcohol lookup for Botan Kamiina Fully Blossoms When Drunk"))
	fmt.Println()

	fmt.Println(style.Heading("USAGE"))
	fmt.Printf("  %s --episode N [--data path/to/episodes.json]\n", style.Accent("kamiina-lookup"))
	fmt.Printf("  %s --query TEXT [--data path/to/episodes.json]\n", style.Accent("kamiina-lookup"))
	fmt.Println()

	fmt.Println(style.Heading("FLAGS"))
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("  %s\n", style.Label("-"+f.Name))
		fmt.Printf("      %s\n", f.Usage)
	})
}

func runSearch(dataset *episode.Dataset, query string) {
	entries := search.BuildEntries(dataset.Episodes)
	results := search.Search(entries, query)

	if len(results) == 0 {
		fmt.Println(style.Muted(fmt.Sprintf("No matches for %q.", query)))
		return
	}

	fmt.Println(style.Title(fmt.Sprintf("%d match(es) for %q:", len(results), query)))
	for _, r := range results {
		fmt.Println()
		fmt.Println(style.Divider(40))
		fmt.Println(style.Muted(fmt.Sprintf("Episode %d", r.EpisodeNumber)))
		printAlcohol(r.Alcohol)
	}
}

func printEpisode(ep episode.Episode) {
	title := fmt.Sprintf("Episode %d", ep.Number)
	if ep.Title != "" {
		title += ": " + ep.Title
	}
	fmt.Println(style.Title(title))
	if ep.AirDate != "" {
		fmt.Println(style.Muted("Aired: " + ep.AirDate))
	}

	if len(ep.Alcohols) == 0 {
		fmt.Println(style.Muted("No alcohol references recorded for this episode."))
		return
	}

	for _, a := range ep.Alcohols {
		fmt.Println()
		fmt.Println(style.Divider(40))
		printAlcohol(a)
	}
}

func printAlcohol(a episode.Alcohol) {
	fmt.Printf("%s %s", style.Bullet(), style.Heading(a.EnglishName))
	if a.JapaneseName != "" {
		fmt.Printf(" %s", style.Accent("("+a.JapaneseName+")"))
	}
	fmt.Println()

	printField("Type", a.Type)
	printField("Brand", a.Brand)
	printField("ABV", a.ABV)
	printField("Origin", a.Origin)
	printField("Image", a.Image)

	for _, s := range a.JPStores {
		if s.URL == "" {
			continue
		}
		printField("JP store", fmt.Sprintf("%s (%s)", s.Name, s.URL))
	}
	for _, s := range a.SGStores {
		if s.URL == "" {
			continue
		}
		printField("SG store", fmt.Sprintf("%s (%s)", s.Name, s.URL))
	}
}

func printField(label, value string) {
	if value == "" {
		return
	}
	fmt.Printf("    %s %s\n", style.Label(label+":"), value)
}
