package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"testing"
)

func newImageFromBitMap(bitmap [][]uint8) *image.NRGBA {
	width, height := 0, len(bitmap)

	if height > 0 {
		width = len(bitmap[0])
	}
	r := image.Rect(0, 0, width, height)
	img := image.NewNRGBA(r)

	for y := 0; y < len(bitmap); y++ {
		for x := 0; x < len(bitmap[y]); x++ {
			if bitmap[y][x] > 0 {
				img.Set(x, y, color.Opaque)
			} else {
				img.Set(x, y, color.Transparent)
			}
		}
	}

	return img
}

var squareOdd = [][]uint8{
	{1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1},
}

var squareEven = [][]uint8{
	{1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1},
}

var paddedSquare = [][]uint8{
	{0, 0, 0, 0, 0},
	{0, 1, 1, 1, 0},
	{0, 1, 1, 1, 0},
	{0, 1, 1, 1, 0},
	{0, 0, 0, 0, 0},
}

var transparent = [][]uint8{
	{0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0},
}

var empty = [][]uint8{}

var emptyWidth = [][]uint8{
	{},
	{},
	{},
	{},
	{},
}

var onePixTransparent = [][]uint8{
	{0},
}

var onePixColored = [][]uint8{
	{1},
}

var oneRow = [][]uint8{
	{0, 0, 0, 1, 0, 1},
}

var oneCol = [][]uint8{
	{0},
	{0},
	{0},
	{1},
	{0},
	{1},
}

var topLeftCorner = [][]uint8{
	{1, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
}

var topRightCorner = [][]uint8{
	{0, 0, 1},
	{0, 0, 0},
	{0, 0, 0},
}

var bottomRightCorner = [][]uint8{
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 1},
}

var bottomLeftCorner = [][]uint8{
	{0, 0, 0},
	{0, 0, 0},
	{1, 0, 0},
}

func randomBitmap(width int, height int) [][]uint8 {
	out := make([][]uint8, height)

	for i := 0; i < height; i++ {
		out[i] = make([]uint8, width)
		for j := 0; j < width; j++ {
			out[i][j] = uint8(rand.Intn(2))
		}
	}

	return out
}

func testFrameFunc(t *testing.T, prefix string, frameFunc func(img *TolerantImage, radius int) (image.Point, image.Point)) {

	expect := func(name string, tlc image.Point, brc image.Point, expetedTlc image.Point, expectedBrc image.Point) {
		if tlc.X != expetedTlc.X || tlc.Y != expetedTlc.Y || brc.X != expectedBrc.X || brc.Y != expectedBrc.Y {
			t.Fatalf("Unexpected frame for %s: Got tlc: %v, brc: %v, Expected tlc: %v brc: %v", fmt.Sprint(prefix, " ", name), tlc, brc, expetedTlc, expectedBrc)
		}
	}

	testBitmap := func(name string, bitmap [][]uint8, expectedTlc image.Point, expectedBrc image.Point) {
		tlc, brc := frameFunc(&TolerantImage{
			newImageFromBitMap(bitmap),
			255,
		}, 0)
		expect(name, tlc, brc, expectedTlc, expectedBrc)
	}

	testBitmap("squareOdd", squareOdd, image.Point{0, 0}, image.Point{4, 4})
	testBitmap("squareEven", squareEven, image.Point{0, 0}, image.Point{5, 5})
	testBitmap("paddedSquare", paddedSquare, image.Point{1, 1}, image.Point{3, 3})
	testBitmap("transparent", transparent, image.Point{5, 5}, image.Point{-1, -1})
	testBitmap("empty", empty, image.Point{0, 0}, image.Point{-1, -1})
	testBitmap("emptyWidth", emptyWidth, image.Point{0, 5}, image.Point{-1, -1})
	testBitmap("onePixTransparent", onePixTransparent, image.Point{1, 1}, image.Point{-1, -1})
	testBitmap("onePixColored", onePixColored, image.Point{0, 0}, image.Point{0, 0})
	testBitmap("oneRow", oneRow, image.Point{3, 0}, image.Point{5, 0})
	testBitmap("oneCol", oneCol, image.Point{0, 3}, image.Point{0, 5})
	testBitmap("topLeftCorner", topLeftCorner, image.Point{0, 0}, image.Point{0, 0})
	testBitmap("topRightCorner", topRightCorner, image.Point{2, 0}, image.Point{2, 0})
	testBitmap("bottomRightCorner", bottomRightCorner, image.Point{2, 2}, image.Point{2, 2})
	testBitmap("bottomLeftCorner", bottomLeftCorner, image.Point{0, 2}, image.Point{0, 2})
}

func TestDeadSimpleFrame(t *testing.T) {
	testFrameFunc(t, "deadSimple", deadSimpleFrame)
}

func TestFrame(t *testing.T) {
	testFrameFunc(t, "frame", frame)
}

func BenchmarkFrame(t *testing.B) {

	img := &TolerantImage{
		newImageFromBitMap(randomBitmap(1920, 1080)),
		255,
	}

	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		frame(img, 0)
	}
}

func BenchmarkDeadSimpleFrame(t *testing.B) {

	img := &TolerantImage{
		newImageFromBitMap(randomBitmap(1920, 1080)),
		255,
	}

	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		deadSimpleFrame(img, 0)
	}
}
