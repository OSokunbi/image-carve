package carver

import (
	"image"
	"image/color"
	"image/draw"
)

// Energy calculates the energy of each pixel in the image
func Energy(img image.Image) [][]float64 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	energies := make([][]float64, height)

	for y := 0; y < height; y++ {
		energies[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			dx := colorDiff(img.At(x-1, y), img.At(x+1, y))
			dy := colorDiff(img.At(x, y-1), img.At(x, y+1))
			energies[y][x] = dx + dy
		}
	}
	return energies
}

// colorDiff calculates the difference between two colors
func colorDiff(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return float64((r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2))
}

// findVerticalSeam finds the vertical seam with the lowest energy
func findVerticalSeam(energies [][]float64) []int {
	height, width := len(energies), len(energies[0])
	seam := make([]int, height)
	dp := make([][]float64, height)

	for i := range dp {
		dp[i] = make([]float64, width)
		for j := range dp[i] {
			dp[i][j] = 1e18
		}
	}

	for j := 0; j < width; j++ {
		dp[0][j] = energies[0][j]
	}

	for i := 1; i < height; i++ {
		for j := 0; j < width; j++ {
			dp[i][j] = energies[i][j] + dp[i-1][j]
			if j > 0 && dp[i-1][j-1]+energies[i][j] < dp[i][j] {
				dp[i][j] = dp[i-1][j-1] + energies[i][j]
			}
			if j < width-1 && dp[i-1][j+1]+energies[i][j] < dp[i][j] {
				dp[i][j] = dp[i-1][j+1] + energies[i][j]
			}
		}
	}

	minIdx := 0
	for j := 1; j < width; j++ {
		if dp[height-1][j] < dp[height-1][minIdx] {
			minIdx = j
		}
	}

	seam[height-1] = minIdx
	for i := height - 2; i >= 0; i-- {
		seam[i] = seam[i+1]
		if seam[i+1] > 0 && dp[i][seam[i+1]-1] < dp[i][seam[i]] {
			seam[i] = seam[i+1] - 1
		}
		if seam[i+1] < width-1 && dp[i][seam[i+1]+1] < dp[i][seam[i]] {
			seam[i] = seam[i+1] + 1
		}
	}

	return seam
}

// removeVerticalSeam removes the specified vertical seam from the image
func removeVerticalSeam(img *image.RGBA, seam []int) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := 0; y < height; y++ {
		for x := seam[y]; x < width-1; x++ {
			img.Set(x, y, img.At(x+1, y))
		}
	}
	img.Rect.Max.X--
}

// Rotate90 rotates the image 90 degrees clockwise
func Rotate90(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	rotated := image.NewRGBA(image.Rect(0, 0, height, width))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rotated.Set(height-y-1, x, img.At(x, y))
		}
	}

	return rotated
}

// RotateMinus90 rotates the image 90 degrees counter-clockwise
func RotateMinus90(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	rotated := image.NewRGBA(image.Rect(0, 0, height, width))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rotated.Set(y, width-x-1, img.At(x, y))
		}
	}

	return rotated
}

// Carve carves an image based on an initial image, a targetheight, and a width
func Carve(initial image.Image, targetWidth, targetHeight int) image.Image {
    rgba := image.NewRGBA(initial.Bounds())
    draw.Draw(rgba, initial.Bounds(), initial, image.Point{}, draw.Src)

    energies := Energy(rgba)

    for rgba.Bounds().Dx() > targetWidth {
        seam := findVerticalSeam(energies)
        removeVerticalSeam(rgba, seam)
        energies = Energy(rgba)
    }

    rotated := Rotate90(rgba)
    energies = Energy(rotated)

    for rotated.Bounds().Dx() > targetHeight {
        seam := findVerticalSeam(energies)
        removeVerticalSeam(rotated, seam)
        energies = Energy(rotated)
    }

    resized := RotateMinus90(rotated)

    return resized
}
