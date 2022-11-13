package core

import (
	"bufio"
	"huffman/src/core/bitsbuffer"
	"huffman/src/core/tree"
	"io"
)

type Decoder struct {
	tree           *tree.Tree
	uniqueWords    uint8
	bitsInLastByte uint8
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (decoder *Decoder) Init(fileBuffer *bitsbuffer.Buffer) error {
	uniqueWords, err := fileBuffer.ReadByte()

	if err != nil {
		return err
	}

	decoder.uniqueWords = uniqueWords

	bitsInLastByte, err := fileBuffer.ReadByte()

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

func (decoder *Decoder) Decode(inputFileBuffer *bitsbuffer.Buffer, writer *bufio.Writer) error {
	outputFileBuffer := bitsbuffer.NewEmptyFlushableBuffer(writer)

	currentCode := bitsbuffer.NewEmptyBuffer()

	err := inputFileBuffer.Scan()

	max := (inputFileBuffer.Length()/8)*8 - (8 - int8(decoder.bitsInLastByte))
	// TODO: do while
	if err == io.EOF {
		var i int8 = 0
		for ; i < max && !inputFileBuffer.IsEmpty(); i++ {
			decoder.processNextBit(inputFileBuffer, outputFileBuffer, currentCode)
		}

		outputFileBuffer.Flush()
		return nil
	}

	var i int8 = 0
	for ; i < max && !inputFileBuffer.IsEmpty(); i++ {
		decoder.processNextBit(inputFileBuffer, outputFileBuffer, currentCode)
	}

	for {
		err := inputFileBuffer.Scan()

		if err == io.EOF {
			var i int8 = 0
			for ; i < max && !inputFileBuffer.IsEmpty(); i++ {
				decoder.processNextBit(inputFileBuffer, outputFileBuffer, currentCode)
			}

			break
		}

		var i int8 = 0
		for ; i < max && !inputFileBuffer.IsEmpty(); i++ {
			decoder.processNextBit(inputFileBuffer, outputFileBuffer, currentCode)
		}
	}

	outputFileBuffer.Flush()

	return nil
}

func (decoder *Decoder) processNextBit(inputFileBuffer *bitsbuffer.Buffer, outputFileBuffer *bitsbuffer.Buffer, currentCode *bitsbuffer.Buffer) {
	bit, _ := inputFileBuffer.ReadBit()

	currentCode.AddBit(bit)

	if word, err := decoder.tree.GetWordByte(currentCode); err == nil {
		outputFileBuffer.AddByte(word)
		currentCode.Reset()
	}
}
