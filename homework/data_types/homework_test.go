package main

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type ConvertionTypes interface {
	uint16 | uint32 | uint64
}

func toLittleEndian[T ConvertionTypes](number T) (res T) {
	switch any(number).(type) {
	case uint16:
		b := make([]byte, 2)
		n := uint16(number)
		binary.BigEndian.PutUint16(b, n)
		res = T(binary.LittleEndian.Uint16(b))
	case uint32:
		b := make([]byte, 4)
		n := uint32(number)
		binary.BigEndian.PutUint32(b, n)
		res = T(binary.LittleEndian.Uint32(b))
	case uint64:
		b := make([]byte, 8)
		n := uint64(number)
		binary.BigEndian.PutUint64(b, n)
		res = T(binary.LittleEndian.Uint64(b))
	}
	return
}

func toLittleEndian32(number uint32) uint32 {
	W := byte(number >> 24)
	X := byte(number >> 16)
	Y := byte(number >> 8)
	Z := byte(number)
	return uint32(W) | uint32(X)<<8 | uint32(Y)<<16 | uint32(Z)<<24
}

func ToLittleEndian(number uint32) uint32 {
	//return toLittleEndian32(number)
	return toLittleEndian(number)
}

func TestСonversion(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
