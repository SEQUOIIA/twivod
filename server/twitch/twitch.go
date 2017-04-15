package twitch

import (
	"net/http"
	"io"
)

type Twitch struct {
	httpCli * http.Client
	clientID string
	apiVersion string
}

const OfficialTwitchAPIEndpoint string = "https://api.twitch.tv"

func (t * Twitch) defaultInit() {
	t.httpCli = http.DefaultClient
	t.apiVersion = "application/vnd.twitchtv.v5+json"
}

func (t * Twitch) newRequest(httpMethod string, url string, body io.Reader) (*http.Request, error){
	req, err := http.NewRequest(httpMethod, url, body)
	req.Header.Add("Client-ID", t.clientID)
	req.Header.Add("Accept", t.apiVersion)

	return req, err
}

func (t * Twitch) Testerino() {

}

func NewClient(clientId string) *Twitch{
	t := &Twitch{}
	t.defaultInit()

	return t
}