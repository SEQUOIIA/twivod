package twitch

import (
	"net/http"
	"io"
	"log"
	"fmt"
	"github.com/sequoiia/twiVod/server/db"
)

type Twitch struct {
	httpCli * http.Client
	clientID string
	apiVersion string
	debug bool
	Db * db.Dbc
}

const OfficialTwitchAPIEndpoint string = "https://api.twitch.tv"

func (t * Twitch) defaultInit() {
	t.httpCli = http.DefaultClient
	t.apiVersion = "application/vnd.twitchtv.v5+json"
	t.debug = false
	t.Db = db.NewDbc("./data.db")
}

func (t * Twitch) debugPrint(method string, value string) {
	if t.debug {
		log.Println(fmt.Sprintf("DEBUG | %s", value))
	}
}

func (t * Twitch) newRequest(httpMethod string, url string, body io.Reader) (*http.Request, error){
	req, err := http.NewRequest(httpMethod, url, body)
	req.Header.Add("Client-ID", t.clientID)
	req.Header.Add("Accept", t.apiVersion)

	return req, err
}

func (t * Twitch) Testerino() {

}

func (t * Twitch) Debug() {
	t.debug = true
}

func NewClient(clientId string) *Twitch{
	t := &Twitch{}
	t.clientID = clientId
	t.defaultInit()

	return t
}