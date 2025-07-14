package utils

import "strings"

type Client struct {
	Name    string `json:"clientName,omitempty"`
	Version string `json:"clientVersion,omitempty"`
	Hl      string `json:"hl",omitempty`
	Url     string `json:"originalUrl",omitempty`
	Visit   string `json:"visitorData",omitempty`
}

type Click struct {
	Param string `json:"clickTrackingParams",omitempty`
}

type Context struct {
	Client Client `json:"client",omitempty`
	Click  Click  `json:"clickTracking",omitempty`
}

type PageType struct {
	PageType string `json:"pageType,omitempty"`
}

type BrowseMusicConfig struct {
	MusicConfig PageType `json:"browseEndpointContextMusicConfig,omitempty"`
}

type Form struct {
	Values []string `json:"selectedValues",omitempty`
}

type WatchMusicConfig struct {
	Panel bool   `json:"hasPersistentPlaylistPanel,omitempty"`
	Type  string `json:"musicVideoType,omitempty"`
}

type BrowseData struct {
	Context     *Context           `json:"context,omitempty"`
	MusicConfig *BrowseMusicConfig `json:"browseEndpointContextMusicConfig,omitempty"`
	BrowseId    string             `json:"browseId,omitempty"`
	Params      *string            `json:"params,omitempty"`
	Form        *Form              `json:"formData,omitempty"`
}

type NextData struct {
	Id          string           `json:"videoId,omitempty"`
	PlaylistId  string           `json:"playlistId",omitempty`
	Context     Context          `json:"context,omitempty"`
	Audio       bool             `json:"isAudioOnly,omitempty"`
	Tuner       string           `json:"tunerSettingValue,omitempty"`
	Panel       bool             `json:"enablePersistentPlaylistPanel,omitempty"`
	MusicConfig WatchMusicConfig `json:"watchEndpointMusicConfig,omitempty"`
}

var BaseContext Context = Context{
	Client: Client{
		Name:    "WEB_REMIX",
		Version: "1.20230320.01.00",
		Hl:      "en",
	},
}

func TypeBrowsePage(id, t string) BrowseData {
	config := BrowseMusicConfig{
		MusicConfig: PageType{
			PageType: "MUSIC_PAGE_TYPE_" + strings.ToUpper(t),
		},
	}

	return BrowseData{
		Context:     &BaseContext,
		MusicConfig: &config,
		BrowseId:    id,
	}
}

func TypeBrowseForm(id, params, form string) BrowseData {
	forms := Form{
		Values: []string{form},
	}

	return BrowseData{
		Context:  &BaseContext,
		BrowseId: id,
		Params:   &params,
		Form:     &forms,
	}
}

func TypeBrowse(id, params string, ct []string) BrowseData {
	data := BrowseData{
		Context:  &BaseContext,
		BrowseId: id,
		Params:   &params,
	}

	if len(ct) > 0 {
		data.Context.Client.Url = "https://music.youtube.com/channel/" + id
		data.Context.Client.Visit = ct[1]
		data.Context.Click = Click{
			Param: ct[0],
		}
	}

	return data
}

func TypeNext(id, pid string) NextData {
	return NextData{
		Id:         id,
		PlaylistId: pid,
		Context:    BaseContext,
		Panel:      true,
		Audio:      true,
		Tuner:      "AUTOMIX_SETTING_NORMAL",
		MusicConfig: WatchMusicConfig{
			Panel: true,
			Type:  "MUSIC_VIDEO_TYPE_ATV",
		},
	}
}
