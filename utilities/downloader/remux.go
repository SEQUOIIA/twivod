package downloader

import (
	"fmt"
	"os/exec"
	"bufio"
	"log"
	"github.com/sequoiia/twiVod/models"
	"os"
)

func Remux(vod *models.TwitchVodOptions) {
	var ffmpegArgs string = fmt.Sprintf("%s.mp4", vod.Name)
	cmd := exec.Command("ffmpeg", "-analyzeduration", "1000000000", "-probesize", "1000000000", "-i" , vod.FileName, "-bsf:a", "aac_adtstoasc", "-c", "copy", ffmpegArgs)
	stdout, err := cmd.StderrPipe()
	r := bufio.NewReader(stdout)
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start(); if err != nil {
		log.Fatal(err)
	}

	log.Println("Should've started here.")

	getProgress(r, 4564564)

	err = os.Remove(vod.FileName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
}