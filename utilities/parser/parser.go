package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sequoiia/twiVod/models"
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

// DetermineStreamType runs a regex expression that determines whether a live stream or a VOD has been given
func DetermineStreamType(url string) models.StreamType {
	r := regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?twitch.tv\/(videos\/(?P<videos>\w+)|(?P<livestream>\w+))`)
	if r.MatchString(url) {
		match := r.FindStringSubmatch(url)
		resMapped := make(map[string]string)

		for key, val := range r.SubexpNames() {
			if key != 0 && strings.Compare(val, "") != 0 {
				if strings.Compare(match[key], "") != 0 {
					resMapped[val] = match[key]
				}
			}
		}

		if len(resMapped) != 0 {
			for k, v := range resMapped {
				//fmt.Printf("%s: %s\n", k, v)
				switch k {
				case "livestream":
					fmt.Printf("Live stream detected. Looking up %s\n", v)
					return models.Livestream
				case "videos":
					fmt.Println("VOD detected.")
					return models.VOD
				}
			}
		} else {
			fmt.Println("No match found.")
		}

	}

	return models.Unknown
}
