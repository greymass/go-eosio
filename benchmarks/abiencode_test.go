package benchmarks

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	eoscanada "github.com/eoscanada/eos-go"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/chain"
)

// sanity checks that we are actually decoding the test data correctly

func TestAbiEncode(t *testing.T) {
	abi := loadAbi(transferAbiJson)

	buf := bytes.NewBuffer(nil)
	err := abi.Encode(buf, "transfer", map[string]interface{}{
		"from":     chain.N("foo"),
		"to":       chain.N("bar"),
		"quantity": *chain.A("1.0000 EOS"),
		"memo":     "hello",
	})
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), testTransferData)
}

func TestAbiEncodeEosCanada(t *testing.T) {
	abi, err := eoscanada.NewABI(bytes.NewReader([]byte(transferAbiJson)))
	assert.NoError(t, err)
	jsonData, err := json.Marshal(map[string]interface{}{
		"from":     chain.N("foo"),
		"to":       chain.N("bar"),
		"quantity": *chain.A("1.0000 EOS"),
		"memo":     "hello",
	})
	assert.NoError(t, err)
	bytes, err := abi.EncodeAction("transfer", jsonData)
	assert.NoError(t, err)
	assert.Equal(t, bytes, testTransferData)
}

// benchmarks

func Benchmark_Encode_AbiDef(b *testing.B) {
	abi := loadAbi(transferAbiJson)
	v := map[string]interface{}{
		"from":     chain.N("foo"),
		"to":       chain.N("bar"),
		"quantity": *chain.A("1.0000 EOS"),
		"memo":     "hello",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := abi.Encode(io.Discard, "transfer", v)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Encode_AbiDef_EosCanada(b *testing.B) {
	abi, err := eoscanada.NewABI(bytes.NewReader([]byte(transferAbiJson)))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsonData, err := json.Marshal(map[string]interface{}{
			"from":     chain.N("foo"),
			"to":       chain.N("bar"),
			"quantity": *chain.A("1.0000 EOS"),
			"memo":     "hello",
		})
		if err != nil {
			b.Fatal(err)
		}
		_, err = abi.EncodeAction("transfer", jsonData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// func BenchmarkAbiDecodeEosCanada(b *testing.B) {
// 	abi, err := eoscanada.NewABI(bytes.NewReader([]byte(transferAbiJson)))
// 	if err != nil {
// 		b.Fatal(err)
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_, err := abi.DecodeAction(testTransferData, "transfer")
// 		if err != nil {
// 			b.Fatal(err)
// 		}

// 	}
// }
