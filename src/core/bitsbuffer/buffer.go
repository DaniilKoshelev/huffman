package bitsbuffer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

const bufferSize = 64
const bufferEmpty = -1

type Buffer struct {
	currentBit int8
	bits       uint64 // бафер должен уметь вмещать до 256 бит, поэтому решение с int64 не самое оптимальное, стоит сделать []byte
	flusher
	reader
}

func (buf *Buffer) Bits() uint64 {
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

func NewBuffer(currentBit int8, bits uint64) *Buffer {
	buf := &Buffer{currentBit: currentBit, bits: bits}
	buf.flusher = newEmptyFlusher(buf)

	return buf
}

func NewEmptyBuffer() *Buffer {
	buf := &Buffer{currentBit: bufferEmpty, bits: 0}
	buf.flusher = newEmptyFlusher(buf)

	return buf
}

func (buf *Buffer) ToString() string {
	return fmt.Sprintf("%d_%d", buf.bits, buf.currentBit)
}

func (buf *Buffer) DropLastBits(bitsToDrop int8) {
	buf.currentBit -= bitsToDrop
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
	newBuf.bits = uint64(byteToAdd) << 56

	buf.AddFromBuffer(newBuf)

	return buf
}

func (buf *Buffer) AddUInt16(uint uint16) *Buffer {
	newBuf := NewEmptyBuffer()
	newBuf.currentBit = 15
	newBuf.bits = uint64(uint) << 48

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

	buf.bits |= uint64(bit) << (bufferSize - 1 - buf.currentBit - 1)
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
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, flusher.buf.bits)
	_, _ = flusher.whereToFlush.Write(bytes) // TODO: process error

	flusher.buf.Reset()
}

//TODO: fix conditions
func (flusher *ioFlusher) flushBufferFinal() {
	if flusher.buf.IsEmpty() {
		// do nothing
	} else if flusher.buf.currentBit <= 7 {
		byteToWrite := byte(flusher.buf.bits >> 56)
		bytes := []byte{byteToWrite}
		_, _ = flusher.whereToFlush.Write(bytes) // TODO: process error
	} else {
		bytesAmount := int(math.Ceil(float64(flusher.buf.currentBit+1) / 8.0))
		var bytes []byte
		for i := 0; i < bytesAmount; i++ {
			bytes = append(bytes, byte(flusher.buf.bits>>(8*(7-uint64(i)))))
		}
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

func (buf *Buffer) IsEmpty() bool {
	return buf.currentBit == bufferEmpty
}

func (buf *Buffer) ReadBit() (uint8, error) {
	if buf.IsEmpty() {
		err := buf.read()

		if err != nil {
			return 0, errors.New(fmt.Sprintln("error: reading from an empty buffer"))
		}
	}

	bit := uint8((buf.bits >> 63) & 1)

	buf.bits <<= 1
	buf.currentBit--

	return bit, nil
}

// Scan TODO: rename (Похоже на костыль)
func (buf *Buffer) Scan() error {
	if buf.currentBit <= 55 {
		err := buf.read()

		return err
	}

	return errors.New(fmt.Sprintln("error: not enough space in buffer"))
}

func (buf *Buffer) ReadByte() (uint8, error) {
	if buf.IsEmpty() || buf.currentBit < 55 {
		err := buf.read()

		if err != nil {
			return 0, errors.New(fmt.Sprintln("error: reading from an empty buffer"))
		}
	}

	byteToReturn := uint8(buf.bits >> 56)

	buf.bits <<= 8
	buf.currentBit -= 8

	return byteToReturn, nil
}
