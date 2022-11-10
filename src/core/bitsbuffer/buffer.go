package bitsbuffer

import (
	"encoding/binary"
	"io"
)

const bufferSize = 16

type Buffer struct {
	currentBit int8
	bits       uint16
	flusher
}

func (buf *Buffer) Bits() uint16 {
	return buf.bits
}

func NewFlushableBuffer(whereToFlush io.Writer) *Buffer {
	buf := &Buffer{currentBit: 0, bits: 0}
	buf.flusher = newIoFlusher(buf, whereToFlush)

	return buf
}

func NewBuffer() *Buffer {
	buf := &Buffer{currentBit: 0, bits: 0}
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
	newBuf := NewBuffer()
	newBuf.currentBit = 7
	newBuf.bits = uint16(byteToAdd) << 8

	buf.AddFromBuffer(newBuf)

	return buf
}

func (buf *Buffer) AddZero() *Buffer {
	buf.bits |= 0 << (bufferSize - 1 - buf.currentBit)
	buf.currentBit++

	if buf.currentBit == bufferSize {
		buf.Flush()
	}

	return buf
}

func (buf *Buffer) AddOne() *Buffer {
	buf.bits |= 1 << (bufferSize - 1 - buf.currentBit)
	buf.currentBit++

	if buf.currentBit == bufferSize {
		buf.Flush()
	}

	return buf
}

func (buf *Buffer) reset() *Buffer {
	buf.bits = 0
	buf.currentBit = 0

	return buf
}

func (buf *Buffer) Flush() *Buffer {
	buf.flusher.FlushBuffer()

	return buf
}

func (flusher *emptyFlusher) FlushBuffer() {
	flusher.buf.reset()
}

func (flusher *ioFlusher) FlushBuffer() {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, flusher.buf.bits)
	_, _ = flusher.whereToFlush.Write(bytes) // TODO: process error

	flusher.buf.reset()
}

func (buf *Buffer) AddFromBuffer(anotherBuf *Buffer) *Buffer {
	currentBitBeforeFlush := buf.currentBit

	l := anotherBuf.bits >> buf.currentBit
	r := anotherBuf.bits << (bufferSize - buf.currentBit)

	buf.bits |= l

	buf.Flush()

	buf.bits = r
	buf.currentBit = currentBitBeforeFlush

	return buf
}

type flusher interface {
	FlushBuffer()
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
