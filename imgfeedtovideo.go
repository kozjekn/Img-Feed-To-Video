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

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/icza/mjpeg"
	"golang.org/x/image/font/gofont/goregular"
)

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
	for _, file := range fls {
		if !file.IsDir() {
			filename := fmt.Sprintf("%s/%s", feedDir, file.Name())
			fd, _ := file.Info()

			reader, err := os.Open(filename)
			newImg, err := processImg(reader, font, fd.ModTime().Format("2006-01-02 15:04:05"))
			reader.Close()

			if err != nil {
				log.Fatal(err)
			}

			err = aw.AddFrame(newImg)
			if err != nil {
				log.Fatal(err)
			}
			rmFiles = append(rmFiles, filename)
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

func processImg(reader *os.File, font *truetype.Font, text string) ([]byte, error) {
	img, _, _ := image.Decode(reader)
	const fontSize int = 70
	size := img.Bounds().Size()
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
	_, err := ctx.DrawString(text, pt)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var newImg bytes.Buffer
	jpeg.Encode(&newImg, imgDraw, &jpeg.Options{Quality: jpeg.DefaultQuality})
	return newImg.Bytes(), nil
}
