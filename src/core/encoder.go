package core

import (
	"bufio"
	"huffman/src/core/tree"
)

type Encoder struct {
	tree *tree.Tree
}

func NewEncoder(reader *bufio.Reader) (*Encoder, error) {
	encoder := new(Encoder)
	createdTree, err := tree.Create(reader)

	if err != nil {
		return nil, err
	}

	encoder.tree = createdTree

	return encoder, nil
}

func (encoder *Encoder) Encode(ch chan []byte) {

}
