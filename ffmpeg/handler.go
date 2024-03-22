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
	fmt.Println(string(version))
	devices, err := exec.Command("v4l2-ctl").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(devices))
	alsaDev, err := exec.Command("arecord", "-L").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(alsaDev))
}

func GetWlan0Ip() string {
	cmdOutput, err := exec.Command("ip", "addr", "show", "wlan0").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(cmdOutput), "\n")
	for _, line := range lines {
		if strings.Contains(line, "inet ") {
			ip := strings.Split(strings.Split(line, " ")[5], "/")[0]
			return ip
		}
	}
	return ""
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

func CreateThumbnail(device string) error {
	fmt.Println("Creating thumbnail for: ", device)
	cmd := exec.Command("ffmpeg", "-f", "v4l2", "-i", "/dev/"+device, "-frames:v", "1", "-vf", "scale=320:240", device+".jpg", "-y")
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("Thumbnail created")
	return nil
}