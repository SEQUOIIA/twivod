package downloader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/grafov/m3u8"
	"github.com/sequoiia/twiVod/models"
	"github.com/sequoiia/twiVod/utilities/parser"
)

var HttpClient *http.Client

func Download(vod *models.TwitchVodOptions) error {
	vodInfo := parser.VodInfo(vod.Url)

	if HttpClient == nil {
		HttpClient = http.DefaultClient
	}

	vodDetails, err := GetVODDetails(vodInfo.ID, HttpClient)
	if err != nil {
		log.Fatal(err)
	}
	vodInfo.Channel = vodDetails.Channel.Name

	if vodInfo.Type != "404" {
		fmt.Printf("Downloading VOD '%v' from Twitch channel '%v'\n", vodInfo.ID, vodInfo.Channel)
		var token models.HlsVodToken = getAccessToken(HttpClient, vodInfo.ID)

		req, err := http.NewRequest("GET", fmt.Sprintf("https://usher.ttvnw.net/vod/%s.m3u8?nauthsig=%s&allow_source=true&allow_spectre=true&nauth=%s", vodInfo.ID, token.Sig, token.Token), nil)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := HttpClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		p, _, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), true)
		if err != nil {
			log.Fatal(err)
		}

		masterPlaylist := p.(*m3u8.MasterPlaylist)

		// Get media playlist
		req, err = http.NewRequest("GET", masterPlaylist.Variants[0].URI, nil)
		if err != nil {
			log.Fatal(err)
		}

		resp, err = HttpClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		pMedia, _, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), true)
		if err != nil {
			log.Fatal(err)
		}

		pMediaPlaylist := pMedia.(*m3u8.MediaPlaylist)
		log.Printf("Total seconds: %s, elasped seconds: %s, concurrent option: %d\n", pMediaPlaylist.TwitchInfo.TotalSeconds, pMediaPlaylist.TwitchInfo.ElapsedSeconds, vod.MaxConcurrentDownloads)

		var w io.WriteCloser
		var endPos int = int(pMediaPlaylist.Count())
		var startPos int = 0
		var concurrentAmount int = 0
		var bytesBuffer bytes.Buffer
		var vodEndpoint string = ""
		var vodEndpointCurrent = len(masterPlaylist.Variants[0].URI) - 1
		var vodEndpointSlashReached bool = false
		var vodEndpointSlashP = 0

		for !vodEndpointSlashReached {
			log.Println(vodEndpointCurrent)
			if (masterPlaylist.Variants[0].URI[vodEndpointCurrent]) == '/' {
				vodEndpointSlashReached = true
				vodEndpointSlashP = vodEndpointCurrent
				break
			}

			vodEndpointCurrent--
		}

		vodEndpointCurrent = 0

		for i := 0; i <= vodEndpointSlashP; i++ {
			bytesBuffer.WriteByte(masterPlaylist.Variants[0].URI[vodEndpointCurrent])
			vodEndpointCurrent++
		}

		vodEndpoint = bytesBuffer.String()
		log.Println(vodEndpoint)

		vod.FileName = fmt.Sprintf("%s.ts", vodInfo.ID)
		vod.Name = fmt.Sprintf("%s_%s", vodInfo.Channel, vodInfo.ID)

		log.Println(endPos)

		file, err := os.Create(vod.FileName)
		if err != nil {
			log.Fatal(err)
		}

		w = file
		vod.Writer = w

		if vod.MaxConcurrentDownloads > endPos {
			concurrentAmount = endPos
		} else {
			concurrentAmount = vod.MaxConcurrentDownloads
		}

		channel := make(chan models.TwitchVodSegment)

		for i := startPos; i <= (concurrentAmount + startPos); i++ {
			go downloadSegment(fmt.Sprintf("%s%s", vodEndpoint, pMediaPlaylist.Segments[i].URI), i, channel, 5)
		}

		buf := make([]io.ReadCloser, endPos)
		pw := startPos
		pd := startPos + concurrentAmount

		for pw < endPos {
			response := <-channel
			buf[response.Id] = response.ResponseBody
			for pw < endPos && buf[pw] != nil {
				_, err := io.Copy(vod.Writer, buf[pw])
				if err != nil {
					log.Fatal(err)
				}
				buf[pw].Close()
				log.Printf("Part %d has been downloaded.", pw)
				pw++
			}
			if pd < endPos {
				go downloadSegment(fmt.Sprintf("%s%s", vodEndpoint, pMediaPlaylist.Segments[pd].URI), pd, channel, 5)
				pd++
			}
		}

	} else {
		return errors.New("VOD not found on Twitch.TV.")
	}

	return nil
}

func downloadSegment(uri string, vodId int, channel chan models.TwitchVodSegment, retries int) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		if retries > 0 {
			downloadSegment(uri, vodId, channel, retries-1)
		}
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		if retries > 0 {
			downloadSegment(uri, vodId, channel, retries-1)
		}
	}

	channel <- models.TwitchVodSegment{vodId, resp.Body}
}

func GetVODDetails(id string, cli *http.Client) (models.VODDetails, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/videos/v%s", id), nil)
	if err != nil {
		return models.VODDetails{}, err
	}

	req.Header.Set("client-id", models.TwitchConfig.Client_id)

	resp, err := cli.Do(req)
	if err != nil {
		return models.VODDetails{}, err
	}

	var payload models.VODDetails

	err = json.NewDecoder(resp.Body).Decode(&payload); if err != nil {
		return models.VODDetails{}, err
	}

	return payload, nil
}
