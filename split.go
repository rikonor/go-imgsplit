package imgsplit

import (
	"errors"
	"fmt"
	"image"
)

const (
	minWidth  = 10
	minHeight = 10
)

// Config defines how to split a image
type Config struct {
	// X and Y count how many images the parent image will be split to
	X, Y int
}

// SplitImageWithIterator will split a given Image based on the given Config
// and return an iterator of the sub images
// This is more memory efficient then returning all images at once
func SplitImageWithIterator(img image.Image, cfg Config) (ImageIterator, error) {
	// Only subImagers can be split
	simgr, ok := img.(subImager)
	if !ok {
		return nil, errors.New("only images implementing SubImage are supported")
	}

	size := simgr.Bounds().Size()

	// approximage how large each image should be
	subImageWidth := size.X / cfg.X
	subImageHeight := size.Y / cfg.Y

	// Check sub-image size
	if subImageWidth < minWidth || subImageHeight < minHeight {
		return nil, fmt.Errorf("sub images size is too small: (%d, %d)", subImageWidth, subImageHeight)
	}

	row, col := 0, -1

	return &mockImageIteratoer{
		NextFn: func() bool {
			col++
			if col == cfg.X {
				col = 0
				row++
			}

			// Stop iterating when row and col have reached their maximum values
			if row == cfg.Y {
				return false
			}
			return true
		},
		GetFn: func() image.Image {
			// Define the position and size of the next subimage
			r := image.Rect(
				col*subImageWidth, row*subImageHeight, // (x0, y0)
				(col+1)*subImageWidth, (row+1)*subImageHeight, // (x1, y1)
			)
			return simgr.SubImage(r)
		},
	}, nil
}
