package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/icza/mjpeg"
	"golang.org/x/image/font/gofont/goregular"
)

type ImgProcessor struct {
	mu     sync.Mutex
	images map[int]([]byte)
	wg     sync.WaitGroup
}

func main() {
	if len(os.Args) != 4 {
		panic("Invalid usage: Use this format: imgfeedtovideo {removefiles(true/false) pathToFeedFolder} {outputPathWithFileName}")
	}

	isRmFiles := os.Args[1] == "true"
	feedDir := os.Args[2]
	outputPath := os.Args[3]

	// Load the font
	font, _ := truetype.Parse(goregular.TTF)

	// Init video
	aw, err := mjpeg.New(outputPath, 200, 100, 1)
	if err != nil {
		log.Fatal(err)
	}

	//Get all files in dir
	fls, err := os.ReadDir(feedDir)
	if err != nil {
		log.Fatal(err)
	}
	rmFiles := []string{}
	processor := ImgProcessor{images: make(map[int][]byte)}

	//Sort by mod date
	sort.Slice(fls, func(i, j int) bool {
		fd1, _ := fls[i].Info()
		fd2, _ := fls[j].Info()
		return fd1.ModTime().Unix() < fd2.ModTime().Unix()
	})

	for i, file := range fls {
		if !file.IsDir() {
			filename := fmt.Sprintf("%s/%s", feedDir, file.Name())
			fd, _ := file.Info()

			processor.wg.Add(1)
			go processImg(filename, font, fd.ModTime().Format("2006-01-02 15:04:05"), &processor, i)

			if err != nil {
				log.Fatal(err)
			}
			rmFiles = append(rmFiles, filename)
		}
	}

	processor.wg.Wait()
	for i, _ := range fls {
		_ = i
		if processor.images[i] != nil {
			err = aw.AddFrame(processor.images[i])
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	if err := aw.Close(); err != nil {
		log.Fatal(err)
	}
	if isRmFiles {
		log.Print("Removing files")
		for _, fileName := range rmFiles {
			if err := os.Remove(fileName); err != nil {
				log.Fatal(err)
			}
		}
	}

}

func processImg(filename string, font *truetype.Font, text string, processor *ImgProcessor, index int) error {

	reader, err := os.Open(filename)
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Printf("File %s skipped as it appears to be corrupted. Decoder error: %s", filename, err.Error())
		processor.images[index] = nil
		return nil
	}
	size := img.Bounds().Size()
	var fontSize int = 22
	if size.X/50 > 22 {
		fontSize = int(size.X / 50)
	}
	imgDraw := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	draw.Draw(imgDraw, imgDraw.Bounds(), img, image.Point{}, draw.Src)

	ctx := freetype.NewContext()
	ctx.SetDst(imgDraw)
	ctx.SetSrc(image.NewUniform(color.White))
	ctx.SetFont(font)
	ctx.SetFontSize(float64(fontSize))
	ctx.SetClip(imgDraw.Bounds())

	// Draw text on the image
	pt := freetype.Pt(0, fontSize)
	_, err = ctx.DrawString(text, pt)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	var newImg bytes.Buffer
	jpeg.Encode(&newImg, imgDraw, &jpeg.Options{Quality: jpeg.DefaultQuality})

	processor.mu.Lock()
	defer processor.wg.Done()
	defer processor.mu.Unlock()
	processor.images[index] = newImg.Bytes()

	return nil
}
