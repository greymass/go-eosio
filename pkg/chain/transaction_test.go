package chain_test

import (
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/chain"
)

func TestTransaction(t *testing.T) {
	th := chain.TransactionHeader{
		Expiration:       chain.TimePointSec(1234567890),
		RefBlockNum:      11,
		RefBlockPrefix:   22,
		MaxNetUsageWords: 33,
		MaxCpuUsageMs:    44,
		DelaySec:         55,
	}
	assert.ABICoding(t, &th, []byte{0xd2, 0x02, 0x96, 0x49, 0x0b, 0x00, 0x16, 0x00, 0x00, 0x00, 0x21, 0x2c, 0x37})
	assert.JSONCoding(t, th, `
		{
			"expiration": "2009-02-13T23:31:30",
			"ref_block_num": 11,
			"ref_block_prefix": 22,
			"max_net_usage_words": 33,
			"max_cpu_usage_ms": 44,
			"delay_sec": 55
		}
	`)

	tx := chain.Transaction{
		TransactionHeader:  th,
		ContextFreeActions: []chain.Action{},
		Actions: []chain.Action{
			{
				Account: chain.N("foo"),
				Name:    chain.N("bar"),
				Authorization: []chain.PermissionLevel{
					{Actor: chain.N("baz"), Permission: chain.N("qux")},
					{Actor: chain.N("quux"), Permission: chain.N("quuz")},
				},
				Data: []byte{0xde, 0xad, 0xbe, 0xef},
			},
		},
		Extensions: []chain.TransactionExtension{},
	}
	assert.ABICoding(t, tx, []byte{
		0xd2, 0x02, 0x96, 0x49, 0x0b, 0x00, 0x16, 0x00, 0x00, 0x00, 0x21, 0x2c, 0x37, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x28, 0x5d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xae, 0x39, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xbe, 0x39, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xba, 0xb6,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xd0, 0xb5, 0xb6, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0xb5, 0xb6,
		0x04, 0xde, 0xad, 0xbe, 0xef, 0x00,
	})
	assert.JSONCoding(t, tx, `
		{
			"expiration": "2009-02-13T23:31:30",
			"ref_block_num": 11,
			"ref_block_prefix": 22,
			"max_net_usage_words": 33,
			"max_cpu_usage_ms": 44,
			"delay_sec": 55,
			"context_free_actions": [],
			"actions": [
				{
					"account": "foo",
					"name": "bar",
					"authorization": [
						{
							"actor": "baz",
							"permission": "qux"
						},
						{
							"actor": "quux",
							"permission": "quuz"
						}
					],
					"data": "deadbeef"
				}
			],
			"transaction_extensions": []
		}
	`)
}
