package abi_test

import (
	"bytes"
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/abi"
)

func noopEncodefunc(enc *abi.Encoder, v interface{}) (done bool, err error) {
	return false, nil
}

func TestEncodeUint64(t *testing.T) {
	var v uint64 = 42
	var b *bytes.Buffer = bytes.NewBuffer(nil)
	err := abi.NewEncoder(b, noopEncodefunc).Encode(v)
	assert.NoError(t, err)
	assert.Equal(t, b.Bytes(), []byte{0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestEncodeStruct(t *testing.T) {
	type TestStruct struct {
		A uint64
		B *uint64 `eosio:"optional"`
		C *uint64 `eosio:"optional"`
	}
	x := uint64(42)
	v := TestStruct{42, nil, &x}
	var b *bytes.Buffer = bytes.NewBuffer(nil)
	err := abi.NewEncoder(b, noopEncodefunc).Encode(v)
	assert.NoError(t, err)
	assert.Equal(t, b.Bytes(), []byte{
		0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
		0x01, 0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
}

func TestEncodeVector(t *testing.T) {
	v := []uint64{42, 43, 0}
	var b *bytes.Buffer = bytes.NewBuffer(nil)
	err := abi.NewEncoder(b, noopEncodefunc).Encode(v)
	assert.NoError(t, err)
	assert.Equal(t, b.Bytes(), []byte{
		0x03,
		0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x2b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
}

func TestEncodeArray(t *testing.T) {
	v := [3]uint64{42, 43, 0}
	var b *bytes.Buffer = bytes.NewBuffer(nil)
	err := abi.NewEncoder(b, noopEncodefunc).Encode(v)
	assert.NoError(t, err)
	assert.Equal(t, b.Bytes(), []byte{
		0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x2b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
}
