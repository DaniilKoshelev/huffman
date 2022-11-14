package core

import (
	"bufio"
	"encoding/binary"
	"errors"
	"huffman/src/core/bitsbuffer"
	"huffman/src/core/tree"
	"io"
)

type Encoder struct {
	tree           *tree.Tree
	bitsInLastByte uint8
	uniqueWords    uint32
	hasPadding     uint8
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

	encoder.calculateMetaParams()

	uniqueWords := make([]byte, 4)
	binary.LittleEndian.PutUint32(uniqueWords, encoder.uniqueWords)

	_, err := writer.Write(uniqueWords)

	if err != nil {
		return err
	}

	err = writer.WriteByte(encoder.bitsInLastByte)

	if err != nil {
		return err
	}

	err = writer.WriteByte(encoder.tree.HasPadding)

	if err != nil {
		return err
	}

	_, err = writer.Write(packedTree.Bytes()) // Запишем дерево в файл

	if err != nil {
		return err
	}

	fileBuffer := bitsbuffer.NewEmptyFlushableBuffer(writer)

	fileBuffer.AddFromBuffer(remainingBuffer)

	//err = writer.Flush() TODO: можно флашнуть буфер после записи дерева

	// Записываем коды в файл
	for {
		newByte1, err := reader.ReadByte()

		if err == io.EOF {
			break
		}

		newByte2, err := reader.ReadByte()

		if err == io.EOF {
			newByte2 = 0
		}

		pair := (uint16(newByte1) << 8) | uint16(newByte2)
		code := encoder.tree.GetCode(pair)
		fileBuffer.AddFromBuffer(code)
	}

	fileBuffer.Flush()

	return nil
}

// Считаем кол-во уникальных символов - чтобы в декодере понять кол-во вершин в дереве
// Считаем кол-во используемых бит под данные в последнем байте файла (сумма всего файла до этого % 8)
func (encoder *Encoder) calculateMetaParams() {
	var bitsInLastByte int64
	var uniqueWords uint32

	for _, word := range encoder.tree.Words {
		if word == nil {
			continue
		}
		uniqueWords++
		lenTotal := int64(word.Length()) * word.Count() % 8
		bitsInLastByte += lenTotal % 8 // Можно не делать %8 каждую итерацию, а только один раз в конце, менее эффективно по памяти
	}

	treeSizeInBits := uniqueWords*10 - 1
	bitsInLastByte = (bitsInLastByte + int64(treeSizeInBits)) % 8

	encoder.bitsInLastByte = uint8(bitsInLastByte)

	encoder.uniqueWords = uniqueWords
}

func (encoder *Encoder) GetAverageBitsPerWord() float64 {
	var sum float64
	var count float64

	for _, code := range encoder.tree.Words {
		if code != nil {
			sum += float64(code.Length()) * float64(code.Count())
			count += float64(code.Count())
		}
	}

	return sum / count
}
