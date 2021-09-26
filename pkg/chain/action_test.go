package chain_test

import (
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/chain"
)

func TestAction(t *testing.T) {
	action := chain.NewAction(
		chain.N("eosio.token"),
		chain.N("transfer"),
		[]chain.PermissionLevel{
			{chain.N("alice"), chain.N("active")},
		},
		chain.Bytes{
			0xbe, 0xef, 0xfa, 0xce,
		},
	)
	assert.ABICoding(t, action, []byte{
		0x00, 0xa6, 0x82, 0x34, 0x03, 0xea, 0x30, 0x55, 0x00, 0x00, 0x00, 0x57, 0x2d, 0x3c, 0xcd, 0xcd,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x85, 0x5c, 0x34, 0x00, 0x00, 0x00, 0x00, 0xa8, 0xed, 0x32,
		0x32, 0x04, 0xbe, 0xef, 0xfa, 0xce,
	})
	assert.JSONCoding(t, action, `
		{
			"account": "eosio.token",
			"name": "transfer",
			"authorization": [
				{
					"actor": "alice",
					"permission": "active"
				}
			],
			"data": "beefface"
		}
	`)
}
