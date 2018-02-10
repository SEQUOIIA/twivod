package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sequoiia/twiVod/models"
	"github.com/sequoiia/twiVod/utilities/downloader"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "twiVod"
	app.Usage = `Provide a Twitch.tv VOD URL as the first argument e.g. "twivod 4 http://www.twitch.tv/videos/638919241"`
	app.Version = "0.9.7"
	app.Action = func(c *cli.Context) error {

		if c.Args().First() == "" {
			fmt.Println(c.App.Name + " " + c.App.Version)
			fmt.Println("  " + c.App.Usage)
		} else {
			fmt.Println(c.App.Name + " " + c.App.Version)
			var vodOptions *models.TwitchVodOptions = &models.TwitchVodOptions{}
			vodOptions.Url = fmt.Sprintf(c.Args().Get(1))
			param0, err := strconv.Atoi(c.Args().Get(0))
			if err != nil {
				log.Fatal("Invalid parameters.")
			}
			vodOptions.MaxConcurrentDownloads = param0

			err = downloader.Download(vodOptions)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(vodOptions.FileName)
			log.Println(vodOptions.Name)
			downloader.Remux(vodOptions)
		}
		return nil
	}

	app.Run(os.Args)
}
