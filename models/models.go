package models

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