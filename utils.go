package imgsplit

import "image"

// ConsumeIterator drains the iterator of images and returns them in a slice
// Note that consuming an entire iterator may cause heavy memory usage
// and usually is a bad idea
func ConsumeIterator(it ImageIterator) []image.Image {
	ms := []image.Image{}
	for it.Next() {
		ms = append(ms, it.Get())
	}
	return ms
}
