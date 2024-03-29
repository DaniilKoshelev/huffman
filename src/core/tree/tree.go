package tree

import (
	"bytes"
	"container/list"
	"errors"
	"huffman/src/core/bitsbuffer"
)

type Tree struct {
	Words            map[uint16]*word
	Codes            map[string]*word
	nodes            *list.List
	nodesCount       uint32
	alreadyReadNodes uint32
	HasPadding       uint8 // TODO: убрать из дерева
}

func newTree() *Tree {
	tree := new(Tree)
	tree.nodes = list.New()
	tree.Codes = make(map[string]*word)
	tree.Words = make(map[uint16]*word)

	return tree
}

func (tree *Tree) buildTree() {
	root := tree.nodes.Front().Value
	code := bitsbuffer.NewEmptyBuffer()

	rootAbstract, ok := tree.nodes.Front().Value.(*abstractNode)

	if !ok {
		code.AddZero()
		word := root.(*initialNode).getWord()
		word.code = code
		tree.Codes[code.ToString()] = word
	} else {
		tree.buildFromNode(rootAbstract, code)
	}
}

func (tree *Tree) buildFromNode(node node, code *bitsbuffer.Buffer) {
	_, ok := node.(*abstractNode)

	if !ok {
		word := node.(*initialNode).getWord()
		word.code = code
		tree.Codes[code.ToString()] = word

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
	buffer := bitsbuffer.NewEmptyFlushableBuffer(packedTree)

	if !ok {
		buffer.AddOne()
		buffer.AddUInt16(root.(*initialNode).getWord().value)

	} else {
		tree.packFromNode(rootAbstract, buffer)
	}

	return packedTree, buffer
}

func (tree *Tree) packFromNode(node node, buffer *bitsbuffer.Buffer) {
	_, ok := node.(*abstractNode)

	if !ok {
		buffer.AddOne()
		buffer.AddUInt16(node.(*initialNode).getWord().value)

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

func (tree *Tree) GetCode(byteForWord uint16) *bitsbuffer.Buffer {
	return tree.Words[byteForWord].code
}

func (tree *Tree) GetWordBytes(codeForWord *bitsbuffer.Buffer) (uint16, error) {
	if word, ok := tree.Codes[codeForWord.ToString()]; ok {
		return word.value, nil
	}

	return 0, errors.New("error: no word found for given code")
}

func (tree *Tree) SetHasPadding() *Tree {
	tree.HasPadding = 1

	return tree
}
