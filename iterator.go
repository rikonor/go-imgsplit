package imgsplit

import "image"

// subImager has a SubImage method
type subImager interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

// ImageIterator is an iterators on a set of images
type ImageIterator interface {
	Next() bool
	Get() image.Image
}

type mockImageIteratoer struct {
	NextFn func() bool
	GetFn  func() image.Image
}

func (it *mockImageIteratoer) Next() bool {
	return it.NextFn()
}

func (it *mockImageIteratoer) Get() image.Image {
	return it.GetFn()
}
