package regen

import (
	crand "crypto/rand"
	"io"
)

type rand struct {
	randomSource io.Reader
}

func NewRand() *rand {
	return newRand(crand.Reader)
}

// Internal initialization function to allow testing to inject its own reader
func newRand(reader io.Reader) *rand {
	return &rand{
		randomSource: reader,
	}
}

// Returns a non-negative random number in [0,max) from the crypto/rand source. Panics if max <= 0
func (r rand) Int31n(max int32) int32 {
	if max <= 0 {
		panic("Max must be greater than 0")
	}

	bytesToRead := byteLenInt32(max)
	bitsToShift := uint(bytesToRead*8 - 1)
	ceil := int32((1 << bitsToShift) - 1 - (1<<bitsToShift)%uint32(max))

	v := r.readBytes(bytesToRead)
	for v > ceil {
		v = r.readBytes(bytesToRead)
	}
	return v % max
}

// Reads byteCount bytes from the internal Reader
func (r rand) readBytes(byteCount int) int32 {
	bytes := make([]byte, byteCount)

	if _, err := r.randomSource.Read(bytes); err != nil {
		panic(err)
	}

	var result uint32
	for index := 0; index < byteCount; index++ {
		result |= uint32(bytes[index]) << uint(8*index)
	}
	signBitIndex := uint(byteCount*8 - 1)
	result &^= (1 << signBitIndex)

	return int32(result)
}

func byteLenInt32(n int32) int {
	if b := n >> 24; b != 0 {
		return 4
	}

	if b := n >> 16; b != 0 {
		return 3
	}

	if b := n >> 8; b != 0 {
		return 2
	}

	return 1
}
