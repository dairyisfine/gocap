package main

import (
	"fmt"
	"os"

	"github.com/dairyisfine/gocap/ffmpeg"
	"github.com/labstack/echo/v4"
)

var server *echo.Echo
var mediaDir = "output/media"
var mediaFiles []string

func web() {

	server = echo.New()

	// media directory serves captures
	server.Static("/media", mediaDir)

	// route for getting file list
	server.GET("/mediafilelist", func(c echo.Context) error {
		reloadMediaFiles()
		return c.JSON(200, mediaFiles)
	})

	// route for getting video devices
	server.GET("/videodevices", func(c echo.Context) error {
		devices := ffmpeg.GetVideoDevices()
		return c.JSON(200, devices)
	})

	// route to capture a thumbnail from the requested device
	server.GET("/createthumbnail/:device", func(c echo.Context) error {
		device := c.Param("device")
		ffmpeg.CreateThumbnail(device)
		return c.String(200, "Thumbnail created for: "+device)
	})

	// start the server
	fmt.Println("Server starting on "+ffmpeg.GetWlan0Ip()+":80")
	server.Start(":80")
}

func reloadMediaFiles() {
	mediaFiles = []string{}
	files, err := os.ReadDir("./")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		mediaFiles = append(mediaFiles, file.Name())
	}
}
