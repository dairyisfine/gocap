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

	// echo server that statically serves the /media directory
	server = echo.New()
	// server.Use(middleware.Logger())
	// server.Use(middleware.Recover())

	// serve the media directory browsable
	server.Static("/media", "")

	server.GET("/files", func(c echo.Context) error {
		reloadMediaFiles()
		return c.JSON(200, mediaFiles)
	})
	server.GET("/videodevices", func(c echo.Context) error {
		devices := ffmpeg.GetVideoDevices()
		return c.JSON(200, devices)
	})
	server.GET("/createthumbnail/:device", func(c echo.Context) error {
		device := c.Param("device")
		ffmpeg.CreateThumbnail(device)
		return c.String(200, "Thumbnail created for: "+device)
	})

	// start the server
	server.Start(":1042")

	// print a message to the console
	fmt.Println("Server started on port 1042")
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
