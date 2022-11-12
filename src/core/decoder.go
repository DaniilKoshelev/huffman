package core

import (
	"bufio"
	"huffman/src/core/bitsbuffer"
	"huffman/src/core/tree"
)

type Decoder struct {
	tree           *tree.Tree
	uniqueWords    uint8
	bitsInLastByte uint8
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (decoder *Decoder) Init(reader *bufio.Reader) error {
	fileBuffer := bitsbuffer.NewEmptyBuffer().SetIoReader(reader)

	uniqueWords, err := reader.ReadByte()

	if err != nil {
		return err
	}

	decoder.uniqueWords = uniqueWords

	bitsInLastByte, err := reader.ReadByte()

	if err != nil {
		return err
	}

	decoder.bitsInLastByte = bitsInLastByte

	createdTree, err := tree.Restore(fileBuffer, uniqueWords)

	if err != nil {
		return err
	}

	decoder.tree = createdTree

	return nil
}

func (decoder *Decoder) Decode(reader *bufio.Reader, writer *bufio.Writer) error {
	//TODO: implement
	return nil
}
