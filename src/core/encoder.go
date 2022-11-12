package core

import (
	"bufio"
	"errors"
	"huffman/src/core/bitsbuffer"
	"huffman/src/core/tree"
	"io"
)

type Encoder struct {
	tree *tree.Tree
}

func NewEncoder() *Encoder {
	return new(Encoder)
}

func (encoder *Encoder) Init(reader *bufio.Reader) error {
	createdTree, err := tree.Create(reader)

	if err != nil {
		return err
	}

	encoder.tree = createdTree

	return nil
}

func (encoder *Encoder) Encode(reader *bufio.Reader, writer *bufio.Writer) error {
	if reader == nil {
		return errors.New("reader is not set")
	}

	packedTree, remainingBuffer := encoder.tree.Pack()

	_, err := writer.Write(packedTree.Bytes())
	fileBuffer := bitsbuffer.NewFlushableBuffer(writer)
	fileBuffer.AddFromBuffer(remainingBuffer)

	if err != nil {
		return err
	}

	//err = writer.Flush()

	if err != nil {
		return err
	}

	for {
		newByte, err := reader.ReadByte()

		if err == io.EOF {
			break
		}

		code := encoder.tree.GetCode(newByte)
		fileBuffer.AddFromBuffer(code)
	}

	fileBuffer.Flush()

	return nil
}
