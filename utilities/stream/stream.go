package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type Status string

const StatusOK Status = "ok"
const StatusError Status = "error"

type Type string

const TypeStage Type = "stage"
const TypeDownloadProgress Type = "downloadProgress"
const TypeRemuxProgress Type = "remuxProgress"

type Stage string

const StageStarted Stage = "started"
const StageDownload Stage = "download"
const StageRemux Stage = "remux"
const StageFinished Stage = "finished"

type Container struct {
	Status  Status
	Type    Type
	Payload interface{}
}

type Client struct {
	Enabled    bool
	handleLock sync.Mutex
}

func (c *Client) Handle(con Container) {
	c.handleLock.Lock()

	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(&con)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(buf.String())

	c.handleLock.Unlock()
}

func (c *Client) HandleErrorFatal(err error) {
	if !c.Enabled {
		log.Fatal(err)
	} else {
		c.Handle(Container{
			Status:  StatusError,
			Payload: err.Error(),
		})
		os.Exit(1)
	}
}

type DownloadProgress struct {
	TotalSegments  int
	CurrentSegment int
}
