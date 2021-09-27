package benchmarks

import (
	"bytes"
	"testing"

	json "encoding/json"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/chain"

	eoscanada "github.com/eoscanada/eos-go"
)

// sanity checks that we are actually decoding the test data correctly

func TestAbiDecode(t *testing.T) {
	abi := loadAbi(transferAbiJson)

	rv, err := abi.Decode(bytes.NewReader(testTransferData), "transfer")
	assert.NoError(t, err)
	assert.Equal(t, rv, map[string]interface{}{
		"from":     chain.N("foo"),
		"to":       chain.N("bar"),
		"quantity": *chain.A("1.0000 EOS"),
		"memo":     "hello",
	})
}

func TestAbiDecodeEosCanada(t *testing.T) {
	abi, err := eoscanada.NewABI(bytes.NewReader([]byte(transferAbiJson)))
	assert.NoError(t, err)
	jsonData, err := abi.DecodeAction(testTransferData, "transfer")
	assert.NoError(t, err)
	var rv map[string]interface{}
	json.Unmarshal(jsonData, &rv)
	assert.Equal(t, rv, map[string]interface{}{
		"from":     "foo",
		"to":       "bar",
		"quantity": "1.0000 EOS",
		"memo":     "hello",
	})
}

// benchmarks
// not apples to apples since eos canada returns JSON bytes directly while ours returns a map

func Benchmark_Decode_AbiDef(b *testing.B) {
	abi := loadAbi(transferAbiJson)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := abi.Decode(bytes.NewReader(testTransferData), "transfer")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Decode_AbiDef_EosCanada(b *testing.B) {
	abi, err := eoscanada.NewABI(bytes.NewReader([]byte(transferAbiJson)))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := abi.DecodeAction(testTransferData, "transfer")
		if err != nil {
			b.Fatal(err)
		}

	}
}
