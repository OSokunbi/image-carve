package main

import (
  "flag"
  "os"
  "image"
  "fmt"
  "image/jpeg"
  "github.com/OSokunbi/image-carve/carver"
)

func main() {
    width := flag.Int("width", 0, "Width of the resized image")
    height := flag.Int("height", 0, "Height of the resized image")
    inputFile := flag.String("filename", "", "Input image file")
    flag.Parse()

    if *width <= 0 || *height <= 0 {
        fmt.Println("Please provide valid width and height")
        return
    }

    if *inputFile == "" {
         fmt.Println("Please provide the input image file")
         return
     }

    file, err := os.Open(*inputFile)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        panic(err)
    }

    resizedImage := carver.Carve(img, *width, *height)

    outputFile, err := os.Create("output.jpg")
    if err != nil {
        panic(err)
    }
    defer outputFile.Close()

    jpeg.Encode(outputFile, resizedImage, nil)
}
