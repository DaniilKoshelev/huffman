package core

import (
	"bufio"
	"errors"
	"io"
)

type ensemble struct {
	bytes map[byte]cell
}

type cell struct {
	count int64
	code  byte
}

type Encoder struct {
	*ensemble
	reader *bufio.Reader
}

func NewEncoder(reader *bufio.Reader) *Encoder {
	encoder := new(Encoder)
	encoder.reader = reader
	encoder.ensemble = new(ensemble)
	encoder.ensemble.bytes = make(map[byte]cell)

	return encoder
}

func (encoder *Encoder) BuildTree() error {
	err := encoder.countBytes()

	if err != nil {
		return err
	}

	return nil
}

func (encoder *Encoder) countBytes() error {
	if encoder.reader == nil {
		return errors.New("reader is not set")
	}

	for {
		newByte, err := encoder.reader.ReadByte()

		if err == io.EOF {
			break
		}

		curCell, exists := encoder.bytes[newByte]

		if exists {
			curCell.count++
		} else {
			encoder.bytes[newByte] = cell{1, 0}
		}
	}

	return nil
}

func (encoder *Encoder) Encode(ch chan []byte) {

}
