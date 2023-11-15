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
	// Load the font
	font, _ := truetype.Parse(goregular.TTF)

	feedDir := "./examples"
	// Video size: 200x100 pixels, FPS: 2
	aw, err := mjpeg.New("test.avi", 200, 100, 1)
	if err != nil {
		log.Fatal(err)
	}

	//Get all files in dir
	fls, err := os.ReadDir(feedDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fls {
		if !file.IsDir() {
			filename := fmt.Sprintf("%s/%s", feedDir, file.Name())
			fd, _ := file.Info()

			reader, err := os.Open(filename)
			img, _, _ := image.Decode(reader)
			newImg, err := processImg(img, font, fd.ModTime().Format("2006-01-02 15:04:05"))
			if err != nil {
				log.Fatal(err)
			}

			err = aw.AddFrame(newImg)
			if err != nil {
				log.Fatal(err)
			}
			reader.Close()
		}
	}

	if err := aw.Close(); err != nil {
		log.Fatal(err)
	}
}

func processImg(img image.Image, font *truetype.Font, text string) ([]byte, error) {
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
