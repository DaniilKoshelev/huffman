package tree

import (
	"bytes"
	"container/list"
)

const maxWords = 255

type Tree struct {
	words [maxWords]*word
	nodes *list.List
}

func newTree() *Tree {
	tree := new(Tree)
	tree.nodes = list.New()

	return tree
}

// TODO: адаптировать под дерево из файла
func (tree *Tree) walkTree() {
	root := tree.nodes.Front().Value
	code := new(bytes.Buffer)

	rootAbstract, ok := tree.nodes.Front().Value.(*abstractNode)

	if !ok {
		code.WriteByte('0')
		root.(*initialNode).getWord().code = code
	} else {
		tree.walkFromNode(rootAbstract, code)
	}
}

// TODO: адаптировать под дерево из файла
func (tree *Tree) walkFromNode(node node, code *bytes.Buffer) {
	_, ok := node.(*abstractNode)

	if !ok {
		node.(*initialNode).getWord().code = code

		return
	}

	left := node.(*abstractNode).left
	right := node.(*abstractNode).right

	if left != nil {
		leftCode := new(bytes.Buffer)
		leftCode.Write(code.Bytes())
		leftCode.WriteByte('1')
		tree.walkFromNode(left, leftCode)
	}

	if right != nil {
		rightCode := new(bytes.Buffer)
		rightCode.Write(code.Bytes())
		rightCode.WriteByte('0')
		tree.walkFromNode(right, rightCode)
	}
}
