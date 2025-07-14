package lib

import (
	"net/url"

	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

func GetAlbumUrl(id string) (map[string]string, int) {
	ctx := utils.TypeBrowsePage(id, "album")
	raw, status := utils.FetchBrowse(ctx)

	uri := gjson.Parse(raw).Get(
		"microformat.microformatDataRenderer.urlCanonical",
	).String()

	u, err := url.Parse(uri)
	if err != nil {
		return map[string]string{
			"error":   "500",
			"message": "failed to parse url",
		}, 500
	}

	return map[string]string{
		"id": u.Query().Get("list"),
	}, status
}
