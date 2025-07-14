package lib

import (
	"encoding/json"
	"sync"

	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type MediaSession struct {
	Thumbnails []Thumbnail `json:"thumbnails"`
	Album      string      `json:"album"`
}

type Next struct {
	Err          error        `json:"error,omitempty"`
	LyricsId     string       `json:"lyricsId"`
	Songs        []Item       `json:"songs"`
	MediaSession MediaSession `json:"mediaSession"`
}

func parseNextSongs(n gjson.Result) []Item {
	panel := n.Get("#.playlistPanelVideoRenderer")

	size := panel.Get("#").Int()
	r := make([]Item, size)

	wg := sync.WaitGroup{}
	wg.Add(int(size))

	panel.ForEach(
		func(n, v gjson.Result) bool {
			go func(i int64, j gjson.Result) {
				defer wg.Done()

				r[i] = Item{
					Id:         j.Get("navigationEndpoint.watchEndpoint.videoId").String(),
					Title:      RunsText(j.Get("title")),
					Sub:        j.Get("longBylineText.runs.2.text").String(),
					Thumbnails: GetThumbnails(j.Get("thumbnail.thumbnails")),
				}
			}(n.Int(), v)

			return true
		},
	)

	wg.Wait()

	return r
}

func parseNext(raw string) Next {

	j := gjson.Parse(raw)

	contents := j.Get(
		"contents.singleColumnMusicWatchNextResultsRenderer.tabbedRenderer.watchNextTabbedResultsRenderer.tabs",
	)
	mediaSession := j.Get(
		"playerOverlays.playerOverlayRenderer.browserMediaSession.browserMediaSessionRenderer",
	)

	upNext := contents.Get(
		"#(tabRenderer.title == Up next).tabRenderer.content",
	).Get("musicQueueRenderer.content.playlistPanelRenderer")

	lyricsId := contents.Get(
		"#(tabRenderer.title == Lyrics).tabRenderer.endpoint.browseEndpoint.browseId",
	).String()

	return Next{
		LyricsId: lyricsId,
		MediaSession: MediaSession{
			Album:      RunsText(mediaSession.Get("album")),
			Thumbnails: GetThumbnails(mediaSession.Get("thumbnailDetails.thumbnails")),
		},
		Songs: parseNextSongs(upNext.Get("contents")),
	}
}

func GetNext(id string, skip bool) (Next, int) {

	pldata, err := json.Marshal(utils.TypeNext(id, ""))
	if err != nil {
		return Next{Err: err}, 500
	}

	plraw, plstatus, err := utils.Fetch("next", pldata)
	if err != nil || plstatus > 399 {
		return Next{Err: err}, plstatus
	}

	if skip {
		return parseNext(plraw), plstatus
	}

	pl := gjson.Parse(plraw).Get(
		"contents.singleColumnMusicWatchNextResultsRenderer." +
			"tabbedRenderer.watchNextTabbedResultsRenderer.tabs.0.tabRenderer.content." +
			"musicQueueRenderer.content.playlistPanelRenderer.contents." +
			"#(automixPreviewVideoRenderer).automixPreviewVideoRenderer." +
			"content.automixPlaylistVideoRenderer.navigationEndpoint." +
			"watchPlaylistEndpoint.playlistId").String()

	data, err := json.Marshal(utils.TypeNext(id, pl))
	if err != nil {
		return Next{Err: err}, 500
	}

	raw, status, err := utils.Fetch("next", data)
	if err != nil {
		return Next{Err: err}, 500
	}

	return parseNext(raw), status
}
