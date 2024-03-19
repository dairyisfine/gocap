package main

import (
	"github.com/dairyisfine/gocap/ffmpeg"
)

func main() {
	ffmpeg.FfmpegStart()
	web()
}
