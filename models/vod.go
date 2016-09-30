package models

import "io"

type TwitchVodOptions struct {
	Name string
	FileName string
	Url string
	MaxConcurrentDownloads int
	SaveFilePath string
	Writer io.Writer
}