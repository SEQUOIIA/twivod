package models

import (
	"bytes"
	"time"
)

// StreamType is a string type to help differentiate between live streams and VODs.
type StreamType string

// Livestream is a StreamType
const Livestream StreamType = "livestream"

// VOD is a StreamType
const VOD StreamType = "videos"

// Unknown is a StreamType
const Unknown StreamType = "unknown"

type StreamTypeResult struct {
	Type  StreamType
	Value string
}

type VODDetails  struct {
	Title              string      `json:"title"`
	Description        interface{} `json:"description"`
	DescriptionHTML    interface{} `json:"description_html"`
	BroadcastID        int64       `json:"broadcast_id"`
	BroadcastType      string      `json:"broadcast_type"`
	Status             string      `json:"status"`
	Language           string      `json:"language"`
	TagList            string      `json:"tag_list"`
	Views              int         `json:"views"`
	CreatedAt          time.Time   `json:"created_at"`
	PublishedAt        time.Time   `json:"published_at"`
	URL                string      `json:"url"`
	ID                 string      `json:"_id"`
	RecordedAt         time.Time   `json:"recorded_at"`
	Game               string      `json:"game"`
	Length             int         `json:"length"`
	Preview            string      `json:"preview"`
	AnimatedPreviewURL string      `json:"animated_preview_url"`
	Thumbnails         []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"thumbnails"`
	Fps struct {
		One60P30   float64 `json:"160p30"`
		Three60P30 float64 `json:"360p30"`
		Four80P30  float64 `json:"480p30"`
		Seven20P30 float64 `json:"720p30"`
		Seven20P60 float64 `json:"720p60"`
		Chunked    float64 `json:"chunked"`
	} `json:"fps"`
	Resolutions struct {
		One60P30   string `json:"160p30"`
		Three60P30 string `json:"360p30"`
		Four80P30  string `json:"480p30"`
		Seven20P30 string `json:"720p30"`
		Seven20P60 string `json:"720p60"`
		Chunked    string `json:"chunked"`
	} `json:"resolutions"`
	Channel struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
	} `json:"channel"`
	Links struct {
		Self    string `json:"self"`
		Channel string `json:"channel"`
	} `json:"_links"`
}

type VODinfo struct {
	Channel string
	Type    StreamType
	ID      string
	Url     string
}

type VODtypeB struct {
	Chunks struct {
		Live []struct {
			Url string `json:"url"`
		} `json:"live"`
		Chunked []struct {
			Url string `json:"url"`
		} `json:"chunked"`
	} `json:"chunks"`
}

type TwitchOauthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type HlsVodToken struct {
	Token string `json:"token"`
	Sig   string `json:"sig"`
}

type VodInfoKraken struct {
	Title         string      `json:"title"`
	Description   interface{} `json:"description"`
	BroadcastID   int64       `json:"broadcast_id"`
	BroadcastType string      `json:"broadcast_type"`
	Status        string      `json:"status"`
	TagList       string      `json:"tag_list"`
	Views         int         `json:"views"`
	CreatedAt     time.Time   `json:"created_at"`
	ID            string      `json:"_id"`
	RecordedAt    time.Time   `json:"recorded_at"`
	Game          string      `json:"game"`
	Length        float64     `json:"length"`
	Preview       string      `json:"preview"`
	Thumbnails    []struct {
		URL  string `json:"url"`
		Type string `json:"type"`
	} `json:"thumbnails"`
	URL string `json:"url"`
	Fps struct {
		AudioOnly float64 `json:"audio_only"`
		Medium    float64 `json:"medium"`
		Mobile    float64 `json:"mobile"`
		High      float64 `json:"high"`
		Low       float64 `json:"low"`
		Chunked   float64 `json:"chunked"`
	} `json:"fps"`
	Resolutions struct {
		Medium  string `json:"medium"`
		Mobile  string `json:"mobile"`
		High    string `json:"high"`
		Low     string `json:"low"`
		Chunked string `json:"chunked"`
	} `json:"resolutions"`
	Links struct {
		Self    string `json:"self"`
		Channel string `json:"channel"`
	} `json:"_links"`
	Channel struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
	} `json:"channel"`
}

type TwitchUser struct {
	DisplayName string    `json:"display_name"`
	ID          int       `json:"_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Logo        string    `json:"logo"`
	Links       struct {
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
		Title         string      `json:"title"`
		Description   interface{} `json:"description"`
		BroadcastID   int64       `json:"broadcast_id"`
		BroadcastType string      `json:"broadcast_type"`
		Status        string      `json:"status"`
		TagList       string      `json:"tag_list"`
		Views         int         `json:"views"`
		CreatedAt     time.Time   `json:"created_at"`
		ID            string      `json:"_id"`
		RecordedAt    time.Time   `json:"recorded_at"`
		Game          string      `json:"game"`
		Length        int         `json:"length"`
		Preview       string      `json:"preview"`
		Thumbnails    []struct {
			URL  string `json:"url"`
			Type string `json:"type"`
		} `json:"thumbnails"`
		URL string `json:"url"`
		Fps struct {
			AudioOnly float64 `json:"audio_only"`
			Medium    float64 `json:"medium"`
			Mobile    float64 `json:"mobile"`
			High      float64 `json:"high"`
			Low       float64 `json:"low"`
			Chunked   float64 `json:"chunked"`
		} `json:"fps"`
		Resolutions struct {
			Medium  string `json:"medium"`
			Mobile  string `json:"mobile"`
			High    string `json:"high"`
			Low     string `json:"low"`
			Chunked string `json:"chunked"`
		} `json:"resolutions"`
		Links struct {
			Self    string `json:"self"`
			Channel string `json:"channel"`
		} `json:"_links"`
		Channel struct {
			Name        string `json:"name"`
			DisplayName string `json:"display_name"`
		} `json:"channel"`
	} `json:"videos"`
}

type TwitchUsersSearch struct {
	Total int `json:"_total"`
	Links struct {
		Self string `json:"self"`
		Next string `json:"next"`
	} `json:"_links"`
	Channels []struct {
		Mature                       bool        `json:"mature"`
		Status                       string      `json:"status"`
		BroadcasterLanguage          string      `json:"broadcaster_language"`
		DisplayName                  string      `json:"display_name"`
		Game                         string      `json:"game"`
		Language                     string      `json:"language"`
		ID                           int         `json:"_id"`
		Name                         string      `json:"name"`
		CreatedAt                    time.Time   `json:"created_at"`
		UpdatedAt                    time.Time   `json:"updated_at"`
		Delay                        interface{} `json:"delay"`
		Logo                         string      `json:"logo"`
		Banner                       interface{} `json:"banner"`
		VideoBanner                  string      `json:"video_banner"`
		Background                   interface{} `json:"background"`
		ProfileBanner                string      `json:"profile_banner"`
		ProfileBannerBackgroundColor interface{} `json:"profile_banner_background_color"`
		Partner                      bool        `json:"partner"`
		URL                          string      `json:"url"`
		Views                        int         `json:"views"`
		Followers                    int         `json:"followers"`
		Links                        struct {
			Self          string `json:"self"`
			Follows       string `json:"follows"`
			Commercial    string `json:"commercial"`
			StreamKey     string `json:"stream_key"`
			Chat          string `json:"chat"`
			Features      string `json:"features"`
			Subscriptions string `json:"subscriptions"`
			Editors       string `json:"editors"`
			Teams         string `json:"teams"`
			Videos        string `json:"videos"`
		} `json:"_links"`
	} `json:"channels"`
}

type TwitchVodSegment struct {
	Id           int
	Buf *bytes.Buffer
}
