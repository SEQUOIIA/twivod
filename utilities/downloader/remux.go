package downloader

import (
	"bufio"
	"fmt"
	"github.com/sequoiia/twivod/models"
	"github.com/sequoiia/twivod/utilities/stream"
	"log"
	"os"
	"os/exec"
)

func Remux(vod *models.TwitchVodOptions, ds *stream.Client) {
	var ffmpegArgs string = fmt.Sprintf("%s%s.mp4", vod.SaveFilePath, vod.Name)
	fmt.Println(fmt.Sprintf("\"%s%s\"", vod.SaveFilePath, vod.FileName))
	cmd := exec.Command("ffmpeg", "-analyzeduration", "1000000000", "-probesize", "1000000000", "-i", fmt.Sprintf("%s%s", vod.SaveFilePath, vod.FileName), "-bsf:a", "aac_adtstoasc", "-c", "copy", ffmpegArgs)
	fmt.Println(cmd.Args)
	stdout, err := cmd.StderrPipe()
	r := bufio.NewReader(stdout)
	if err != nil {
		ds.HandleErrorFatal(err)
	}
	err = cmd.Start()
	if err != nil {
		ds.HandleErrorFatal(err)
	}

	getProgress(r, ds)

	err = os.Remove(fmt.Sprintf("%s%s", vod.SaveFilePath, vod.FileName))
	if err != nil {
		ds.HandleErrorFatal(err)
	}

	if !ds.Enabled {
		log.Println("Done!")
	} else {
		ds.Handle(stream.Container{
			Status:  stream.StatusOK,
			Type:    stream.TypeStage,
			Payload: stream.StageFinished,
		})
	}
}
