package downloader

import (
	"github.com/sequoiia/twiVod/models"
	"github.com/sequoiia/twiVod/utilities/parser"
	"errors"
	"fmt"
	"net/http"
	"log"
	"github.com/grafov/m3u8"
	"bufio"
	"io"
	"os"
	"bytes"
)

var HttpClient *http.Client

func Download(vod *models.TwitchVodOptions) (error){
	vodInfo := parser.VodInfo(vod.Url)

	if HttpClient == nil {
		HttpClient = http.DefaultClient
	}

	if vodInfo.Type != "404" {
		fmt.Printf("Downloading VOD '%v' from Twitch channel '%v'\n",  vodInfo.ID, vodInfo.Channel)
		var token models.HlsVodToken = getAccessToken(HttpClient, vodInfo.ID)
//		var vodKraken models.VodInfoKraken = getVodInfo(HttpClient, vodInfo)

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

		log.Println(masterPlaylist.Variants[0].URI)

		resp, err = HttpClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		pMedia, _, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), true)
		if err != nil {
			log.Fatal(err)
		}

		pMediaPlaylist := pMedia.(*m3u8.MediaPlaylist)
		log.Printf("Total seconds: %s, elasped seconds: %s\n", pMediaPlaylist.TwitchInfo.TotalSeconds, pMediaPlaylist.TwitchInfo.ElapsedSeconds)

		var w io.WriteCloser
		var endPos int = int(pMediaPlaylist.Count())
		var startPos int = 0
		var concurrentAmount int = 0
		var bytesBuffer bytes.Buffer
		var vodEndpoint string = ""
		var vodEndpointCurrent = 0

		for vodEndpointCurrent = 0; vodEndpointCurrent <= (len(masterPlaylist.Variants[0].URI) - 15); vodEndpointCurrent++ {
			bytesBuffer.WriteByte(masterPlaylist.Variants[0].URI[vodEndpointCurrent])
		}
		vodEndpoint = bytesBuffer.String()

		vod.FileName = fmt.Sprintf("%s.ts", vodInfo.ID)
		vod.Name = fmt.Sprintf("%s_%s", vodInfo.Channel,vodInfo.ID)

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
			go downloadSegment(fmt.Sprintf("%s%s", vodEndpoint, pMediaPlaylist.Segments[i].URI), i, channel)
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
				go downloadSegment(fmt.Sprintf("%s%s", vodEndpoint, pMediaPlaylist.Segments[pd].URI), pd, channel)
				pd++
			}
		}

		log.Println("This text should be here")

	} else {
		return errors.New("VOD not found on Twitch.TV.")
	}

	return nil
}

func downloadSegment(uri string, vodId int, channel chan models.TwitchVodSegment) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	channel <- models.TwitchVodSegment{vodId, resp.Body}
}