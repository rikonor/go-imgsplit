package imgsplit

import (
	"image"
	"testing"
)

func TestConsumeIterator(t *testing.T) {
	// Mock an iterator that returns 3 images
	itms := []image.Image{
		// ConsumeIterator doesn't care if the images are nil
		nil, nil, nil,
	}
	it := &mockImageIteratoer{
		NextFn: func() bool {
			return len(itms) > 0
		},
		GetFn: func() image.Image {
			m := itms[0]
			itms = itms[1:]
			return m
		},
	}

	ms := ConsumeIterator(it)
	expectedCount := 3

	if len(ms) != expectedCount {
		t.Fatalf("expected %d images, got %d", expectedCount, len(ms))
	}
}
