package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"github.com/sequoiia/twiVod/utilities/downloader"
	"github.com/sequoiia/twiVod/models"
	"log"
)

func main() {
	app := cli.NewApp()
	app.Name =  "twiVod"
	app.Usage = `Provide a Twitch.tv VOD URL as the first argument e.g. "twivod http://www.twitch.tv/dexteritybonus/b/638919241"`
	app.Version = "0.9.7"
	app.Action = func(c *cli.Context) error{

		if c.Args().First() == "" {
			fmt.Println(c.App.Name + " " + c.App.Version)
			fmt.Println("  " + c.App.Usage)
		} else {
			fmt.Println(c.App.Name + " " + c.App.Version)
			var vodOptions *models.TwitchVodOptions = new(models.TwitchVodOptions)
			vodOptions.Url = fmt.Sprintf(c.Args().First())
			vodOptions.MaxConcurrentDownloads = 4
			err := downloader.Download(vodOptions)
			if err !=  nil {
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


