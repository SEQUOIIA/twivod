package downloader

import (
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/sequoiia/twiVod/models"
	"github.com/sequoiia/twiVod/utilities/parser"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sync"
    	"github.com/grafov/m3u8"
    	"bufio"
	"log"
)

func legacydl(url string, filename string, wg *sync.WaitGroup) {
	out, err := os.Create(filename + ".flv")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bar := pb.New(int(resp.ContentLength)).SetUnits(pb.U_BYTES)

	bar.ShowSpeed = true
	bar.Start()

	writer := io.MultiWriter(out, bar)
	io.Copy(writer, resp.Body)
	bar.Finish()

	wg.Done()
}

func Get(urlarg string) {
	vod := parser.VodInfo(fmt.Sprintf(urlarg))
	if vod.Type != "404" {
		fmt.Printf("\nDownloading VOD '%v' from Twitch channel '%v'\n",  vod.ID, vod.Channel)
		cli := http.DefaultClient
		var token models.HlsVodToken = getAccessToken(cli, vod.ID)

		req, err := http.NewRequest("GET", fmt.Sprintf("https://usher.ttvnw.net/vod/%s.m3u8?nauthsig=%s&allow_source=true&allow_spectre=true&nauth=%s", vod.ID, token.Sig, token.Token), nil)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := cli.Do(req)

		p, _, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), true)
		if err != nil {
			log.Fatal(err)
		}

		masterPlayList := p.(*m3u8.MasterPlaylist)
		fmt.Println(masterPlayList.TwitchInfos[0])
	}
}

func getAccessToken(cli *http.Client, vodId string) models.HlsVodToken {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/api/vods/%s/access_token", vodId), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	tmpbody, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	var returnModel models.HlsVodToken

	err = json.Unmarshal(tmpbody, &returnModel)
	if err != nil {
		log.Fatal(err)
	}

	return returnModel
}

func LegacyGet(urlarg string) {
	vod := parser.VodInfo(fmt.Sprintf(urlarg))
	if vod.Type != "404" {

		fmt.Printf("\nDownloading VOD '%v' from Twitch channel '%v'\n", vod.ID, vod.Channel)

		if vod.Type == "b" {
			cli := http.DefaultClient
			endpoint := "https://api.twitch.tv/api/videos/a" + vod.ID
			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				panic(err)
			}
			req.Header.Set("User-Agent", "twiVod - https://github.com/equoia/twivod")

			resp, err := cli.Do(req)
			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()
			var apiresponse models.VODtypeB
			tmpbody, err := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(tmpbody, &apiresponse)
			if err != nil {
				panic(err)
			}

			if len(apiresponse.Chunks.Live) == 0 {
				fmt.Println("Looks like this VOD is subscriber only, you will need to authenticate before proceeding.\n Go to http://localhost:7261")
				tmpdatafile, err := os.Create("vod_oauth")
				if err != nil {
				    panic(err)
				}

				jsonedvod, err := json.Marshal(vod)
				if err != nil {
				    panic(err)
				}

				tmpdatafile.Write([]byte(jsonedvod))
						// Add oauth authentication here
				Oauth(vod)
			} else {
				//var vodurls []string

				var wg sync.WaitGroup
				wg.Add(len(apiresponse.Chunks.Live))
				for _, data := range apiresponse.Chunks.Live {
					r := regexp.MustCompile(`.*.tv\/.*?(live.*)\.`)
					go legacydl(data.Url, r.FindStringSubmatch(data.Url)[1], &wg)
					//vodurls = append(vodurls, data.Url)
				}
				//fmt.Println(len(vodurls))
				wg.Wait()

				for _, data := range apiresponse.Chunks.Live {
					r := regexp.MustCompile(`.*.tv\/.*?(live.*)\.`)
					filenameflv := r.FindStringSubmatch(data.Url)[1] + ".flv"
                    filenamemp4 := r.FindStringSubmatch(data.Url)[1] + ".mp4"
                    cmd := exec.Command("ffmpeg", "-i", filenameflv, "-vcodec", "copy", "-acodec", "copy", filenamemp4)
                    cmd.Stdout = os.Stdout
                    cmd.Stdin = os.Stdin
                    cmd.Stderr = os.Stderr
                    cmd.Run()

					//vodurls = append(vodurls, data.Url)
				}

				remuxfile, err := os.Create("demux.txt")
				if err != nil {
					panic(err)
				}

				for _, data := range apiresponse.Chunks.Live {
					r := regexp.MustCompile(`.*.tv\/.*?(live.*)\.`)
					fullstring := "file '" + r.FindStringSubmatch(data.Url)[1] + ".mp4'\n"
					_, err := remuxfile.WriteString(fullstring)
					if err != nil {
						panic(err)
					}
				}

                filenamemp4 := vod.Channel + "_" + vod.ID + ".mp4"
				cmd := exec.Command("ffmpeg", "-f", "concat", "-i", "demux.txt", "-c", "copy", filenamemp4)
                cmd.Stdout = os.Stdout
                cmd.Stdin = os.Stdin
                cmd.Stderr = os.Stderr
                cmd.Run()

				for _, data := range apiresponse.Chunks.Live {
					r := regexp.MustCompile(`.*.tv\/.*?(live.*)\.`)
					os.Remove(r.FindStringSubmatch(data.Url)[1] + ".mp4")
					os.Remove(r.FindStringSubmatch(data.Url)[1] + ".flv")
				}
				os.Remove("demux.txt")

				fmt.Println("Done!")
			}

		} else if vod.Type == "v" {

            cli := http.DefaultClient

			endpoint := "https://api.twitch.tv/api/videos/v" + vod.ID
			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				panic(err)
			}
			req.Header.Set("User-Agent", "twiVod - https://github.com/equoia/twivod")

			resp, err := cli.Do(req)
			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()
			var apiresponse models.VODtypeB
			tmpbody, err := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(tmpbody, &apiresponse)
			if err != nil {
				panic(err)
			}

			if len(apiresponse.Chunks.Chunked) == 0 {
				fmt.Println("Looks like this VOD is subscriber only, you will need to authenticate before proceeding.\n Go to http://localhost:7261")
				tmpdatafile, err := os.Create("vod_oauth")
				if err != nil {
					panic(err)
				}

				jsonedvod, err := json.Marshal(vod)
				if err != nil {
					panic(err)
				}

				tmpdatafile.Write([]byte(jsonedvod))
				// Add oauth authentication here
				Oauth(vod)
			} else {




            endpoint = "https://api.twitch.tv/api/vods/" + vod.ID + "/access_token"
            req, err = http.NewRequest("GET", endpoint, nil)
            if err != nil {
                panic(err)
            }
            req.Header.Set("User-Agent", "twiVod - https://github.com/equoia/twivod")

            rsp, err := cli.Do(req)
            if err != nil {
                panic(err)
            }

            defer rsp.Body.Close()
            var vodToken models.HlsVodToken

            tmpbody, err = ioutil.ReadAll(rsp.Body)
            if err != nil {
                panic(err)
            }

            err = json.Unmarshal(tmpbody, &vodToken)
            if err != nil {
                panic(err)
            }

            endpoint = "http://usher.justin.tv/vod/" + vod.ID + "?nauthsig=" + vodToken.Sig + "&nauth=" + vodToken.Token + "&allow_source=true"
            req, err = http.NewRequest("GET", endpoint, nil)
            if err != nil {
                panic(err)
            }
            req.Header.Set("User-Agent", "twiVod - https://github.com/equoia/twivod")

            rsp, err = cli.Do(req)
            if err != nil {
                panic(err)
            }

            defer rsp.Body.Close()

            p, listType, err := m3u8.DecodeFrom(bufio.NewReader(rsp.Body), true)
            if err != nil {
                panic(err)
            }

            switch listType {
                case m3u8.MEDIA:
                mediapl := p.(*m3u8.MediaPlaylist)
                fmt.Printf("%+v\n", mediapl)
                case m3u8.MASTER:
                masterpl := p.(*m3u8.MasterPlaylist)
                //fmt.Printf("%+v\n", masterpl.Variants[5])
                for _, data := range masterpl.Variants {
                    if data.Video == "chunked" {
                        fmt.Println(data.URI)
                        ffmpegargs := vod.Channel + "_" + vod.ID + ".mp4"
                        cmd := exec.Command("ffmpeg", "-analyzeduration", "1000000000", "-probesize", "1000000000", "-i" , data.URI, "-bsf:a", "aac_adtstoasc", "-c", "copy", ffmpegargs)
                        cmd.Stdout = os.Stdout
                        cmd.Stdin = os.Stdin
                        cmd.Stderr = os.Stderr
                        cmd.Run()
                        fmt.Println("Done!")
                    }
                }
            }
        }
		}
	}

}
