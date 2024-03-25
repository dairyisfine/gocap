package ffmpeg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var cwd string
var err error

var activeRecording bool
var activeRecordingProcess *exec.Cmd

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
		fmt.Println("cannot retrieve wireless IP")
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
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/") {
			deviceName := strings.Replace(line, "/dev/", "", 1)
			devices = append(devices, deviceName)
		}
	}
	return devices
}

func CreateThumbnail(device string) error {
	fmt.Println("Creating thumbnail for: ", device)
	cmd := exec.Command("ffmpeg", "-f", "v4l2", "-i", "/dev/"+device, "-vf", "fps=1/2", "-frames:v", "120", "-update", "1", "-vf", "scale=320:240", "thumbnails/"+device+".jpg", "-y")
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("Thumbnail created")
	return nil
}

func StartCapture(device string) error {
	fmt.Println("Starting recording for: ", device)
	fileName := device + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".mkv"

	activeRecordingProcess = exec.Command("ffmpeg", "-f", "v4l2", "-framerate", "30", "-video_size", "640x480", "-i", "/dev/"+device, fileName, "-vf", "fps=1/3,scale=320:240", "-update", "1", "thumbnails/"+device+".jpg", "-y")

	activeRecordingProcess.Start()
	time.Sleep(2 * time.Second)
	_, err := os.Stat(fileName)
	if err != nil {
		activeRecordingProcess.Process.Kill()
		activeRecordingProcess.Wait()
		return err
	}
	activeRecording = true
	fmt.Println("Recording started")
	return nil
}

func StopCapture() error {
	fmt.Println("Stopping recording")
	err := activeRecordingProcess.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}
	activeRecordingProcess.Wait()
	activeRecording = false
	return nil
}

func IsActiveRecording() bool {
	return activeRecording
}

