package tree

import (
	"bytes"
	"container/list"
	"huffman/src/core/bitsbuffer"
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

func (tree *Tree) buildTree() {
	root := tree.nodes.Front().Value
	code := new(bytes.Buffer)

	rootAbstract, ok := tree.nodes.Front().Value.(*abstractNode)

	if !ok {
		code.WriteByte('0')
		root.(*initialNode).getWord().code = code
	} else {
		tree.buildFromNode(rootAbstract, code)
	}
}

func (tree *Tree) buildFromNode(node node, code *bytes.Buffer) {
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
		tree.buildFromNode(left, leftCode)
	}

	if right != nil {
		rightCode := new(bytes.Buffer)
		rightCode.Write(code.Bytes())
		rightCode.WriteByte('0')
		tree.buildFromNode(right, rightCode)
	}
}

func (tree *Tree) Pack() *bytes.Buffer {
	root := tree.nodes.Front().Value
	packedTree := new(bytes.Buffer)
	rootAbstract, ok := tree.nodes.Front().Value.(*abstractNode)
	buffer := bitsbuffer.NewBuffer(0, 0, packedTree)

	if !ok {
		buffer.AddOne()
		buffer.AddByte(root.(*initialNode).getWord().value)
	} else {
		tree.packFromNode(rootAbstract, buffer)
	}

	buffer.Flush()
	// TODO: последний байт используется не полностью

	return packedTree
}

func (tree *Tree) packFromNode(node node, buffer *bitsbuffer.Buffer) {
	_, ok := node.(*abstractNode)

	if !ok {
		buffer.AddOne()
		buffer.AddByte(node.(*initialNode).getWord().value)

		return
	}

	buffer.AddZero()

	left := node.(*abstractNode).left
	right := node.(*abstractNode).right

	if left != nil {
		tree.packFromNode(left, buffer)
	}

	if right != nil {
		tree.packFromNode(right, buffer)
	}
}
