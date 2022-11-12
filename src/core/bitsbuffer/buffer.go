package bitsbuffer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const bufferSize = 16
const bufferEmpty = -1

type Buffer struct {
	currentBit int8
	bits       uint16
	flusher
	reader
}

func (buf *Buffer) Bits() uint16 {
	return buf.bits
}

func (buf *Buffer) SetIoReader(whereToReadFrom io.ByteReader) *Buffer {
	buf.reader = newIoReader(buf, whereToReadFrom)

	return buf
}

func NewEmptyFlushableBuffer(whereToFlush io.Writer) *Buffer {
	buf := &Buffer{currentBit: bufferEmpty, bits: 0}
	buf.flusher = newIoFlusher(buf, whereToFlush)

	return buf
}

func NewBuffer(currentBit int8, bits uint16) *Buffer {
	buf := &Buffer{currentBit: currentBit, bits: bits}
	buf.flusher = newEmptyFlusher(buf)

	return buf
}

func NewEmptyBuffer() *Buffer {
	buf := &Buffer{currentBit: bufferEmpty, bits: 0}
	buf.flusher = newEmptyFlusher(buf)

	return buf
}

func (buf *Buffer) ToInt() int {
	return int(buf.bits) | int(buf.currentBit)<<16
}

func From(buffer *Buffer) *Buffer {
	// TODO: скопировать средствами golang
	newBuffer := new(Buffer)

	newBuffer.bits = buffer.bits
	newBuffer.currentBit = buffer.currentBit
	newBuffer.flusher = buffer.flusher

	return newBuffer
}

func (buf *Buffer) AddByte(byteToAdd byte) *Buffer {
	newBuf := NewEmptyBuffer()
	newBuf.currentBit = 7
	newBuf.bits = uint16(byteToAdd) << 8

	buf.AddFromBuffer(newBuf)

	return buf
}

func (buf *Buffer) AddZero() *Buffer {
	return buf.AddBit(0)
}

func (buf *Buffer) AddOne() *Buffer {
	return buf.AddBit(1)
}

func (buf *Buffer) AddBit(bit uint8) *Buffer {
	if buf.currentBit == bufferSize-1 {
		buf.flush()
	}

	buf.bits |= uint16(bit) << (bufferSize - 1 - buf.currentBit - 1)
	buf.currentBit++

	return buf
}

func (buf *Buffer) Reset() *Buffer {
	buf.bits = 0
	buf.currentBit = bufferEmpty

	return buf
}

// Эффективный внутренний флаш только при переполнении буфера (флаш 2 байт сразу)
func (buf *Buffer) flush() *Buffer {
	buf.flusher.flushBuffer()

	return buf
}

// Flush Неэффективный флаш, должен использоваться однократно, когда заканчиваем работу с буфером
// Пишет 0, 1, 2 байта, в зависимости от заполненности
func (buf *Buffer) Flush() *Buffer {
	buf.flusher.flushBufferFinal()

	return buf
}

func (flusher *emptyFlusher) flushBuffer() {
	//flusher.buf.Reset()
}

func (flusher *emptyFlusher) flushBufferFinal() {
	//flusher.buf.Reset()
}

func (flusher *ioFlusher) flushBuffer() {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, flusher.buf.bits)
	_, _ = flusher.whereToFlush.Write(bytes) // TODO: process error

	flusher.buf.Reset()
}

func (flusher *ioFlusher) flushBufferFinal() {
	if flusher.buf.isEmpty() {
		// do nothing
	} else if flusher.buf.currentBit <= 7 {
		byteToWrite := byte(flusher.buf.bits >> 8)
		bytes := []byte{byteToWrite}
		_, _ = flusher.whereToFlush.Write(bytes) // TODO: process error
	} else {
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, flusher.buf.bits)
		_, _ = flusher.whereToFlush.Write(bytes) // TODO: process error
	}

	flusher.buf.Reset()
}

func (buf *Buffer) AddFromBuffer(anotherBuf *Buffer) *Buffer {
	currentBitBeforeFlush := buf.currentBit + 1 + anotherBuf.currentBit + 1 - bufferSize - 1
	l := anotherBuf.bits >> (buf.currentBit + 1)
	r := anotherBuf.bits << (bufferSize - buf.currentBit - 1)

	buf.bits |= l

	if buf.currentBit+1+anotherBuf.currentBit+1 > bufferSize {
		buf.flush()
		buf.bits = r
		buf.currentBit = currentBitBeforeFlush

		return buf
	}

	buf.currentBit += anotherBuf.currentBit + 1

	return buf
}

func (buf *Buffer) Length() int8 {
	return buf.currentBit + 1
}

type flusher interface {
	flushBuffer()
	flushBufferFinal()
}

type reader interface {
	read() error
}

type ioReader struct {
	buf             *Buffer
	whereToReadFrom io.ByteReader
}

func newIoReader(buf *Buffer, whereToReadFrom io.ByteReader) *ioReader {
	return &ioReader{buf: buf, whereToReadFrom: whereToReadFrom}
}

func (reader *ioReader) read() error {
	newByte, err := reader.whereToReadFrom.ReadByte()

	if err != nil {
		return err
	}

	reader.buf.AddByte(newByte)

	return nil
}

type emptyFlusher struct {
	buf *Buffer
}

func newEmptyFlusher(buf *Buffer) *emptyFlusher {
	return &emptyFlusher{buf: buf}
}

type ioFlusher struct {
	buf          *Buffer
	whereToFlush io.Writer
}

func newIoFlusher(buf *Buffer, whereToFlush io.Writer) *ioFlusher {
	return &ioFlusher{buf: buf, whereToFlush: whereToFlush}
}

func (buf *Buffer) isEmpty() bool {
	return buf.currentBit == bufferEmpty
}

func (buf *Buffer) ReadBit() (uint8, error) {
	if buf.isEmpty() {
		err := buf.read()

		if err != nil {
			return 0, errors.New(fmt.Sprintln("error: reading from an empty buffer"))
		}
	}

	bit := uint8((buf.bits >> 15) & 1)

	buf.bits <<= 1
	buf.currentBit--

	return bit, nil
}

// Scan TODO: rename (Похоже на костыль)
func (buf *Buffer) Scan() error {
	if buf.currentBit <= 7 {
		err := buf.read()

		return err
	}

	return errors.New(fmt.Sprintln("error: not enough space in buffer"))
}

func (buf *Buffer) ReadByte() (uint8, error) {
	if buf.isEmpty() || buf.currentBit < 7 {
		err := buf.read()

		if err != nil {
			return 0, errors.New(fmt.Sprintln("error: reading from an empty buffer"))
		}
	}

	byteToReturn := uint8(buf.bits >> 8)

	buf.bits <<= 8
	buf.currentBit -= 8

	return byteToReturn, nil
}
