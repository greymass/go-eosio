package benchmarks

import (
	"encoding/json"

	"github.com/greymass/go-eosio/pkg/abi"
	"github.com/greymass/go-eosio/pkg/chain"
)

func loadAbi(v string) *chain.Abi {
	var rv chain.Abi
	err := json.Unmarshal([]byte(v), &rv)
	if err != nil {
		panic(err)
	}
	return &rv
}

func noopDecode(dec *abi.Decoder, v interface{}) (done bool, err error) {
	return false, nil
}

func noopEncode(enc *abi.Encoder, v interface{}) (done bool, err error) {
	return false, nil
}
