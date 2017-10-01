package quadtreego

import (
	"testing"
	"math/rand"
	"fmt"
	"image"
	"time"
)

func assertEqual(t *testing.T, expected interface{}, actual interface{}, message string) {
	if expected == actual {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", expected, actual)
	}
	t.Fatal(message)
}

var img = image.NewRGBA(image.Rect(0, 0, 101, 101))

func drawBBOX(x float64, y float64, w float64, h float64, pt Point) {

	drawRect(int(x), int(y), int(x)+int(w), int(y)+int(h), img)
}

func drawPT(x float64, y float64, w float64, h float64, pt Point) {

	drawPoint(int(pt.X), int(pt.Y), img)
}

func getTree() *Quadtree {
	tr := NewQuadTree(0, 0, 100, 100)
	tr.set(5, 20, "Foo")
	tr.set(50, 32, "Bar")
	tr.set(47, 96, "Baz")
	tr.set(50, 50, "Bing")
	tr.set(12, 0, "Bong")
	return tr
}

func getRandTree() *Quadtree {
	tr := NewQuadTree(0, 0, 100, 100)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 100; i++ {

		tr.set(float64(rand.Intn(100)), float64(rand.Intn(100)), float64(i))
	}
	return tr
}

func TestCreateQuadTree(t *testing.T) {

	qt := getTree()
	assertEqual(t, 5, qt.Count, "Count must be 5")
}

func TestGet(t *testing.T) {

	qt := getTree()
	val := qt.get(47, 96)
	assertEqual(t, "Baz", val, "expected val not returned")

	val2 := qt.get(47, 91)
	assertEqual(t, nil, val2, "no value returned for wrong point")

}

func TestSearchQuadTree(t *testing.T) {

	qt := getTree()
	result := qt.find(qt.Root, 12, 0)
	assertEqual(t, 12.0, result.Point.X, "point X much match")
	assertEqual(t, 0.0, result.Point.Y, "point Y much match")

	result2 := qt.find(qt.Root, 12, 10)

	assertEqual(t, nil, result2.Point.Weight, "expected nil  ")

}

func TestTraverseTree(t *testing.T) {

	qt := getRandTree()

	var f traverseFunc
	f = drawBBOX
	qt.traverse(qt.Root, f)

	var f2 traverseFunc
	f2 = drawPT
	qt.traverse(qt.Root, f2)

	render(img, "QT.png")

	qt.balance(qt.Root)

}
