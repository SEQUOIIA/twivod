package models

import "time"

type VODinfo struct {
	Channel		string
	Type		string
	ID			string
	Url			string
}

type VODtypeB struct {
	Chunks	struct {
		Live	[]struct {
			Url	string		`json:"url"`
		}	`json:"live"`
		Chunked	[]struct {
			Url	string		`json:"url"`
		}	`json:"chunked"`
	}	`json:"chunks"`
}

type TwitchOauthResponse struct {
    AccessToken     string      `json:"access_token"`
    RefreshToken    string      `json:"refresh_token"`
}

type HlsVodToken struct {
    Token   string  `json:"token"`
    Sig     string  `json:"sig"`
}

type VodInfoKraken struct {
	Title string `json:"title"`
	Description interface{} `json:"description"`
	BroadcastID int64 `json:"broadcast_id"`
	BroadcastType string `json:"broadcast_type"`
	Status string `json:"status"`
	TagList string `json:"tag_list"`
	Views int `json:"views"`
	CreatedAt time.Time `json:"created_at"`
	ID string `json:"_id"`
	RecordedAt time.Time `json:"recorded_at"`
	Game string `json:"game"`
	Length float64 `json:"length"`
	Preview string `json:"preview"`
	Thumbnails []struct {
		URL string `json:"url"`
		Type string `json:"type"`
	} `json:"thumbnails"`
	URL string `json:"url"`
	Fps struct {
		      AudioOnly float64 `json:"audio_only"`
		      Medium float64 `json:"medium"`
		      Mobile float64 `json:"mobile"`
		      High float64 `json:"high"`
		      Low float64 `json:"low"`
		      Chunked float64 `json:"chunked"`
	      } `json:"fps"`
	Resolutions struct {
		      Medium string `json:"medium"`
		      Mobile string `json:"mobile"`
		      High string `json:"high"`
		      Low string `json:"low"`
		      Chunked string `json:"chunked"`
	      } `json:"resolutions"`
	Links struct {
		      Self string `json:"self"`
		      Channel string `json:"channel"`
	      } `json:"_links"`
	Channel struct {
		      Name string `json:"name"`
		      DisplayName string `json:"display_name"`
	      } `json:"channel"`
}

type TwitchUser struct {
	DisplayName string `json:"display_name"`
	ID int `json:"_id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Bio string `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Logo string `json:"logo"`
	Links struct {
			    Self string `json:"self"`
		    } `json:"_links"`
}

type TwitchUserVods struct {
	Total int `json:"_total"`
	Links struct {
		      Self string `json:"self"`
		      Next string `json:"next"`
	      } `json:"_links"`
	Videos []struct {
		Title string `json:"title"`
		Description interface{} `json:"description"`
		BroadcastID int64 `json:"broadcast_id"`
		BroadcastType string `json:"broadcast_type"`
		Status string `json:"status"`
		TagList string `json:"tag_list"`
		Views int `json:"views"`
		CreatedAt time.Time `json:"created_at"`
		ID string `json:"_id"`
		RecordedAt time.Time `json:"recorded_at"`
		Game string `json:"game"`
		Length int `json:"length"`
		Preview string `json:"preview"`
		Thumbnails []struct {
			URL string `json:"url"`
			Type string `json:"type"`
		} `json:"thumbnails"`
		URL string `json:"url"`
		Fps struct {
			      AudioOnly float64 `json:"audio_only"`
			      Medium float64 `json:"medium"`
			      Mobile float64 `json:"mobile"`
			      High float64 `json:"high"`
			      Low float64 `json:"low"`
			      Chunked float64 `json:"chunked"`
		      } `json:"fps"`
		Resolutions struct {
			      Medium string `json:"medium"`
			      Mobile string `json:"mobile"`
			      High string `json:"high"`
			      Low string `json:"low"`
			      Chunked string `json:"chunked"`
		      } `json:"resolutions"`
		Links struct {
			      Self string `json:"self"`
			      Channel string `json:"channel"`
		      } `json:"_links"`
		Channel struct {
			      Name string `json:"name"`
			      DisplayName string `json:"display_name"`
		      } `json:"channel"`
	} `json:"videos"`
}