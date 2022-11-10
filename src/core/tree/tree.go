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
	code := bitsbuffer.NewBuffer()

	rootAbstract, ok := tree.nodes.Front().Value.(*abstractNode)

	if !ok {
		code.AddZero()
		root.(*initialNode).getWord().code = code
	} else {
		tree.buildFromNode(rootAbstract, code)
	}
}

func (tree *Tree) buildFromNode(node node, code *bitsbuffer.Buffer) {
	_, ok := node.(*abstractNode)

	if !ok {
		node.(*initialNode).getWord().code = code

		return
	}

	left := node.(*abstractNode).left
	right := node.(*abstractNode).right

	if left != nil {
		leftCode := bitsbuffer.From(code)
		leftCode.AddOne()
		tree.buildFromNode(left, leftCode)
	}

	if right != nil {
		rightCode := bitsbuffer.From(code)
		rightCode.AddZero()
		tree.buildFromNode(right, rightCode)
	}
}

func (tree *Tree) Pack() (*bytes.Buffer, *bitsbuffer.Buffer) {
	root := tree.nodes.Front().Value
	packedTree := new(bytes.Buffer)
	rootAbstract, ok := tree.nodes.Front().Value.(*abstractNode)
	buffer := bitsbuffer.NewFlushableBuffer(packedTree)

	if !ok {
		buffer.AddOne()
		buffer.AddByte(root.(*initialNode).getWord().value)
	} else {
		tree.packFromNode(rootAbstract, buffer)
	}

	return packedTree, buffer
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

func (tree *Tree) GetCode(byteForWord byte) *bitsbuffer.Buffer {
	return tree.words[byteForWord].code
}
