package quadtreego

import (
	"fmt"
	"errors"
)

type Point struct {
	X      float64
	Y      float64
	Weight interface{}
}

type Node struct {
	Left     float64
	Top      float64
	Width    float64
	Height   float64
	Parent   *Node
	NW       *Node
	NE       *Node
	SW       *Node
	SE       *Node
	Point    Point
	NodeType NodeType
}

type NodeType int

const (
	EMPTY   NodeType = iota
	LEAF
	POINTER
)

type Quadtree struct {
	Root  Node
	Count int
}

func NewQuadTree(minX float64, minY float64, maxX float64, maxY float64) *Quadtree {
	root := Node{Left: minX, Top: minY, Width: maxX - minX, Height: maxY - minY}

	return &Quadtree{
		root, 0,
	}
}

func (tree *Quadtree) find(node Node, x float64, y float64) Node {
	var response Node
	switch node.NodeType {
	case EMPTY:
		{
			break
		}
	case LEAF:
		{
			if node.Point.X == x && node.Point.Y == y {
				response = node
			}
		}
	case POINTER:
		{
			response = tree.find(*getQuadrantForPoint(&node, x, y), x, y);
		}
	}

	return response
}

func (tree *Quadtree) set(x float64, y float64, value interface{}) error {
	if x < tree.Root.Left || y < tree.Root.Top || x > tree.Root.Left+tree.Root.Width || y > tree.Root.Top+tree.Root.Height {
		return errors.New(fmt.Sprintf("Out of bounds :( %d, %d )", x, y))
	}
	if insert(&tree.Root, Point{x, y, value}) {
		tree.Count++;
	}

	return nil
}

func (tree *Quadtree) get(x float64, y float64) interface{} {
	node := tree.find(tree.Root, x, y)
	return node.Point.Weight
}

func (tree *Quadtree) remove(x float64, y float64) {
	node := tree.find(tree.Root, x, y)
	if node.NodeType != EMPTY {
		node.NodeType = EMPTY
		tree.Count--
		tree.balance(node)
	}
}

func (tree *Quadtree) clear() {
	tree.Root.NW = nil
	tree.Root.NE = nil
	tree.Root.SW = nil
	tree.Root.SE = nil
	tree.Root.NodeType = EMPTY
	tree.Root.Point = Point{}
	tree.Count = 0

}

func (tree *Quadtree) traverse(node Node, TF traverseFunc) {

	switch node.NodeType {

	case LEAF:
		{
			TF(node.Left, node.Top, node.Width, node.Height, node.Point)
		}
	case POINTER:
		{
			tree.traverse(*node.NE, TF)
			tree.traverse(*node.SE, TF)
			tree.traverse(*node.SW, TF)
			tree.traverse(*node.NW, TF)
		}
	case EMPTY:
		{
			//TF(node.Left, node.Top, node.Width, node.Height, node.Point)

		}
	}

}

func (tree *Quadtree) balance(node Node) {
	switch node.NodeType {
	case EMPTY:
	case LEAF:
		{
			if node.Parent != nil {
				tree.balance(*node.Parent)
			}
		}
	case POINTER:
		{
			nw := node.NW
			ne := node.NE
			sw := node.SW
			se := node.SE

			var firstLeaf *Node
			if nw.NodeType != EMPTY {

				firstLeaf = nw
			}
			if ne.NodeType != EMPTY {
				if firstLeaf != nil {
					break
				}
				firstLeaf = ne
			}

			if sw.NodeType != EMPTY {
				if firstLeaf != nil {
					break
				}
				firstLeaf = sw
			}

			if se.NodeType != EMPTY {
				if firstLeaf != nil {
					break
				}
				firstLeaf = se
			}

			if firstLeaf == nil {
				node.NodeType = EMPTY
				node.NW = nil
				node.NE = nil
				node.SW = nil
				node.SE = nil
			} else if firstLeaf.NodeType == POINTER {
				break
			} else {
				node.NodeType = LEAF
				node.NW = nil
				node.NE = nil
				node.SW = nil
				node.SE = nil
				node.Point = firstLeaf.Point

			}

			if node.Parent != nil {
				tree.balance(*node.Parent)
			}

		}
	}
}


/*============================ PRIVATE FUNCS ============================ */

func insert(parent *Node, point Point) bool {
	var result = false

	switch parent.NodeType {
	case EMPTY:
		setPointForNode(parent, point)
		result = true
	case LEAF:
		if parent.Point.X == point.X && parent.Point.Y == point.Y {
			setPointForNode(parent, point)
			result = true
		} else {
			split(parent)
			result = insert(parent, point)
		}
		result = true

	case POINTER:
		result = insert(getQuadrantForPoint(parent, point.X, point.Y), point)
	}

	return result
}

type traverseFunc func(float64, float64, float64, float64, Point)

func getQuadrantForPoint(parent *Node, x float64, y float64) *Node {
	mx := parent.Left + parent.Width/2
	my := parent.Top + parent.Height/2
	if x < mx {
		if y < my {
			return parent.NW
		} else {
			return parent.SW
		}
	} else {
		if y < my {
			return parent.NE
		} else {
			return parent.SE
		}
	}
}

func split(node *Node) {
	oldPoint := node.Point
	node.Point = Point{}
	node.NodeType = POINTER

	x := node.Left
	y := node.Top
	subWidth := node.Width / 2
	subHeight := node.Height / 2

	node.NW = &Node{Left: x, Top: y, Width: subWidth, Height: subHeight, Parent: node}
	node.NE = &Node{Left: x + subWidth, Top: y, Width: subWidth, Height: subHeight, Parent: node}
	node.SW = &Node{Left: x, Top: y + subHeight, Width: subWidth, Height: subHeight, Parent: node}
	node.SE = &Node{Left: x + subWidth, Top: y + subHeight, Width: subWidth, Height: subHeight, Parent: node}

	insert(node, oldPoint)

}

func setPointForNode(node *Node, point Point) {
	if node.NodeType == POINTER {
		panic("Can not set point for node of type POINTER");
	}
	node.NodeType = LEAF
	node.Point = point

}
