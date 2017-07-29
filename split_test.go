package imgsplit

import (
	"image"
	"image/jpeg"
	"os"
	"testing"
)

func getTestImage() (image.Image, error) {
	f, err := os.Open("./resources/test.jpg")
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func TestSplitImageWithIterator(t *testing.T) {
	in, err := getTestImage()
	if err != nil {
		t.Fatalf("failed to load test image: %s", err)
	}

	// set cfg
	cfg := Config{X: 3, Y: 3}

	// split
	it, err := SplitImageWithIterator(in, cfg)
	if err != nil {
		t.Fatalf("unexpected failure: %s", err)
	}

	// drain the parts to check expected count
	outs := []image.Image{}
	for it.Next() {
		outs = append(outs, it.Get())
	}

	if len(outs) != 9 {
		t.Fatalf("wrong number of parts: %d", len(outs))
	}
}
