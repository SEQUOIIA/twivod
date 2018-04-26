package downloader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sequoiia/twivod/internal/github.com/grafov/m3u8"
	"github.com/sequoiia/twivod/models"
	"github.com/sequoiia/twivod/utilities/parser"
	"github.com/sequoiia/twivod/utilities/stream"
)

var HttpClient *http.Client

func Download(vod *models.TwitchVodOptions, bwlimit int64, ds *stream.Client) error {
	vodInfo := parser.VodInfo(vod.Url)

	if vodInfo.Type == models.Unknown {
		return errors.New("unknown URL")
	}

	if ds.Enabled {
		ds.Handle(stream.Container{
			Status:  stream.StatusOK,
			Type:    stream.TypeStage,
			Payload: stream.StageStarted,
		})
	}

	if HttpClient == nil {
		HttpClient = &http.Client{Timeout: 15 * time.Second}
	}

	if vodInfo.Type == models.VOD {
		vodDetails, err := GetVODDetails(vodInfo.ID, HttpClient)
		if err != nil {
			ds.HandleErrorFatal(err)
		}
		vodInfo.Channel = vodDetails.Channel.Name

		if !ds.Enabled {
			fmt.Printf("Downloading VOD '%v' from Twitch channel '%v'\n", vodInfo.ID, vodInfo.Channel)
		} else {
			ds.Handle(stream.Container{
				Status:  stream.StatusOK,
				Type:    stream.TypeStage,
				Payload: stream.StageDownload,
			})
		}
		var token models.HlsVodToken = getAccessToken(HttpClient, vodInfo.ID)

		req, err := http.NewRequest("GET", fmt.Sprintf("https://usher.ttvnw.net/vod/%s.m3u8?nauthsig=%s&allow_source=true&allow_spectre=true&nauth=%s", vodInfo.ID, token.Sig, token.Token), nil)
		if err != nil {
			ds.HandleErrorFatal(err)
		}

		resp, err := HttpClient.Do(req)
		if err != nil {
			ds.HandleErrorFatal(err)
		}

		p, _, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), true)
		if err != nil {
			ds.HandleErrorFatal(err)
		}

		masterPlaylist := p.(*m3u8.MasterPlaylist)

		// Get media playlist
		req, err = http.NewRequest("GET", masterPlaylist.Variants[0].URI, nil)
		if err != nil {
			ds.HandleErrorFatal(err)
		}

		resp, err = HttpClient.Do(req)
		if err != nil {
			ds.HandleErrorFatal(err)
		}

		pMedia, _, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), true)
		if err != nil {
			ds.HandleErrorFatal(err)
		}

		pMediaPlaylist := pMedia.(*m3u8.MediaPlaylist)

		if !ds.Enabled {
			log.Printf("Total seconds: %s, elasped seconds: %s, concurrent option: %d\n", pMediaPlaylist.TwitchInfo.TotalSeconds, pMediaPlaylist.TwitchInfo.ElapsedSeconds, vod.MaxConcurrentDownloads)
		}

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
			//log.Println(vodEndpointCurrent)
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
		//log.Println(vodEndpoint)

		vod.FileName = fmt.Sprintf("%s.ts", vodInfo.ID)
		vod.Name = fmt.Sprintf("%s_%s", vodInfo.Channel, vodInfo.ID)

		//log.Println(endPos)

		file, err := os.Create(fmt.Sprintf("%s%s", vod.SaveFilePath, vod.FileName))
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

		bwLimitCorrected := bwlimit / int64(concurrentAmount)

		for i := startPos; i <= (concurrentAmount + startPos); i++ {
			go downloadSegment(fmt.Sprintf("%s%s", vodEndpoint, pMediaPlaylist.Segments[i].URI), i, channel, 5, bwLimitCorrected)
		}

		buf := make([]*bytes.Buffer, endPos)
		pw := startPos
		pd := startPos + concurrentAmount

		for pw < endPos {
			response := <-channel
			buf[response.Id] = response.Buf
			for pw < endPos && buf[pw] != nil {
				_, err := io.Copy(vod.Writer, buf[pw])
				if err != nil {
					ds.HandleErrorFatal(err)
				}

				if !ds.Enabled {
					log.Printf("Part %d has been downloaded.", pw)
				} else {
					ds.Handle(stream.Container{
						Status: stream.StatusOK,
						Type:   stream.TypeDownloadProgress,
						Payload: stream.DownloadProgress{
							TotalSegments:  endPos,
							CurrentSegment: pw + 1,
						},
					})
				}

				buf[pw] = nil
				pw++
			}
			if pd < endPos {
				go downloadSegment(fmt.Sprintf("%s%s", vodEndpoint, pMediaPlaylist.Segments[pd].URI), pd, channel, 5, bwLimitCorrected)
				pd++
			}
		}

	} else {
		return errors.New("VOD not found on Twitch.TV.")
	}

	vod.Writer.(*os.File).Close()
	return nil
}

func downloadSegment(uri string, vodId int, channel chan models.TwitchVodSegment, retries int, bwLimit int64) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		if retries > 0 {
			downloadSegment(uri, vodId, channel, retries-1, bwLimit)
		}
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		if retries > 0 {
			downloadSegment(uri, vodId, channel, retries-1, bwLimit)
		}
	}

	buf := bytes.NewBuffer(nil)
	for range time.Tick(1 * time.Second) {
		_, err := io.CopyN(buf, resp.Body, bwLimit)
		if err != nil {
			break
		}
	}

	resp.Body.Close()
	channel <- models.TwitchVodSegment{vodId, buf}
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

	defer resp.Body.Close()

	var payload models.VODDetails

	tmpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.VODDetails{}, err
	}

	err = json.Unmarshal(tmpBody, &payload)
	if err != nil {
		return models.VODDetails{}, errors.New(string(tmpBody))
	}

	return payload, nil
}
