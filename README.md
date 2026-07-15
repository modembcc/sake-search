# kamiina-lookup

Terminal app that looks up every real-world alcohol referenced in a given episode of the anime _Botan Kamiina Fully Blossoms When Drunk_ (上伊那ぼたん、酔へる姿は百合の花) — Japanese and English names, ABV, brand, origin, and where to buy it in Japan and Singapore. Built with Go.

## Usage

```
go run . -episode 1
```

```
Episode 1
Aired: 2026-04-11

[1] Whiskey Highball (ハイボール)
    Type:   Cocktail
    ABV:    7% to 9% in ready-to-drink cans
    Origin: Japan
...
```

## Data

Curated by hand in [data/episodes.json](data/episodes.json) — one entry per episode, each holding an array of alcohol references (an episode can feature more than one drink). Schema lives in [internal/episode/episode.go](internal/episode/episode.go):

- `japaneseName`, `englishName` — native and romanized/English name
- `image` — reference image URL
- `type` — sake, beer, wine, cocktail, etc.
- `abv` — alcohol content
- `brand` — producer or product line
- `origin` — region the drink comes from
- `jpStores` / `sgStores` — arrays of `{name, url}` retailer links

## Status

Work in progress:

- [x] Episode lookup CLI (`-episode N`)
- [x] Episodes 1–12
- [x] Fuzzy/partial search across all referenced drinks
- [x] Stylized UI
- [ ] JP/EN display toggle
