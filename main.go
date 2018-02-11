package main

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/sequoiia/twivod/cmd"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("twivod")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath(os.Getenv("APPDATA"))

	err := viper.ReadInConfig()
	if err != nil {
		defaultConfig := []byte("twitchclientid: undefined\n")
		if runtime.GOOS == "windows" {
			err := ioutil.WriteFile(os.Getenv("APPDATA")+"\\twivod.yaml", defaultConfig, 0770)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := ioutil.WriteFile(os.Getenv("HOME")+"/.config/twivod.yaml", defaultConfig, 0770)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = viper.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	cmd.Execute()
}
