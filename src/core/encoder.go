package core

import (
	"bufio"
	"errors"
	"io"
)

type ensemble struct {
	bytes map[byte]int64
}

type Encoder struct {
	reader *bufio.Reader
}

func NewEncoder(reader *bufio.Reader) *Encoder {
	encoder := new(Encoder)
	encoder.reader = reader

	return encoder
}

func (encoder *Encoder) BuildTree() error {
	if encoder.reader == nil {
		return errors.New("reader is not set")
	}

	for {
		newByte, err := encoder.reader.ReadByte()

		if err == io.EOF {
			break
		}
	}

	return nil
}

func (encoder *Encoder) Encode(ch chan []byte) {

}
