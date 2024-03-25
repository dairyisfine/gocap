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

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func web() {

	server = echo.New()

	// allow all origins
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			return next(c)
		}
	})

	// media directory serves captures
	server.Static("/media", "")

	// ui directory serves the web interface
	server.Static("/", "../../ui/dist/")

	// route for getting file list
	server.GET("/mediafilelist", func(c echo.Context) error {
		reloadMediaFiles()
		// filter "thumbnails"
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
		err := ffmpeg.CreateThumbnail(device)
		if err != nil {
			return c.JSON(500, Response{Success: false, Message: "Failed to create thumbnail for: "+device})
		}
		return c.JSON(200, Response{Success: true, Message: "Thumbnail created for: "+device})
	})

	server.GET("/thumbnail/:device", func(c echo.Context) error {
		device := c.Param("device")
		return c.File("thumbnails/"+device+".jpg")
	})

	// route to start capturing from the requested device
	server.GET("/startcapture/:device", func(c echo.Context) error {
		device := c.Param("device")
		err := ffmpeg.StartCapture(device)
		if err != nil {
			return c.JSON(500, Response{Success: false, Message: "Failed to start capture for: "+device+": "+err.Error()})
		}
		return c.JSON(200, Response{Success: true, Message: "Capture started for: "+device})
	})

	server.GET("/activerecording", func(c echo.Context) error {
		return c.JSON(200, Response{Success: ffmpeg.IsActiveRecording(), Message: ""})
	})


	// route to stop capturing from the requested device
	server.GET("/stopcapture", func(c echo.Context) error {
		err := ffmpeg.StopCapture()
		if err != nil {
			return c.JSON(500, Response{Success: false, Message: "Failed to stop capture"})
		}
		return c.JSON(200, Response{Success: true, Message: "Capture stopped"})
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
		if !file.IsDir() {
			mediaFiles = append(mediaFiles, file.Name())
		}
	}
}
