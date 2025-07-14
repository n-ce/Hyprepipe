package lib

import (
	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Options struct {
	Default string `json:"default"`
	All     []Item `json:"all"`
}

type Charts struct {
	Options  Options `json:"options"`
	Artists  []Item  `json:"artists"`
	Trending []Item  `json:"trending"`
}

func parseCharts(raw string) Charts {

	j := gjson.Parse(raw)

	c := j.Get("contents.singleColumnBrowseResultsRenderer.tabs.0").Get("tabRenderer.content.sectionListRenderer.contents")

	o := c.Get(
		"0.musicShelfRenderer.subheaders.0.musicSideAlignedItemRenderer",
	).Get(
		"startItems.0.musicSortFilterButtonRenderer",
	)

	a := c.Get(
		"#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Top artists)",
	).Get("musicCarouselShelfRenderer.contents")

	t := c.Get(
		"#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Trending)",
	).Get("musicCarouselShelfRenderer.contents")

	opts := o.Get(
		"menu.musicMultiSelectMenuRenderer.options.#.musicMultiSelectMenuItemRenderer",
	)
	ref := j.Get("frameworkUpdates.entityBatchUpdate.mutations")

	return Charts{
		Options: Options{
			Default: RunsText(o.Get("title")),
			All:     MultiSelectMenuItemRenderer(opts, ref),
		},
		Artists:  ResponsiveListItemRendererCH(a),
		Trending: ResponsiveListItemRenderer(t),
	}
}

func GetCharts(params, code string) (Charts, int) {

	context := utils.TypeBrowseForm("FEmusic_charts", params, code)

	raw, status := utils.FetchBrowse(context)

	res := parseCharts(raw)

	return res, status
}
