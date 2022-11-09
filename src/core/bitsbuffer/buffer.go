package bitsbuffer

import "io"

type Buffer struct {
	currentBit   int8
	bits         byte
	whereToFlush io.ByteWriter
}

func NewBuffer(currentBit int8, bits byte, whereToFlush io.ByteWriter) *Buffer {
	return &Buffer{currentBit: currentBit, bits: bits, whereToFlush: whereToFlush}
}

func (buf *Buffer) AddByte(byteToAdd byte) *Buffer {
	currentBitBeforeFlush := buf.currentBit

	// Splitting byte into two parts
	l := byteToAdd >> buf.currentBit
	r := byteToAdd << (8 - buf.currentBit)

	buf.bits |= l

	buf.Flush()

	buf.bits = r
	buf.currentBit = currentBitBeforeFlush

	return buf
}

func (buf *Buffer) AddZero() *Buffer {
	buf.bits |= 0 << (7 - buf.currentBit)
	buf.currentBit++

	if buf.currentBit == 8 {
		buf.Flush()
	}

	return buf
}

func (buf *Buffer) AddOne() *Buffer {
	buf.bits |= 1 << (7 - buf.currentBit)
	buf.currentBit++

	if buf.currentBit == 8 {
		buf.Flush()
	}

	return buf
}

func (buf *Buffer) reset() *Buffer {
	buf.bits = 0
	buf.currentBit = 0

	return buf
}

func (buf *Buffer) Flush() {
	_ = buf.whereToFlush.WriteByte(buf.bits) // TODO: process error

	buf.reset()
}
