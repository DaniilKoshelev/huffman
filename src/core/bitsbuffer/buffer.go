package bitsbuffer

import (
	"encoding/binary"
	"io"
)

const bufferSize = 16
const bufferEmpty = -1

type Buffer struct {
	currentBit int8
	bits       uint16
	flusher
}

func (buf *Buffer) Bits() uint16 {
	return buf.bits
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
	buf.bits |= 0 << (bufferSize - 1 - buf.currentBit - 1)
	buf.currentBit++

	if buf.currentBit == bufferSize {
		buf.flush()
	}

	return buf
}

func (buf *Buffer) AddOne() *Buffer {
	buf.bits |= 1 << (bufferSize - 1 - buf.currentBit - 1)
	buf.currentBit++

	if buf.currentBit == bufferSize {
		buf.flush()
	}

	return buf
}

func (buf *Buffer) reset() *Buffer {
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
	flusher.buf.reset()
}

func (flusher *emptyFlusher) flushBufferFinal() {
	flusher.buf.reset()
}

func (flusher *ioFlusher) flushBuffer() {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, flusher.buf.bits)
	_, _ = flusher.whereToFlush.Write(bytes) // TODO: process error

	flusher.buf.reset()
}

func (flusher *ioFlusher) flushBufferFinal() {
	if flusher.buf.currentBit == bufferEmpty {
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

	flusher.buf.reset()
}

func (buf *Buffer) AddFromBuffer(anotherBuf *Buffer) *Buffer {
	currentBitBeforeFlush := buf.currentBit + 1 + anotherBuf.currentBit + 1 - bufferSize - 1
	l := anotherBuf.bits >> (buf.currentBit + 1)
	r := anotherBuf.bits << (bufferSize - buf.currentBit - 1)

	buf.bits |= l

	if buf.currentBit+1+anotherBuf.currentBit+1 >= bufferSize {
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
