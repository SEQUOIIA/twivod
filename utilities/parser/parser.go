package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sequoiia/twiVod/models"
)

func VodInfo(url string) models.VODinfo {
	urlType := DetermineStreamType(url)

	if urlType.Type == models.VOD {
		return models.VODinfo{Channel: "", Type: urlType.Type, ID: urlType.Value}
	} else {
		return models.VODinfo{Type: models.Unknown}
	}
}

// DetermineStreamType runs a regex expression that determines whether a live stream or a VOD has been given
func DetermineStreamType(url string) models.StreamTypeResult {
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
					return models.StreamTypeResult{Type: models.Livestream}
				case "videos":
					fmt.Println("VOD detected.")
					return models.StreamTypeResult{Type: models.VOD, Value: v}
				}
			}
		} else {
			fmt.Println("No match found.")
		}

	}

	return models.StreamTypeResult{Type: models.Unknown}
}
