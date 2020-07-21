package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"strings"

	"github.com/nfnt/resize"

	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
)

const scaleFactor float64 = 0.5
const charWidth int = 10;
const charHeight int = 18;
const resizeWidth uint = 140;

func getImageFromFile(dir string) (image.Image, string, error) {
	imgFile, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer imgFile.Close()

	return image.Decode(imgFile)
}

// Get character from luminosity/grayscale value
func getChar(val float64) string {
	asciiMatrix := strings.Split("`^\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$", "")
 
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

func main() {
	// Get image
	img, _, err := getImageFromFile("pic.jpg")
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()

	// Resize to save comp power
	img = resize.Resize(resizeWidth, 0, img, resize.Lanczos3)

	bounds = img.Bounds()
	
	// Test resized
	out, err := os.Create("test_resized.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, img, nil)

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

			asciiImage[y][x] = getChar(lum) 
		}
	}

	// // Print the results.
	printImage(asciiImage)
}