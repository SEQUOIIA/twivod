package parser

import (
	"fmt"
	"github.com/equoia/twiVod/models"
	"regexp"
)

func VodInfo(url string) models.VODinfo {
	r := regexp.MustCompile(`.*twitch.tv\/(.*?)\/(.)\/(.*[0-9])`)
	if r.MatchString(url) == true {
		vod := models.VODinfo{Channel: r.FindStringSubmatch(url)[1], Type: r.FindStringSubmatch(url)[2], ID: r.FindStringSubmatch(url)[3], Url: r.FindStringSubmatch(url)[0]}
		return vod
	} else {
		fmt.Println("Sorry, but the URL you submitted doesn't look like a Twitch.TV VOD URL. Try again.")
		vod := models.VODinfo{Type: "404"}
		return vod
	}
}
