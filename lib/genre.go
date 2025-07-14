package lib

import (
	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Genre struct {
	Title     string                 `json:"title"`
	Spotlight []Item                 `json:"spotlight"`
	Featured  []Item                 `json:"featured"`
	Community []Item                 `json:"community"`
	Shelf     map[string]interface{} `json:"shelf"`
}

func parseGenre(raw string) Genre {
	j := gjson.Parse(raw)

	c := j.Get(
		"contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents",
	)

	g := c.Get("#.gridRenderer")

	s := g.Get("#(header.gridHeaderRenderer.title.runs.0.text == Spotlight)")
	f := g.Get("#(header.gridHeaderRenderer.title.runs.0.text == Featured playlists)")
	cp := g.Get("#(header.gridHeaderRenderer.title.runs.0.text == Community playlists)")

	return Genre{
		Title:     RunsText(j.Get("header.musicHeaderRenderer.title")),
		Spotlight: TwoRowItemRenderer(s.Get("items"), false),
		Featured:  TwoRowItemRenderer(f.Get("items"), false),
		Community: TwoRowItemRenderer(cp.Get("items"), false),
		Shelf: ShelfRenderer(
			c.Get("#.musicCarouselShelfRenderer"),
		),
	}
}

func GetGenre(param string) (Genre, int) {

	context := utils.TypeBrowse("FEmusic_moods_and_genres_category", param, []string{})

	raw, status := utils.FetchBrowse(context)

	res := parseGenre(raw)

	return res, status
}
