package tree

import "bytes"

type word struct {
	value byte
	count int64
	code  *bytes.Buffer
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

func newInitialNode(p int64) *initialNode {
	return &initialNode{count: p}
}

func newAbstractNode(p int64) *abstractNode {
	return &abstractNode{count: p}
}

func (node *initialNode) getCount() int64 {
	return node.count
}

func (node *abstractNode) getCount() int64 {
	return node.count
}
