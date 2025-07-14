package lib

import (
	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Genres struct {
	Moods  []Item `json:"moods"`
	Genres []Item `json:"genres"`
}

func parseGenres(raw string) Genres {

	j := gjson.Parse(raw)

	c := j.Get(
		"contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer",
	).Get(
		"content.sectionListRenderer.contents.#.gridRenderer",
	)

	m := c.Get("#(header.gridHeaderRenderer.title.runs.0.text == Moods & moments)")
	g := c.Get("#(header.gridHeaderRenderer.title.runs.0.text == Genres)")

	return Genres{
		Moods:  NavigationButton(m.Get("items")),
		Genres: NavigationButton(g.Get("items")),
	}
}

func GetGenres() (Genres, int) {

	context := utils.TypeBrowse("FEmusic_moods_and_genres", "", []string{})

	raw, status := utils.FetchBrowse(context)

	res := parseGenres(raw)

	return res, status
}
