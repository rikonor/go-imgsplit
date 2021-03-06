package imgsplit

import (
	"image"
	"image/color"
	"testing"

	imgutil "github.com/rikonor/go-imgutil"
)

func TestQuadTreeTrivialCase(t *testing.T) {
	// The trivial case is for an image that shouldn't cause any partitioning
	// e.g a uniform color image
	r := image.Rect(0, 0, 100, 100)
	m := image.NewNRGBA64(r)

	c := color.White
	imgutil.Iterate(m, func(x int, y int) {
		m.Set(x, y, c)
	})

	it, err := QuadTreeIterator(m, 4, 2000.0)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ms := ConsumeIterator(it)
	expectedCount := 1

	if len(ms) != expectedCount {
		t.Fatalf("expected %d images but got %d", expectedCount, len(ms))
	}
}

func TestQuadTreeThresholdNotReached(t *testing.T) {
	// The threshold not reached case is for an image where the sub images
	// are similar enough to not merit partitioning them further
	r := image.Rect(0, 0, 100, 100)
	m := image.NewNRGBA64(r)

	// Define four colors
	baseR, baseG := 32767, 32767
	delta := 1000
	c1 := color.NRGBA64{R: uint16(baseR - delta), G: uint16(baseG - delta), A: 65535}
	c2 := color.NRGBA64{R: uint16(baseR - delta), G: uint16(baseG + delta), A: 65535}
	c3 := color.NRGBA64{R: uint16(baseR + delta), G: uint16(baseG - delta), A: 65535}
	c4 := color.NRGBA64{R: uint16(baseR + delta), G: uint16(baseG + delta), A: 65535}

	imgutil.Iterate(m, func(x int, y int) {
		switch {
		// top left quadrant
		case x >= 0 && x < 50 && y >= 0 && y < 50:
			m.Set(x, y, c1)

		// top right quadrant
		case x >= 50 && x < 100 && y >= 0 && y < 50:
			m.Set(x, y, c2)

		// bottom left quadrant
		case x >= 0 && x < 50 && y >= 50 && y < 100:
			m.Set(x, y, c3)

		// bottom right quadrant
		case x >= 50 && x < 100 && y >= 50 && y < 100:
			m.Set(x, y, c4)
		}
	})

	it, err := QuadTreeIterator(m, 4, 2000.0)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ms := ConsumeIterator(it)
	expectedCount := 1

	if len(ms) != expectedCount {
		t.Fatalf("expected %d images but got %d", expectedCount, len(ms))
	}
}

func TestQuadTreeOneLevelIn(t *testing.T) {
	// the one level case is for an image that is composed of 4 uniformly colored squares
	// the result of this iterator should be 4 sub-images, one for each square
	r := image.Rect(0, 0, 100, 100)
	m := image.NewNRGBA64(r)

	// Define four colors
	c1 := color.White
	c2 := color.Black
	c3 := color.NRGBA64{R: 65535, A: 65535} // Red
	c4 := color.NRGBA64{B: 65535, A: 65535} // Blue

	imgutil.Iterate(m, func(x int, y int) {
		switch {
		// top left quadrant
		case x >= 0 && x < 50 && y >= 0 && y < 50:
			m.Set(x, y, c1)

		// top right quadrant
		case x >= 50 && x < 100 && y >= 0 && y < 50:
			m.Set(x, y, c2)

		// bottom left quadrant
		case x >= 0 && x < 50 && y >= 50 && y < 100:
			m.Set(x, y, c3)

		// bottom right quadrant
		case x >= 50 && x < 100 && y >= 50 && y < 100:
			m.Set(x, y, c4)
		}
	})

	it, err := QuadTreeIterator(m, 4, 2000.0)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ms := ConsumeIterator(it)
	expectedCount := 4

	if len(ms) != expectedCount {
		t.Fatalf("expected %d images but got %d", expectedCount, len(ms))
	}
}

func TestQuadTreeMaxDepthReached(t *testing.T) {
	// the max depth reached case is for hitting the max depth and not triggering
	// any more tree levels, the result should be just one level (in this case)
	r := image.Rect(0, 0, 100, 100)
	m := image.NewNRGBA64(r)

	// Define four colors
	// this image would usually trigger 2 levels
	c1 := color.White
	c2 := color.Black
	c3 := color.NRGBA64{R: 65535, A: 65535} // Red
	c4 := color.NRGBA64{B: 65535, A: 65535} // Blue

	imgutil.Iterate(m, func(x int, y int) {
		switch {
		// top left quadrant
		case x >= 0 && x < 50 && y >= 0 && y < 50:
			m.Set(x, y, c1)

		// top right quadrant
		case x >= 50 && x < 100 && y >= 0 && y < 50:
			m.Set(x, y, c2)

		// bottom left quadrant
		case x >= 0 && x < 50 && y >= 50 && y < 100:
			m.Set(x, y, c3)

		// bottom right quadrant
		case x >= 50 && x < 100 && y >= 50 && y < 100:
			m.Set(x, y, c4)
		}
	})

	it, err := QuadTreeIterator(m, 1, 2000.0)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ms := ConsumeIterator(it)
	expectedCount := 1

	if len(ms) != expectedCount {
		t.Fatalf("expected %d images but got %d", expectedCount, len(ms))
	}
}
