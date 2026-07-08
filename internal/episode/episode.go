// Package episode loads the curated per-episode alcohol-reference dataset.
package episode

import (
	"encoding/json"
	"fmt"
	"os"
)

// StoreLink is a place where an alcohol can be bought.
type StoreLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Alcohol is a single alcohol reference spotted in an episode.
type Alcohol struct {
	JapaneseName string      `json:"japaneseName"`
	EnglishName  string      `json:"englishName"`
	Image        string      `json:"image"`
	Type         string      `json:"type"`
	ABV          string     `json:"abv"`
	Brand        string      `json:"brand"`
	Origin       string      `json:"origin"`
	JPStores     []StoreLink `json:"jpStores"`
	SGStores     []StoreLink `json:"sgStores"`
}

// Episode is one anime episode and every alcohol referenced in it.
type Episode struct {
	Number   int       `json:"number"`
	Title    string    `json:"title"`
	AirDate  string    `json:"airDate"`
	Alcohols []Alcohol `json:"alcohols"`
}

// Dataset is the full curated collection loaded from disk.
type Dataset struct {
	Series   string    `json:"series"`
	Episodes []Episode `json:"episodes"`
}

// Load reads and parses the episode dataset from path.
func Load(path string) (*Dataset, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("episode: read %s: %w", path, err)
	}

	var ds Dataset
	if err := json.Unmarshal(data, &ds); err != nil {
		return nil, fmt.Errorf("episode: parse %s: %w", path, err)
	}
	return &ds, nil
}

// ByNumber returns the episode with the given number, if present.
func (d *Dataset) ByNumber(n int) (Episode, bool) {
	for _, ep := range d.Episodes {
		if ep.Number == n {
			return ep, true
		}
	}
	return Episode{}, false
}
