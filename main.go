package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"flag"

	"github.com/aybabtme/rgbterm"
	"github.com/nfnt/resize"

	_ "image/jpeg"
	_ "image/png"
)

// Constants
const scaleFactor float64 = 0.5
const charWidth int = 10
const charHeight int = 18
const resizeWidth uint = 100

// Flag variables
var imageFileName string
var isColored bool


func getImageFromFile(dir string) (image.Image, string, error) {
	imgFile, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer imgFile.Close()

	return image.Decode(imgFile)
}

// Get character from luminosity/grayscale value
func getChar(val float64) byte {
	asciiMatrix := []byte(" `^\":;Il!iXYUJCLQ0OZ#MW8B@$")
 
	divisor := float64(256)/float64(len(asciiMatrix))

	return asciiMatrix[int(math.Floor(float64(val)/divisor))]
}

func getLuminosityPt(x int, y int, img image.Image) float64 {
	r, g, b, _ := img.At(x, y).RGBA()
	return (0.2126*float64(r>>8) + 0.7152*float64(g>>8) + 0.0722*float64(b>>8))
}

func getGrayscalePt(x int, y int, img image.Image) float64 {
	r, g, b, _ := img.At(x, y).RGBA()
	return (float64(r>>8 + g>>8 + b>>8)/3)
}

func printImage(img[][] string) {
	for _, dY := range img {
		for _, dX := range dY {
			fmt.Printf("%s ", dX)
		} 
		fmt.Println()
	}
}

func applyColor(r uint8, g uint8, b uint8, character byte) string {
	return rgbterm.FgString(string([]byte{character}), r, g, b)
}

func init() {
	const (
		defaultImageFileName = "images/bonin.jpg"
		defaultIsColored = false
	)
	flag.StringVar(&imageFileName,
		"f",
		defaultImageFileName,
		"Image file")
	flag.BoolVar(&isColored,
		"c",
		defaultIsColored,
		"Color the text")
}

func main() {
	// Parse flags
	flag.Parse()

	// Get image
	img, _, err := getImageFromFile(imageFileName)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()

	// Resize to save comp power
	img = resize.Resize(resizeWidth, 0, img, resize.Lanczos3)

	bounds = img.Bounds()

	var w int = bounds.Max.X
	var h int = bounds.Max.Y

	// Initialize new image
	asciiImage := make([][]string, h)
	for i := range asciiImage {
		asciiImage[i] = make([]string, w)
	}

	// Convert to ascii
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			lum := getLuminosityPt(x, y, img)
			character := getChar(lum)

			r, g, b, _ := color.NRGBAModel.Convert(img.At(x, y)).RGBA()

			if isColored {
				asciiImage[y][x] = applyColor(uint8(r), uint8(g), uint8(b), character)
			} else {
				asciiImage[y][x] = string(character)
			}
		}
	}

	// // Print the results.
	printImage(asciiImage)
}