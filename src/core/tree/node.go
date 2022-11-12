package tree

import (
	"huffman/src/core/bitsbuffer"
)

type word struct {
	value byte
	count int64
	code  *bitsbuffer.Buffer
}

type node interface {
	getCount() int64
}

type abstractNode struct {
	count int64 // Количество раз, сколько встретился определенный байт
	left  node
	right node
}

type initialNode struct {
	count int64
	word  *word
}

func (node *initialNode) getWord() *word {
	return node.word
}

func (node *initialNode) setWord(word *word) {
	node.word = word
}

func newInitialNode() *initialNode {
	return &initialNode{}
}

func newAbstractNode(left node, right node) *abstractNode {
	return &abstractNode{left: left, right: right}
}

func newInitialNodeCount(p int64) *initialNode {
	return &initialNode{count: p}
}

func newAbstractNodeCount(p int64) *abstractNode {
	return &abstractNode{count: p}
}

func (node *initialNode) getCount() int64 {
	return node.count
}

func (node *abstractNode) getCount() int64 {
	return node.count
}

func (word *word) Length() int8 {
	return word.code.Length()
}

func (word *word) Count() int64 {
	return word.count
}
