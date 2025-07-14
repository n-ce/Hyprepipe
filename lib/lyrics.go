package lib

import (
	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Lyrics struct {
	Text   string `json:"text"`
	Source string `json:"source"`
}

func parseLyrics(raw string) Lyrics {

	j := gjson.Parse(raw)

	d := j.Get("contents.sectionListRenderer.contents.0.musicDescriptionShelfRenderer")

	l := d.Get("description")
	s := d.Get("footer")

	return Lyrics{
		Text:   RunsText(l),
		Source: RunsText(s),
	}
}

func GetLyrics(id string) (Lyrics, int) {

	context := utils.TypeBrowsePage(id, "lyrics")

	raw, status := utils.FetchBrowse(context)

	res := parseLyrics(raw)

	return res, status
}
