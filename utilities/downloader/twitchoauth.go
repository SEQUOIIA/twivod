package downloader

import (
    "net/http"
    "net/url"
    "html/template"
    "io/ioutil"
    "bytes"
    "sync"
    "github.com/sequoiia/twivod/models"
    "encoding/json"
    "fmt"
    "regexp"
    "github.com/skratchdot/open-golang/open"
    "os/exec"
    "os"
    "log"
    "github.com/GeertJohan/go.rice"
    "github.com/sequoiia/twivod/internal/github.com/grafov/m3u8"
    "bufio"
    "github.com/gorilla/mux"
    "github.com/codegangsta/negroni")

const CtxSpanID = 0

func root(w http.ResponseWriter, r *http.Request) {
    p := ""


    templateBox, err := rice.FindBox("views/src")
    if err != nil {
        panic(err)
    }

    templateString, err := templateBox.String("index.html")
    if err != nil {
        panic(err)
    }

    htmlTemplate, err := template.New("index").Parse(templateString)
    if err != nil {
        panic(err)
    }

    htmlTemplate.Execute(w, p)
    /*
    t, _ := template.ParseFiles("views/index.html")
    t.Execute(w, p)
    */
}

func dlvod(twitchoauthresponse []byte) {
    vodjson, err := ioutil.ReadFile("vod_oauth")
    if err != nil {
        panic(err)
    }
    var vod models.VODinfo
    var twitchtokens models.TwitchOauthResponse
    err = json.Unmarshal(vodjson, &vod)
    if err != nil {
        panic(err)
    }


    err = json.Unmarshal(twitchoauthresponse, &twitchtokens)
    if err != nil {
        panic(err)
    }

    fmt.Println("Access token acquired! Trying to download VOD again.")


    cli := http.DefaultClient

    if vod.Type == "v" {
        req, err := http.NewRequest("GET", "https://api.twitch.tv/api/vods/" + vod.ID + "/access_token?oauth_token=" + twitchtokens.AccessToken, nil)
        if err != nil {
            panic(err)
        }
        req.Header.Set("User-Agent", "twiVod - https://github.com/equoia/twivod")

        resp, err := cli.Do(req)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()
        var vodTokenresponse struct{
            Token   string  `json:"token"`
            Sig     string  `json:"sig"`
        }

        tmpbody, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            panic(err)
        }

        err = json.Unmarshal(tmpbody, &vodTokenresponse)
        if err != nil {
            panic(err)
        }

        endpoint := "http://usher.justin.tv/vod/" + vod.ID + "?nauthsig=" + vodTokenresponse.Sig + "&nauth=" + vodTokenresponse.Token + "&allow_source=true"
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
                }
            }
        }
        os.Remove("vod_oauth")
        log.Fatal("Done!")
    } else {


    endpoint := "https://api.twitch.tv/api/videos/a" + vod.ID + "?oauth_token=" + twitchtokens.AccessToken
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
        os.Remove("vod_oauth")
        log.Fatal("Doesn't look like you have access to the VOD, check if you are subscribed to the channel.")
    } else {

    var wg sync.WaitGroup
    wg.Add(len(apiresponse.Chunks.Live))
    for _, data := range apiresponse.Chunks.Live {
        r := regexp.MustCompile(`.*.tv\/.*?(live.*)\.`)
        go legacydl(data.Url, r.FindStringSubmatch(data.Url)[1], &wg)
        //vodurls = append(vodurls, data.Url)
    }

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

    ffmpegargs := vod.Channel + "_" + vod.ID + ".mp4"
    cmd := exec.Command("ffmpeg", "-f", "concat", "-i", "demux.txt", "-c", "copy", ffmpegargs)
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    cmd.Run()

    }
    for _, data := range apiresponse.Chunks.Live {
        r := regexp.MustCompile(`.*.tv\/.*?(live.*)\.`)
        os.Remove(r.FindStringSubmatch(data.Url)[1] + ".mp4")
        os.Remove(r.FindStringSubmatch(data.Url)[1] + ".flv")
    }
    os.Remove("demux.txt")
    os.Remove("vod_oauth")

    log.Fatal("Done!")
    }
}

func success(w http.ResponseWriter, r *http.Request) {
    ctx := make(map[string]interface{})

    ctx["code"] = r.URL.Query().Get("code")

    cli := http.DefaultClient
    payload := url.Values{}
    payload.Add("client_id", models.TwitchConfig.Client_id)
    payload.Add("client_secret", models.TwitchConfig.Client_secret)
    payload.Add("grant_type", "authorization_code")
    payload.Add("redirect_uri", "http://localhost:7261/success")
    payload.Add("code", r.URL.Query().Get("code"))

    req, err := http.NewRequest("POST", "https://api.twitch.tv/kraken/oauth2/token", bytes.NewBufferString(payload.Encode()))
    if err != nil {
        panic(err)
    }

    rsp, err := cli.Do(req)
    if err != nil {
        panic(err)
    }
    defer rsp.Body.Close()

    tmpbody, err := ioutil.ReadAll(rsp.Body)
    if err != nil {
        panic(err)
    }

    ctx["payload"] = string(tmpbody)

    go dlvod(tmpbody)

    templateBox, err := rice.FindBox("views/src")
    if err != nil {
        panic(err)
    }

    templateString, err := templateBox.String("success.html")
    if err != nil {
        panic(err)
    }

    htmlTemplate, err := template.New("success").Parse(templateString)
    if err != nil {
        panic(err)
    }

    htmlTemplate.Execute(w, ctx)

    /*
    t, _ := template.ParseFiles("views/success.html")
    t.Execute(w, ctx)
    */
}

func Oauth(vod models.VODinfo) {
    router := mux.NewRouter()
    router.HandleFunc("/", root)
    router.HandleFunc("/success", success)
	cssfileServer := http.StripPrefix("/css/", http.FileServer(rice.MustFindBox("views/css").HTTPBox()))
	router.PathPrefix("/css/").Handler(cssfileServer)

    n := negroni.Classic()
    n.UseHandler(router)

    open.Start("http://localhost:7261")
    n.Run("0.0.0.0:7261")

    /*http.HandleFunc("/", root)
    http.HandleFunc("/success", success)
    cssfileServer := http.StripPrefix("/css/", http.FileServer(rice.MustFindBox("views/css").HTTPBox()))
    http.Handle("/css/", cssfileServer)
    http.ListenAndServe(":7261", nil)
    */
}
