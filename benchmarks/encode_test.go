package benchmarks

import (
	"bytes"
	"io"
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/abi"
	"github.com/greymass/go-eosio/pkg/chain"

	eoscanada "github.com/eoscanada/eos-go"
)

// sanity checks that we are actually decoding the test data correctly

func TestEncode(t *testing.T) {
	var err error

	b := bytes.NewBuffer(nil)
	err = chain.NewEncoder(b).Encode(testTransaction)
	assert.NoError(t, err)
	assert.Equal(t, b.Bytes(), testTransactionData)
}

func TestEncodeEosCanada(t *testing.T) {
	var err error
	b := bytes.NewBuffer(nil)
	err = eoscanada.NewEncoder(b).Encode(testTransactionCanada)
	assert.NoError(t, err)
	assert.Equal(t, b.Bytes(), testTransactionData)
}

// benchmarks

func Benchmark_Encode(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = chain.NewEncoder(io.Discard).Encode(testTransaction)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Encode_NoOptimize(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = abi.NewEncoder(io.Discard, noopEncode).Encode(testTransactionNoOptimize)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Encode_EosCanada(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = eoscanada.NewEncoder(io.Discard).Encode(testTransactionCanada)
		if err != nil {
			b.Fatal(err)
		}
	}
}
