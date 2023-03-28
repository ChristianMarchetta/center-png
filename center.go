package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
)

// TolerantImage returns color.Transparent for pixels with an alpha value < tolerance, color.Opaque.
type TolerantImage struct {
	// *image.NRGBA
	image.Image
	tolerance uint8
}

// Portions of the At() function are taken from the go std/image package.
// Please refer to their licence.
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the vendor/go-licence file.
func (t *TolerantImage) At(x, y int) color.Color {
	if !(image.Point{x, y}).In(t.Bounds()) {
		return color.Transparent
	}

	c := t.Image.At(x, y)
	_, _, _, a := color.NRGBAModel.Convert(c).RGBA()

	if a <= uint32(t.tolerance) {
		return color.Transparent
	}

	return color.Opaque
}

// Frame computes a new image frame (TL corner, BR corner) so that it is centered.
func Frame(img image.Image, tolerance uint8, radius int) (image.Point, image.Point) {
	timg := &TolerantImage{img, tolerance}

	return frame(timg, radius)
	// return deadSimpleFrame(timg, radius)
}

func deadSimpleFrame(img *TolerantImage, radius int) (image.Point, image.Point) {

	tlc := image.Point{img.Bounds().Dx(), img.Bounds().Dy()}
	brc := image.Point{-1, -1}

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if img.At(x, y) == color.Opaque {
				if x < tlc.X {
					tlc.X = x
				}
				if x > brc.X {
					brc.X = x
				}
				if y < tlc.Y {
					tlc.Y = y
				}
				if y > brc.Y {
					brc.Y = y
				}
			}
		}
	}

	return tlc, brc
}

// center returns the center of the image.
func frame(img *TolerantImage, radius int) (image.Point, image.Point) {

	if img.Bounds().Dx() == 0 && img.Bounds().Dy() == 0 {
		return image.Point{}, image.Point{-1, -1}
	}

	if img.Bounds().Dx() == 0 {
		return image.Point{0, img.Bounds().Dy()}, image.Point{-1, -1}
	}

	directions := []image.Point{
		{1, 0},  // TOP
		{0, 1},  // RIGHT
		{-1, 0}, // BOTTOM
		{0, -1}, // LEFT
	}

	// The corners at which we will stop
	stopCorners := []image.Point{
		{img.Bounds().Dx() - 1, 0},
		{img.Bounds().Dx() - 1, img.Bounds().Dy() - 1},
		{0, img.Bounds().Dy() - 1},
		{0, 0},
	}

	cornerIncrs := []image.Point{
		{-1, 1},
		{-1, -1},
		{1, -1},
		{1, 1},
	}

	currDir := 0
	prevDir := 3

	// We always want the image dimensions to be even to allow for easy stop conditions in the below loop

	// if the height is odd
	if img.Bounds().Dy()%2 != 0 {
		// if the height is odd we move the top row up by one pixel,
		// but also go to the next direction to avoid computing over transparent pixels
		currDir = 1
		prevDir = 0
		stopCorners[0] = stopCorners[0].Add(image.Point{0, -1})
		stopCorners[3] = stopCorners[3].Add(image.Point{0, -1})
	}

	// if the width is odd
	if img.Bounds().Dx()%2 != 0 {
		// we move the right column right by one pixel,
		// if there was an odd heght, we are going to also skip this direction as well
		// otherwie we'll have to go throw a row for nothing.
		if currDir == 1 {
			currDir = 2
			prevDir = 1
		}
		// TODO: avoid going through a transparent column if currDir == 0

		stopCorners[0] = stopCorners[0].Add(image.Point{1, 0})
		stopCorners[1] = stopCorners[1].Add(image.Point{1, 0})
	}

	pos := stopCorners[prevDir]

	// The found corners
	tlc := image.Point{img.Bounds().Dx(), img.Bounds().Dy()}
	brc := image.Point{-1, -1}

	founds := make([]bool, 4)
	foundCount := 0

	// thresholdRectangle := image.Rect(-radius, -radius, radius, radius)
	// area := thresholdRectangle.Dx() * thresholdRectangle.Dy()
	// threshold := area / 2

	// thresholdRectOffsets := []image.Point{
	// 	{radius,}
	// }

	for {

		// if img.At(pos.X, pos.Y) == color.Opaque && satisfiesThreshold(img, thresholdRectangle.Add(image.Point{}.), threshold) {
		if img.At(pos.X, pos.Y) == color.Opaque {

			// thresholdRectangle := image.Rect(
			// 	math.Max(pos.X -radius,
			// 	-radius,
			// 	radius,
			// 	radius)
			// area := thresholdRectangle.Dx() * thresholdRectangle.Dy()
			// threshold := area / 2
			// if {

			if pos.X < tlc.X {
				tlc.X = pos.X
			}
			if pos.X > brc.X {
				brc.X = pos.X
			}
			if pos.Y < tlc.Y {
				tlc.Y = pos.Y
			}
			if pos.Y > brc.Y {
				brc.Y = pos.Y
			}

			// there cannot be anything more at the directions[currDir] than this point
			if !founds[currDir] {
				founds[currDir] = true
				foundCount += 1
				if foundCount == 4 {
					break
				}
			}

			// am I at the starting corner?
			if pos == stopCorners[prevDir] {
				// there cannot be anything more at the directions[prevDir] than this point
				if !founds[prevDir] {
					founds[prevDir] = true
					foundCount += 1
					if foundCount == 4 {
						break
					}
				}
			}

			//am I at the stopping corner?
			if pos == stopCorners[currDir] {
				nextDir := (currDir + 1) % 4
				// there cannot be anything more at the directions[nextDir] than this point
				if !founds[nextDir] {
					founds[nextDir] = true
					foundCount += 1
					if foundCount == 4 {
						break
					}
				}
			}
			// }
		}

		pos = pos.Add(directions[currDir])
		//fmt.Println("\t\t", pos)

		if pos == stopCorners[currDir] {
			prevDir = currDir
			currDir = (currDir + 1) % 4
			if currDir == 0 {
				for i := range stopCorners {
					stopCorners[i] = stopCorners[i].Add(cornerIncrs[i])
				}

				pos = stopCorners[3]

				width := stopCorners[0].X - stopCorners[3].X
				height := stopCorners[1].Y - stopCorners[0].Y
				// we're done
				if width < 0 || height < 0 {
					break
				}
			}
		}
	}

	return tlc, brc
}

func satisfiesThreshold(img *TolerantImage, rect image.Rectangle, threshold int) bool {

	visible := 0

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			if img.At(x, y) == color.Opaque {
				visible += 1
			}
		}
	}

	return visible > threshold
}

func Cut(img image.Image, tlc image.Point, brc image.Point, padding PaddingArgs) image.Image {

	paddings := convertPaddings(image.Rect(tlc.X, tlc.Y, brc.X, brc.Y), padding)

	newWidth := brc.X - tlc.X + paddings.Left + paddings.Right
	newHeight := brc.Y - tlc.Y + paddings.Top + paddings.Bottom

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := tlc.Y; y <= brc.Y; y++ {
		for x := tlc.X; x <= brc.X; x++ {
			newImg.Set(x-tlc.X+paddings.Left, y-tlc.Y+paddings.Top, img.At(x, y))
		}
	}

	return newImg
}

type PaddingPixels struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

func convertPaddings(rect image.Rectangle, paddings PaddingArgs) PaddingPixels {
	conv := func(pad PaddingArg, dim int) int {
		switch pad.Type {
		case Pixel:
			return pad.Value
		case Percent:
			return int(float64(dim*pad.Value) / 100)
		default:
			panic("invalid padding argument")
		}
	}

	return PaddingPixels{
		Top:    conv(paddings.Top, rect.Dy()),
		Right:  conv(paddings.Right, rect.Dx()),
		Bottom: conv(paddings.Bottom, rect.Dy()),
		Left:   conv(paddings.Left, rect.Dx()),
	}
}

var errEmptyImage = errors.New("empty image")
var errNotPng = errors.New("not a png image")

func Process(inFile string, outFile string, tolerance uint8, padding PaddingArgs) error {

	f, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return err
	}

	if format != "png" {
		return errNotPng
	}

	tlc, brc := Frame(img, tolerance, 0)

	if rect := image.Rect(tlc.X, tlc.Y, brc.X, brc.Y); rect.Empty() {
		return errEmptyImage
	}

	newImg := Cut(img, tlc, brc, padding)

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, newImg)

}
