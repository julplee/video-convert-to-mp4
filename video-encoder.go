package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var files []string

	root := "./video-to-encode"
	supportedExtensions := map[string]struct{}{
		".ts":  {},
		".mp4": {},
		".wmv": {},
		".avi": {},
		".mkv": {},
	}

	log.Println("Reading the video to encode directory...")
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		extension := strings.ToLower(filepath.Ext(path))
		if info.IsDir() || !isSupportedExtension(extension, supportedExtensions) {
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
		if err := encodeVideoToMP4(file); err != nil {
			log.Println("Skipping rename because encoding failed:", file)
			continue
		}
		renameOriginFile(file)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Encoding done")
	reader.ReadString('\n')
}

func isSupportedExtension(extension string, supportedExtensions map[string]struct{}) bool {
	_, ok := supportedExtensions[extension]
	return ok
}

func encodeVideoToMP4(inputVideoFilepath string) error {
	var filename = filepath.Base(inputVideoFilepath)
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]

	// use ffmpeg to encode video to MP4
	ffmpegPath, err := resolveFFmpegPath()
	if err != nil {
		log.Println("ERROR resolving ffmpeg:", err)
		return err
	}
	outputPath := filepath.Join("video-encoded", name+".mp4")

	cmd := exec.Command(
		ffmpegPath,
		"-i", filepath.Clean(inputVideoFilepath),
		"-preset", "slow",
		"-c:a", "aac",
		"-c:v", "libx264",
		"-maxrate", "3.5M",
		"-bufsize", "1.5M",
		"-profile:v", "main",
		outputPath)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Println("Encode video " + inputVideoFilepath)

	start := time.Now()
	err = cmd.Run()
	elapsed := time.Since(start)

	log.Printf("Recording took %s", elapsed)

	if err != nil {
		log.Println("ERROR: ", inputVideoFilepath, stderr.String())
		log.Println("ERROR executing ffmpeg:", err)
		return err
	}

	return nil
}

func resolveFFmpegPath() (string, error) {
	candidates := []string{"ffmpeg", "ffmpeg.exe"}

	for _, candidate := range candidates {
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}

	for _, candidate := range candidates {
		if path, err := filepath.Abs(candidate); err == nil {
			if _, statErr := os.Stat(path); statErr == nil {
				return path, nil
			}
		}
	}

	return "", errors.New("ffmpeg executable not found in PATH or project root")
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
