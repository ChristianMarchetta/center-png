package main

import (
	"image"
	"image/color"
)

// TolerantImage returns color.Transparent for pixels with an alpha value < tolerance, color.Opaque.
type TolerantImage struct {
	*image.NRGBA
	tolerance uint8
}

// Portions of the At() function are taken from the go std/image package.
// Please refer to their licence.
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the vendor/go-licence file.
func (t *TolerantImage) At(x, y int) color.Color {
	if !(image.Point{x, y}).In(t.NRGBA.Rect) {
		return color.Transparent
	}

	i := t.NRGBA.PixOffset(x, y)
	a := t.NRGBA.Pix[i+3]

	if a < t.tolerance {
		return color.Transparent
	}

	return color.Opaque
}

// Frame computes a new image frame (TL corner, BR corner) so that it is centered.
func Frame(img image.Image, tolerance uint8, radius int) (image.Point, image.Point) {
	timg := &TolerantImage{image.NewNRGBA(img.Bounds()), tolerance}

	return frame(timg, radius)
}

// center returns the center of the image.
func frame(img *TolerantImage, radius int) (image.Point, image.Point) {

	directions := []image.Point{
		{1, 0},  // TOP
		{0, 1},  // RIGHT
		{-1, 0}, // BOTTOM
		{0, -1}, // LEFT
	}

	// The corners at which we will stop
	stopCorners := []image.Point{
		{img.Rect.Dx() - 1, 0},
		{img.Rect.Dx() - 1, img.Rect.Dy() - 1},
		{0, img.Rect.Dy() - 1},
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
	if img.Rect.Dy()%2 != 0 {
		// if the height is odd we move the top row up by one pixel,
		// but also go to the next direction to avoid computing over transparent pixels
		currDir = 1
		prevDir = 0
		stopCorners[0] = stopCorners[0].Add(image.Point{0, -1})
		stopCorners[3] = stopCorners[3].Add(image.Point{0, -1})
	}

	// if the width is odd
	if img.Rect.Dx()%2 != 0 {
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
	tlc := image.Point{img.Rect.Dx(), img.Rect.Dy()}
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
		}

		pos = pos.Add(directions[currDir])

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

// func satisfiesThreshold(img *TolerantImage, rect image.Rectangle, threshold int) bool {

// 	visible := 0

// 	for y := rect.Min.Y; y < rect.Max.Y; y++ {
// 		for x := rect.Min.X; x < rect.Max.X; x++ {
// 			if img.At(x, y) == color.Opaque {
// 				visible += 1
// 			}
// 		}
// 	}

// 	return visible > threshold
// }
