package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"github.com/equoia/twiVod/utilities/downloader"
)

func main() {
	app := cli.NewApp()
	app.Name =  "twiVod"
	app.Usage = `Provide a Twitch.tv VOD URL as the first argument e.g. "twivod http://www.twitch.tv/dexteritybonus/b/638919241"`
	app.Version = "0.9.7"
	app.Action = func(c *cli.Context) {

		if c.Args().First() == "" {
			fmt.Println(c.App.Name + " " + c.App.Version)
			fmt.Println("  " + c.App.Usage)
		} else {
			fmt.Println(c.App.Name + " " + c.App.Version)
			downloader.Get(fmt.Sprintf(c.Args().First()))
		}

	}

	app.Run(os.Args)
}


