package core

import (
	"bufio"
	"errors"
	"huffman/src/core/bitsbuffer"
	"huffman/src/core/tree"
	"io"
)

type Encoder struct {
	tree           *tree.Tree
	bitsInLastByte uint8
	uniqueWords    uint8
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

	_, err := writer.Write(
		[]byte{
			encoder.uniqueWords, // 1-й байт - кол-во уникальных символов

			// TODO: использовать 3 бита вместо целого байта
			encoder.bitsInLastByte, // 2-й байт - кол-во полезных бит в последнем байте файла (info: возможно неправильно, придется потом заменить на общее число слов в изаначальном файле)
		},
	)

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

// Считаем кол-во уникальных символов - чтобы в декодере понять кол-во вершин в дереве
// Считаем кол-во используемых бит под данные в последнем байте файла (сумма всего файла до этого % 8)
func (encoder *Encoder) calculateMetaParams() {
	var bitsInLastByte int64
	var uniqueWords uint16

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

	// TODO: костыль
	if uniqueWords == 256 {
		uniqueWords = 0
	}

	encoder.uniqueWords = uint8(uniqueWords)
}
