package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sequoiia/twivod/models"
	"github.com/sequoiia/twivod/utilities/downloader"
	"github.com/spf13/cobra"
)

var (
	Verbose bool
	ConcurrentDownloads int
)

var rootCmd = &cobra.Command{
	Use:   "twivod https://www.twitch.tv/videos/227241536",
	Short: "twiVod is a simple tool for downloading Twitch VODs and live streams.",
	Long: `twiVod is a simple tool for downloading Twitch VODs and live streams.`,
	Run: func(cmd *cobra.Command, args []string) {
		vodOptions := &models.TwitchVodOptions{Url: args[0], MaxConcurrentDownloads: ConcurrentDownloads}

		err := downloader.Download(vodOptions)
		if err != nil {
			log.Fatal(err)
		}

		downloader.Remux(vodOptions)
	},
	Args: func(cmd *cobra.Command, args[]string) error{
		if len(args) < 1 {
			return errors.New("Requires at least one arg(url)")
		}
		return nil
	},
}

func Execute() {
	Init()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Init() {
	// rootCmd
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.LocalFlags().IntVarP(&ConcurrentDownloads, "concurrent-download", "c", 1, "The amount of video segments downloaded at the same time.")
}