package chain_test

import (
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/chain"
)

func TestBytes(t *testing.T) {
	bytes := chain.Bytes([]byte{0xbe, 0xef, 0xfa, 0xce})

	assert.ABICoding(t, bytes, []byte{0x04, 0xbe, 0xef, 0xfa, 0xce})
	assert.JSONCoding(t, bytes, `"beefface"`)
}
