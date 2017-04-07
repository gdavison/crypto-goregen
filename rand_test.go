package regen

import (
	"bytes"
	"math"
	"testing"
	"testing/quick"
)

func TestInt31nShouldPanicWhenMaxZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected Int31n() to panic")
		}
	}()
	theRand := NewRand()
	theRand.Int31n(0)
}

func TestInt31nShouldReturnTheValueFromStreamIfLessThanMax(t *testing.T) {
	value := byte(0x02)
	slice := []byte{value}
	reader := bytes.NewReader(slice)
	theRand := newRand(reader)

	result := theRand.Int31n(3)
	if result != int32(value) {
		t.Fatalf("Expected %d, was %d", value, result)
	}
}

func TestInt31nShouldReturnAModuleValueFromStreamIfGreaterThanMax(t *testing.T) {
	value := byte(0x04)
	slice := []byte{value}
	reader := bytes.NewReader(slice)
	theRand := newRand(reader)
	expected := int32(value % 3)

	result := theRand.Int31n(3)
	if result != expected {
		t.Fatalf("Expected %d, was %d", expected, result)
	}
}

func TestInt31nShouldReadNextValueFromStreamIfThereWouldBeModuloBias(t *testing.T) {
	value := byte(0x65)
	slice := []byte{0xff, value}
	reader := bytes.NewReader(slice)
	theRand := newRand(reader)
	expected := int32(value % 3)

	result := theRand.Int31n(3)
	if result != expected {
		t.Fatalf("Expected %d, was %d", expected, result)
	}
}

func TestReadingSingleBytes(t *testing.T) {
	assertion := func(x byte) bool {
		slice := []byte{x}
		reader := bytes.NewReader(slice)
		theRand := newRand(reader)
		expected := int32(x &^ (1 << 7))

		result := theRand.readBytes(1)
		return result == expected
	}
	if err := quick.Check(assertion, nil); err != nil {
		t.Error(err)
	}
}

func TestReadingSingleByteRemovesSignBit(t *testing.T) {
	value := byte(0xff)
	slice := []byte{value}
	reader := bytes.NewReader(slice)
	theRand := newRand(reader)
	expected := int32(0x7f)

	result := theRand.readBytes(1)
	if result != expected {
		t.Fatalf("Expected %d, was %d", expected, result)
	}
}

func TestReadingTwoBytes(t *testing.T) {
	assertion := func(x [2]byte) bool {
		reader := bytes.NewReader(x[:])
		theRand := newRand(reader)
		var expected int32
		expected = int32(x[0]) | (int32(x[1]&^(1<<7)) << 8)

		result := theRand.readBytes(2)
		return result == int32(expected)
	}
	if err := quick.Check(assertion, nil); err != nil {
		t.Error(err)
	}
}

func TestReadingThreeBytes(t *testing.T) {
	assertion := func(x [3]byte) bool {
		reader := bytes.NewReader(x[:])
		theRand := newRand(reader)
		var expected int32
		expected = int32(uint32(x[0]) | (uint32(x[1]) << 8) | (uint32(x[2]&^(1<<7)) << 16))

		result := theRand.readBytes(3)
		return result == int32(expected)
	}
	if err := quick.Check(assertion, nil); err != nil {
		t.Error(err)
	}
}

func TestReadingFourBytes(t *testing.T) {
	assertion := func(x [4]byte) bool {
		reader := bytes.NewReader(x[:])
		theRand := newRand(reader)
		var expected int32
		expected = int32(uint32(x[0]) | (uint32(x[1]) << 8) | (uint32(x[2]) << 16) | (uint32(x[3]&^(1<<7)) << 24))

		result := theRand.readBytes(4)
		return result == int32(expected)
	}
	if err := quick.Check(assertion, nil); err != nil {
		t.Error(err)
	}
}

func TestIntnShouldPanicWhenMaxZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected Intn() to panic")
		}
	}()
	theRand := NewRand()
	theRand.Intn(0)
}

func TestIntnHasMaxRangeInt32(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected Intn() to panic")
		}
	}()
	theRand := NewRand()
	theRand.Intn(int(math.MaxInt32 + 1))
}
