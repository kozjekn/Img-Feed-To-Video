package main

import (
	"fmt"
	"log"
	"os"

	"github.com/icza/mjpeg"
)

func main() {
	feedDir := "./examples"
	// Video size: 200x100 pixels, FPS: 2
	aw, err := mjpeg.New("test.avi", 200, 100, 2)
	if err != nil {
		log.Fatal(err)
	}

	fls, err := os.ReadDir(feedDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fls {
		if !file.IsDir() {
			filename := fmt.Sprintf("%s/%s", feedDir, file.Name())
			data, err := os.ReadFile(filename)

			if err != nil {
				log.Fatal(err)
			}

			err = aw.AddFrame(data)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := aw.Close(); err != nil {
		log.Fatal(err)
	}
}
