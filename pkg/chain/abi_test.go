package chain_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
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

var tokenAbi = loadAbi(`
{
    "version": "eosio::abi/1.1",
    "types": [],
    "structs": [
        {
            "name": "account",
            "base": "",
            "fields": [
                {
                    "name": "balance",
                    "type": "asset"
                }
            ]
        },
        {
            "name": "banana",
            "base": "",
            "fields": [
                {
                    "name": "moo",
                    "type": "name"
                }
            ]
        },
        {
            "name": "create",
            "base": "",
            "fields": [
                {
                    "name": "issuer",
                    "type": "name"
                },
                {
                    "name": "maximum_supply",
                    "type": "asset"
                }
            ]
        },
        {
            "name": "currency_stats",
            "base": "",
            "fields": [
                {
                    "name": "supply",
                    "type": "asset"
                },
                {
                    "name": "max_supply",
                    "type": "asset"
                },
                {
                    "name": "issuer",
                    "type": "name"
                }
            ]
        },
        {
            "name": "issue",
            "base": "",
            "fields": [
                {
                    "name": "to",
                    "type": "name"
                },
                {
                    "name": "quantity",
                    "type": "asset"
                },
                {
                    "name": "memo",
                    "type": "string"
                }
            ]
        },
        {
            "name": "open",
            "base": "",
            "fields": [
                {
                    "name": "owner",
                    "type": "name"
                },
                {
                    "name": "symbol",
                    "type": "symbol"
                },
                {
                    "name": "ram_payer",
                    "type": "name"
                }
            ]
        },
        {
            "name": "megatransfer",
            "base": "transfer",
            "fields": [
                {
                    "name": "extra",
                    "type": "mega"
                },
				        {
                    "name": "extra2",
                    "type": "banana[]"
                }
            ]
        },
        {
            "name": "transfer",
            "base": "",
            "fields": [
                {
                    "name": "from",
                    "type": "name"
                },
                {
                    "name": "to",
                    "type": "name"
                },
                {
                    "name": "quantity",
                    "type": "asset"
                },
                {
                    "name": "memo",
                    "type": "string"
                }
            ]
        }
    ],
    "actions": [
        {
            "name": "close",
            "type": "close",
            "ricardian_contract": ""
        },
        {
            "name": "create",
            "type": "create",
            "ricardian_contract": ""
        },
        {
            "name": "issue",
            "type": "issue",
            "ricardian_contract": ""
        },
        {
            "name": "open",
            "type": "open",
            "ricardian_contract": ""
        },
        {
            "name": "retire",
            "type": "retire",
            "ricardian_contract": ""
        },
        {
            "name": "transfer",
            "type": "transfer",
            "ricardian_contract": ""
        }
    ],
    "tables": [
        {
            "name": "accounts",
            "index_type": "i64",
            "key_names": [],
            "key_types": [],
            "type": "account"
        },
        {
            "name": "stat",
            "index_type": "i64",
            "key_names": [],
            "key_types": [],
            "type": "currency_stats"
        }
    ],
    "ricardian_clauses": [],
    "variants": [
        {
            "name": "mega",
            "types": ["uint64", "string"]
        }
    ]
}
`)

var transferData = []byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x28, 0x5d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xae, 0x39,
	0x10, 0x27, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x45, 0x4f, 0x53, 0x00, 0x00, 0x00, 0x00,
	0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
	// extra variant
	0x01,
	// utf8 string "foo"
	0x03, 0x66, 0x6f, 0x6f,
	// extra2 array
	0x01,                                           // 1 item
	0x80, 0xb1, 0x91, 0x5e, 0x5d, 0x26, 0x8d, 0xca, // name teamgreymass
}

func TestAbiDecode(t *testing.T) {
	rv, err := tokenAbi.Decode(bytes.NewReader(transferData), "megatransfer")
	assert.NoError(t, err)
	assert.Equal(t, rv, map[string]interface{}{
		"from":     chain.N("foo"),
		"to":       chain.N("bar"),
		"quantity": *chain.A("1.0000 EOS"),
		"memo":     "hello",
		"extra":    []interface{}{"string", "foo"},
		"extra2": []interface{}{
			map[string]interface{}{"moo": chain.N("teamgreymass")},
		},
	})
}

func TestAbiEncode(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	err := tokenAbi.Encode(buf, "megatransfer", map[string]interface{}{
		"from":     chain.N("foo"),
		"to":       chain.N("bar"),
		"quantity": *chain.A("1.0000 EOS"),
		"memo":     "hello",
		"extra":    []interface{}{"string", "foo"},
		"extra2": []interface{}{
			map[string]interface{}{"moo": chain.N("teamgreymass")},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), transferData)
}
