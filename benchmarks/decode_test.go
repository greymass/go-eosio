package benchmarks

import (
	"bytes"
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/abi"
	"github.com/greymass/go-eosio/pkg/chain"

	eoscanada "github.com/eoscanada/eos-go"
)

// sanity checks that we are actually decoding the test data correctly

func TestDecode(t *testing.T) {
	var err error
	var tx chain.Transaction

	err = chain.NewDecoder(bytes.NewReader(testTransactionData)).Decode(&tx)
	assert.NoError(t, err)
	assert.Equal(t, tx.Actions[len(tx.Actions)-1].Account.String(), "greymass")
}

func TestDecodeEosCanada(t *testing.T) {
	var err error
	var tx eoscanada.Transaction

	err = eoscanada.NewDecoder(testTransactionData).Decode(&tx)
	assert.NoError(t, err)
	assert.Equal(t, tx.Actions[len(tx.Actions)-1].Account, eoscanada.AccountName("greymass"))
}

// benchmarks

func Benchmark_Decode(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		var tx chain.Transaction
		err = chain.NewDecoder(bytes.NewReader(testTransactionData)).Decode(&tx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Decode_NoOptimize(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		var tx Transaction
		err = abi.NewDecoder(bytes.NewReader(testTransactionData), noopDecode).Decode(&tx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Decode_EosCanada(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		var tx eoscanada.Transaction
		err = eoscanada.NewDecoder(testTransactionData).Decode(&tx)
		if err != nil {
			b.Fatal(err)
		}
	}
}
