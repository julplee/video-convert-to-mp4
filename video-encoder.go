package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var files []string

	root := "./video-to-encode"

	log.Println("Reading the video to encode directory...")
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || (filepath.Ext(path) != ".ts" && filepath.Ext(path) != ".mp4" && filepath.Ext(path) != ".wmv" && filepath.Ext(path) != ".avi") {
			if path != root {
				log.Println("WARNING: The file " + path + " will not be encoded")
			}
			return nil
		}

		log.Println("The file " + path + " will be encoded")
		files = append(files, path)
		return nil
	})

	if err != nil {
		log.Println("ERROR: ", err)
	}

	log.Println("Now let's encode each video:")
	for _, file := range files {
		encodeVideoToMP4(file)
		renameOriginFile(file)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Encoding done")
	reader.ReadString('\n')
}

func encodeVideoToMP4(inputVideoFilepath string) {
	var filename = filepath.Base(inputVideoFilepath)
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]

	// use ffmpeg to encode video to MP4
	cmd := exec.Command(
		"ffmpeg.exe",
		"-i", inputVideoFilepath,
		"-preset", "slow",
		"-c:a", "aac",
		"-c:v", "libx264",
		"-maxrate", "3.5M",
		"-bufsize", "1.5M",
		"-profile:v", "main",
		"video-encoded\\"+name+".mp4")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Println("Encode video " + inputVideoFilepath)
	err := cmd.Run()

	if err != nil {
		log.Println("ERROR: ", inputVideoFilepath, stderr.String())
	}
}

func renameOriginFile(inputVideoFilepath string) {
	var filename = filepath.Base(inputVideoFilepath)
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	var newfilevideopath = strings.Replace(inputVideoFilepath, name, name+" orig", 1)

	e := os.Rename(inputVideoFilepath, newfilevideopath)
	if e != nil {
		log.Println("ERROR: ", e)
	}
}
