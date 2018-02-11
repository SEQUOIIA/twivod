package stream

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
