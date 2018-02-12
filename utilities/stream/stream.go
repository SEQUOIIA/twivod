package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

// Status is used for determining whether something went as expected, or if something unexpected occured.
type Status string

// StatusOK use in case everything went as expected.
const StatusOK Status = "ok"

// StatusError use if something unexpected occured.
const StatusError Status = "error"

// Type is used to determine what kind of payload is to be expected.
type Type string

// TypeStage is used to determine what stage twiVod is at.
const TypeStage Type = "stage"

// TypeDownloadProgress is to be used if payload is DownloadProgress
const TypeDownloadProgress Type = "downloadProgress"

// TypeRemuxProgress is to be used if payload is RemuxProgress
const TypeRemuxProgress Type = "remuxProgress"

// Stage contains what twiVod is currently doing
type Stage string

// StageStarted is to be used as soon as the application is started
const StageStarted Stage = "started"

// StageDownload is to be used just before the download process begins
const StageDownload Stage = "download"

// StageRemux is to be used after StageDownload has finished
const StageRemux Stage = "remux"

// StageFinished is to be used after StageRemux has finished
const StageFinished Stage = "finished"

// Container is the data structure used for storing and transferring data through DataStream
type Container struct {
	Status  Status
	Type    Type
	Payload interface{}
}

// Client can handle Container objects, and encode them to JSON, for easy consumption through the DataStream flag
type Client struct {
	Enabled    bool
	handleLock sync.Mutex
}

// Handle :
// Takes a Container as parameter, encodes that object to JSON, and prints it to stdout
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

// HandleErrorFatal can be used instead of log.Fatal(), so if one uses the DataStream flag, one can still get the error message encapsulated in a Container
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

// DownloadProgress - Used for download stats
type DownloadProgress struct {
	TotalSegments  int
	CurrentSegment int
}

// RemuxProgress - WIP
type RemuxProgress struct {
}
