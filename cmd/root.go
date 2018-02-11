package cmd

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/sequoiia/twivod/models"
	"github.com/sequoiia/twivod/utilities/downloader"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Verbose             bool
	ConcurrentDownloads int
	BandwidthLimit      int64
	PrintVersion        bool
	SetTwitchClientID   string
	TwitchClientID      string
)

var rootCmd = &cobra.Command{
	Use:   "twivod https://www.twitch.tv/videos/227241536",
	Short: "twiVod is a simple tool for downloading Twitch VODs and live streams.",
	Long:  `twiVod is a simple tool for downloading Twitch VODs and live streams.`,
	Run: func(cmd *cobra.Command, args []string) {
		vodOptions := &models.TwitchVodOptions{Url: args[0], MaxConcurrentDownloads: ConcurrentDownloads}

		if BandwidthLimit != math.MaxInt64 {
			BandwidthLimit = BandwidthLimit * 1000
		}

		if viper.GetString("twitchclientid") == "undefined" {
			if TwitchClientID == "" {
				fmt.Println("Twitch client-id has not been set. Run 'twivod --help' for further help.")
				os.Exit(1)
			}

		} else {
			models.TwitchConfig.Client_id = viper.GetString("twitchclientid")
		}

		if TwitchClientID != "" {
			models.TwitchConfig.Client_id = TwitchClientID
		}

		err := downloader.Download(vodOptions, BandwidthLimit)
		if err != nil {
			log.Fatal(err)
		}

		downloader.Remux(vodOptions)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if PrintVersion {
			versionCmd.Run(cmd, args)
			os.Exit(1)
		}

		if SetTwitchClientID != "" {
			viper.Set("twitchclientid", SetTwitchClientID)
			err := viper.WriteConfig()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Client-ID saved in config.")
			os.Exit(1)
		}

		if len(args) < 1 {
			return errors.New("requires at least one arg(url)")
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
	rootCmd.PersistentFlags().BoolVar(&PrintVersion, "version", false, "Print software version")
	rootCmd.PersistentFlags().IntVarP(&ConcurrentDownloads, "concurrent-download", "c", 1, "The amount of video segments downloaded at the same time.")
	rootCmd.PersistentFlags().Int64VarP(&BandwidthLimit, "bwlimit", "b", math.MaxInt64, "Limits download speed, value is in kb/s")
	rootCmd.PersistentFlags().StringVar(&SetTwitchClientID, "set-client-id", "", "Set client-id in config")
	rootCmd.PersistentFlags().StringVar(&TwitchClientID, "client-id", "", "Use client-id for only this command")

	rootCmd.AddCommand(versionCmd)
}
