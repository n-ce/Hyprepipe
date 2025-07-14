package lib

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/tidwall/gjson"
)

type Thumbnail struct {
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
	Url    string `json:"url"`
}

type Item struct {
	Id         string      `json:"id"`
	Title      string      `json:"title"`
	Sub        string      `json:"subtitle,omitempty"`
	SubId      string      `json:"subId,omitempty"`
	Thumbnails []Thumbnail `json:"thumbnails,omitempty"`
}

func RunsText(j gjson.Result) (s string) {
	a := j.Get("runs.#.text").Array()
	for i := 0; i < len(a); i++ {
		s += a[i].String()
	}
	return s
}

func parseUrl(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	prehost := u.Host
	u.Scheme = "https"
	u.Host = os.Getenv("HYP_PROXY")

	q := u.Query()
	q.Set("host", prehost)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func GetThumbnails(j gjson.Result) []Thumbnail {
	size := j.Get("#").Int()
	t := make([]Thumbnail, size)

	wg := sync.WaitGroup{}
	wg.Add(int(size))

	j.ForEach(
		func(n, j gjson.Result) bool {
			go func(i int64, v gjson.Result) {
				defer wg.Done()

				u, err := parseUrl(v.Get("url").String())
				if err != nil {
					log.Fatal(err)
				}

				t[i] = Thumbnail{
					Url:    u,
					Width:  v.Get("width").Int(),
					Height: v.Get("height").Int(),
				}
			}(n.Int(), j)
			return true
		},
	)

	wg.Wait()

	return t
}

func TwoRowItemRenderer(a gjson.Result, t bool) []Item {

	id := "menu.menuRenderer.items.#(menuNavigationItemRenderer.text.runs.0.text == Shuffle play).menuNavigationItemRenderer.navigationEndpoint.watchPlaylistEndpoint.playlistId"

	if !t {
		id = "navigationEndpoint.browseEndpoint.browseId"
	}

	v := a.Get("#.musicTwoRowItemRenderer")
	size := v.Get("#").Int()

	r := make([]Item, size)
	wg := sync.WaitGroup{}
	wg.Add(int(size))

	v.ForEach(
		func(n, s gjson.Result) bool {
			go func(i int64, j gjson.Result) {
				defer wg.Done()

				gid := j.Get(id).String()
				if len(gid) > 2 && gid[:2] == "VL" {
					gid = gid[2:]
				}

				r[i] = Item{
					Id:    gid,
					Title: RunsText(j.Get("title")),
					Sub:   RunsText(j.Get("subtitle")),
					Thumbnails: GetThumbnails(
						j.Get("thumbnailRenderer.musicThumbnailRenderer.thumbnail.thumbnails"),
					),
				}

			}(n.Int(), s)

			return true
		},
	)

	wg.Wait()

	return r
}

func ResponsiveListItemRenderer(s gjson.Result) []Item {
	size := s.Get("#").Int()
	r := make([]Item, size)

	wg := sync.WaitGroup{}
	wg.Add(int(size))

	s.ForEach(
		func(n, s gjson.Result) bool {
			go func(i int64, v gjson.Result) {
				defer wg.Done()

				j := v.Get("musicResponsiveListItemRenderer")
				flex := j.Get(
					"flexColumns.#.musicResponsiveListItemFlexColumnRenderer.text.runs.0",
				)

				r[i] = Item{
					Id:    j.Get("playlistItemData.videoId").String(),
					Title: flex.Get("#(navigationEndpoint.watchEndpoint.videoId).text").String(),
					SubId: flex.Get(
						"#.navigationEndpoint.browseEndpoint.browseId",
					).Get("0").String(),
					Sub: flex.Get(
						"#(navigationEndpoint.browseEndpoint.browseId).text",
					).String(),
					Thumbnails: GetThumbnails(
						j.Get("thumbnail.musicThumbnailRenderer.thumbnail.thumbnails"),
					),
				}

			}(n.Int(), s)

			return true
		},
	)

	wg.Wait()

	return r
}

func ResponsiveListItemRendererCH(s gjson.Result) []Item {
	size := s.Get("#").Int()
	r := make([]Item, size)

	wg := sync.WaitGroup{}
	wg.Add(int(size))

	s.ForEach(
		func(n, s gjson.Result) bool {
			go func(i int64, v gjson.Result) {
				defer wg.Done()

				j := v.Get("musicResponsiveListItemRenderer")
				flex := j.Get(
					"flexColumns.#.musicResponsiveListItemFlexColumnRenderer.text.runs.0.text",
				)

				r[i] = Item{
					Id:    j.Get("navigationEndpoint.browseEndpoint.browseId").String(),
					Title: flex.Get("0").String() + " • " + flex.Get("1").String(),
					Thumbnails: GetThumbnails(
						j.Get("thumbnail.musicThumbnailRenderer.thumbnail.thumbnails"),
					),
				}

			}(n.Int(), s)

			return true
		},
	)

	wg.Wait()

	return r
}

func NavigationButton(s gjson.Result) []Item {
	size := s.Get("#").Int()
	r := make([]Item, size)

	wg := sync.WaitGroup{}
	wg.Add(int(size))

	s.ForEach(
		func(n, s gjson.Result) bool {
			go func(i int64, v gjson.Result) {
				defer wg.Done()

				j := v.Get("musicNavigationButtonRenderer")
				color := j.Get("solid.leftStripeColor").Uint() & 0xffffff

				r[i] = Item{
					Id:    j.Get("clickCommand.browseEndpoint.params").String(),
					Title: RunsText(j.Get("buttonText")),
					Sub:   fmt.Sprintf("#%06x", color),
				}
			}(n.Int(), s)

			return true
		},
	)

	wg.Wait()

	return r
}

func MultiSelectMenuItemRenderer(j, ref gjson.Result) []Item {
	size := j.Get("#").Int()
	r := make([]Item, size)

	wg := sync.WaitGroup{}
	wg.Add(int(size))

	j.ForEach(
		func(n, s gjson.Result) bool {
			go func(i int64, v gjson.Result) {
				defer wg.Done()

				steg := v.Get("formItemEntityKey").String()
				id := ref.Get(
					"#(entityKey == " + steg + ")",
				).Get("payload.musicFormBooleanChoice.opaqueToken")

				r[i] = Item{
					Id:    id.String(),
					Title: RunsText(v.Get("title")),
				}
			}(n.Int(), s)

			return true
		},
	)

	wg.Wait()

	return r
}

func ShelfRenderer(j gjson.Result) map[string]interface{} {

	r := make(map[string]interface{})
	wg := sync.WaitGroup{}

	j.ForEach(
		func(n, s gjson.Result) bool {
			wg.Add(1)

			go func(v gjson.Result) {
				defer wg.Done()

				title := RunsText(
					v.Get("header.musicCarouselShelfBasicHeaderRenderer.title"),
				)

				r[title] = TwoRowItemRenderer(v.Get("contents"), true)
			}(s)

			return true
		},
	)

	wg.Wait()

	return r
}
