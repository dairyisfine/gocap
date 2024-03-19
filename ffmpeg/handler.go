package ffmpeg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var cwd string
var err error

func FfmpegStart() {
	fmt.Println("ffmpegStart")
	os.Chdir("output/media")
	cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("current directory: ", cwd)
	version, err := exec.Command("ffmpeg", "-version").Output()
	if err != nil {
		log.Fatal(err)
	}
	err = exec.Command("v4l2-ctl", "--list-devices").Run()
	if err != nil {
		log.Fatal(err)
	}
	err = exec.Command("arecord", "-L").Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(version))
}

func GetVideoDevices() []string {
	cmdOutput, err := exec.Command("v4l2-ctl", "--list-devices").Output()
	if err != nil {
		log.Fatal(err)
	}

	devices := []string{}
	lines := strings.Split(string(cmdOutput), "\n")
	for _, line := range lines {
		// remove whitespace
		line = strings.TrimSpace(line)
		// if the first character is a /, it's a device
		if strings.HasPrefix(line, "/") {
			deviceName := strings.Replace(line, "/dev/", "", 1)
			devices = append(devices, deviceName)
		}
	}
	return devices
}

func CreateThumbnail(device string) {
	fmt.Println("Creating thumbnail for: ", device)
	cmd := exec.Command("ffmpeg", "-f", "v4l2", "-i", "/dev/"+device, "-frames:v", "1", "-vf", "scale=320:240", device+".jpg", "-y")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Thumbnail created")
}