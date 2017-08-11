package imgsplit

import "image"

// QuadTreeIterator ...
func QuadTreeIterator(img image.Image) (ImageIterator, error) {
	return &mockImageIteratoer{
		NextFn: func() bool {
			return false
		},
		GetFn: func() image.Image {
			return nil
		},
	}, nil
}
