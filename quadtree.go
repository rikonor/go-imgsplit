package imgsplit

import (
	"image"
	"math"

	"github.com/gonum/stat"
	imgutil "github.com/rikonor/go-imgutil"
)

// QuadTreeIterator iterates over an image using a quadtree mesh
// m 								-> image to iterate over
// maxDepth 				-> max depth of quadtree
// maxDissimilarity -> max amount of dissimilarity allowed, anything above will trigger another tree level
func QuadTreeIterator(m image.Image, maxDepth int, maxDissimilarity float64) (ImageIterator, error) {
	root := buildSubtree(m, maxDepth, 1, maxDissimilarity)
	images := flattenTree(root)
	return mockIteratorFromImages(images), nil
}

type node struct {
	children []*node

	// for leaf
	m image.Image
}

func buildSubtree(m image.Image, maxDepth, currentLevel int, maxDissimilarity float64) *node {
	// Stop if max level reached
	if currentLevel == maxDepth {
		return &node{m: m}
	}

	// Divide the given image into a set of subimages
	subimgs := partitionImage(m, 4)

	// Calculate the average red/green/blue for every subimage
	avgRs, avgGs, avgBs := []float64{}, []float64{}, []float64{}
	for _, subimg := range subimgs {
		avgR, avgG, avgB := avgRGB(subimg)
		avgRs = append(avgRs, avgR)
		avgGs = append(avgGs, avgG)
		avgBs = append(avgBs, avgB)
	}

	// Calculate the standard deviation of the avgs
	stdDevR := stat.StdDev(avgRs, nil)
	stdDevG := stat.StdDev(avgGs, nil)
	stdDevB := stat.StdDev(avgBs, nil)

	// Check if subimages similarity is within given threshold
	simDiff := math.Sqrt(stdDevR*stdDevR + stdDevG*stdDevG + stdDevB*stdDevB)
	if simDiff <= maxDissimilarity {
		// If similarity is within threshold, that means all the subimages
		// are fairly similar and there's no need to divide them further

		// return a node with no children
		return &node{m: m}
	}

	// divide the image further
	children := []*node{}
	for _, subimg := range subimgs {
		child := buildSubtree(subimg, maxDepth, currentLevel+1, maxDissimilarity)
		children = append(children, child)
	}

	return &node{children: children}
}

func partitionImage(m image.Image, childrenCount int) []image.Image {
	// Children count has to be a square of an integer (4, 9, 25)
	sq := math.Sqrt(float64(childrenCount))
	if float64(int(sq)) != sq {
		panic("images can only be partitioned into N where N is a square of an integer")
	}
	sqi := int(sq)

	dx := m.Bounds().Dx()
	dy := m.Bounds().Dy()

	subImageWidth := int(float64(dx) / sq)
	subImageHeight := int(float64(dy) / sq)

	// Only subImagers can be split
	simgr, ok := m.(subImager)
	if !ok {
		panic("only images implementing SubImage are supported")
	}

	subimgs := []image.Image{}
	for row := 0; row < sqi; row++ {
		for col := 0; col < sqi; col++ {
			// Get the coords for the subimage
			x0, y0 := row*subImageWidth, col*subImageHeight
			x1, y1 := (row+1)*subImageWidth, (col+1)*subImageHeight

			// Offset to account for location in original image
			pMin := simgr.Bounds().Min
			x0, y0 = x0+pMin.X, y0+pMin.Y
			x1, y1 = x1+pMin.X, y1+pMin.Y

			r := image.Rect(
				x0, y0,
				x1, y1,
			)
			subimg := simgr.SubImage(r)
			subimgs = append(subimgs, subimg)
		}
	}

	return subimgs
}

func avgRGB(m image.Image) (avgR, avgG, avgB float64) {
	avgR, avgG, avgB = .0, .0, .0

	imgutil.Iterate(m, func(x int, y int) {
		r, g, b, _ := m.At(x, y).RGBA()

		avgR += float64(r)
		avgG += float64(g)
		avgB += float64(b)
	})

	dx := m.Bounds().Dx()
	dy := m.Bounds().Dy()
	pxCount := float64(dx * dy)

	avgR /= pxCount
	avgG /= pxCount
	avgB /= pxCount

	return avgR, avgG, avgB
}

func flattenTree(n *node) []image.Image {
	if len(n.children) == 0 {
		return []image.Image{n.m}
	}

	images := []image.Image{}
	for _, child := range n.children {
		images = append(images, flattenTree(child)...)
	}

	return images
}
