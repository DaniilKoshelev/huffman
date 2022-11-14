package core

import (
	"bufio"
	"encoding/binary"
	"huffman/src/core/bitsbuffer"
	"huffman/src/core/tree"
	"io"
)

type Decoder struct {
	tree           *tree.Tree
	uniqueWords    uint32
	bitsInLastByte uint8
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (decoder *Decoder) Init(fileBuffer *bitsbuffer.Buffer) error {
	byte1, err := fileBuffer.ReadByte()
	byte2, err := fileBuffer.ReadByte()
	byte3, err := fileBuffer.ReadByte()
	byte4, err := fileBuffer.ReadByte()

	if err != nil {
		return err
	}

	var uniqueWords = binary.LittleEndian.Uint32([]byte{byte1, byte2, byte3, byte4})

	decoder.uniqueWords = uniqueWords

	bitsInLastByte, err := fileBuffer.ReadByte()

	if err != nil {
		return err
	}

	hasPadding, err := fileBuffer.ReadByte()

	if err != nil {
		return err
	}

	if bitsInLastByte == 0 {
		bitsInLastByte = 8
	}

	decoder.bitsInLastByte = bitsInLastByte

	createdTree, err := tree.Restore(fileBuffer, uniqueWords)

	if err != nil {
		return err
	}

	if hasPadding == 1 {
		createdTree.SetHasPadding()
	}

	decoder.tree = createdTree

	return nil
}

func (decoder *Decoder) Decode(inputFileBuffer *bitsbuffer.Buffer, writer *bufio.Writer) error {
	outputFileBuffer := bitsbuffer.NewEmptyFlushableBuffer(writer)

	currentCode := bitsbuffer.NewEmptyBuffer()

	for inputFileBuffer.Length() > 8 {
		decoder.processNextBit(inputFileBuffer, outputFileBuffer, currentCode)
	}

	for {
		err := inputFileBuffer.Scan()

		if err == io.EOF {
			var i int8 = 0

			max := inputFileBuffer.Length() - (8 - int8(decoder.bitsInLastByte))
			for ; i < max && !inputFileBuffer.IsEmpty(); i++ {
				decoder.processNextBit(inputFileBuffer, outputFileBuffer, currentCode)
			}

			break
		}

		err = inputFileBuffer.Scan()

		if err == io.EOF {
			var i int8 = 0

			max := inputFileBuffer.Length() - (8 - int8(decoder.bitsInLastByte))
			for ; i < max && !inputFileBuffer.IsEmpty(); i++ {
				decoder.processNextBit(inputFileBuffer, outputFileBuffer, currentCode)
			}

			break
		}

		var i int8 = 0
		for ; i < 8; i++ {
			decoder.processNextBit(inputFileBuffer, outputFileBuffer, currentCode)
		}
	}

	outputFileBuffer.Flush()

	return nil
}

func (decoder *Decoder) processNextBit(inputFileBuffer *bitsbuffer.Buffer, outputFileBuffer *bitsbuffer.Buffer, currentCode *bitsbuffer.Buffer) {
	bit, _ := inputFileBuffer.ReadBit()

	currentCode.AddBit(bit)

	if word, err := decoder.tree.GetWordBytes(currentCode); err == nil {
		outputFileBuffer.AddUInt16(word)
		currentCode.Reset()
	}
}
