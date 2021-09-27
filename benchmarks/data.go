package benchmarks

import (
	"github.com/greymass/go-eosio/pkg/chain"

	eoscanada "github.com/eoscanada/eos-go"
)

var testTransaction = chain.Transaction{
	TransactionHeader: chain.TransactionHeader{
		Expiration:       1519257492,
		RefBlockNum:      1,
		RefBlockPrefix:   2,
		MaxNetUsageWords: 3,
		MaxCpuUsageMs:    4,
		DelaySec:         5,
	},
	ContextFreeActions: []chain.Action{
		{
			Account: 7338027470446133248,
			Name:    11323884548116185088,
			Authorization: []chain.PermissionLevel{
				{Actor: 7241788200811757568, Permission: 3617214756542218240},
			},
			Data: []byte{
				0xbe, 0xef,
			},
		},
	},
	Actions: []chain.Action{
		{
			Account: 7338027470446133248,
			Name:    14829575318809870336,
			Authorization: []chain.PermissionLevel{
				{Actor: 3580161049631916032, Permission: 12224602554813644800},
				{Actor: 5307181291334008832, Permission: 12224602554813644800},
				{Actor: 9014633609356640256, Permission: 12224602554813644800},
				{Actor: 10500712387475144704, Permission: 12224602554813644800},
				{Actor: 10926484768698662912, Permission: 12224602554813644800},
				{Actor: 13990885797717344256, Permission: 12224602554813644800},
			},
			Data: []byte{
				0xde, 0xad, 0xc0, 0xde,
			},
		},
		{
			Account: 7338027470446133248,
			Name:    6292810045348380672,
			Authorization: []chain.PermissionLevel{
				{Actor: 14595364149838066048, Permission: 3617214756542218240},
			},
			Data: []byte{
				0xba, 0xbe, 0x92, 0x3c, 0x59, 0xd5, 0x14, 0x5b, 0xc3, 0x13, 0x03, 0x93, 0x35, 0xf5, 0x9f, 0x3b,
				0xc7, 0x55, 0xfd, 0xe1, 0xde, 0xaf, 0xa1, 0x0e, 0x62, 0x43, 0xff, 0xf4, 0x23, 0x46, 0xbe, 0xb4,
				0xb1, 0xe7, 0x81, 0x88, 0x5f, 0x1b, 0x6c, 0x82, 0x42, 0x60, 0x79, 0xcc, 0xb2, 0x7d, 0x9e, 0x74,
				0x2f, 0x3f, 0x7f, 0x4e, 0x1b, 0x7b, 0xd0, 0xb9, 0x50, 0x82, 0x6d, 0x44, 0x3b, 0x50, 0xc2, 0xe6,
				0xde, 0x34, 0xc0, 0x84, 0x6f, 0xcd, 0x84, 0xfa, 0x73, 0x6e, 0x70, 0x0e, 0xc5, 0x0b, 0x6b, 0xce,
				0xbf, 0x36, 0x75, 0x41, 0x1d, 0x45, 0x48, 0x26, 0x07, 0xe1, 0x92, 0x2b, 0xcf, 0x8f, 0x9d, 0xf8,
				0x5b, 0xc9, 0x8c, 0xb7, 0x1e, 0xcf, 0xa1, 0x67, 0x05, 0x36, 0xe3, 0x34, 0x0e, 0xd9, 0xc5, 0x9a,
				0xe9, 0x54, 0xa6, 0x91, 0x6d, 0xed, 0x90, 0xa9, 0xe7, 0x88, 0x1e, 0xf1, 0xbb, 0x41, 0x1c, 0x05,
			},
		},
	},
}

var testTransactionNoOptimize = Transaction{
	TransactionHeader: TransactionHeader{
		Expiration:       uint32(testTransaction.Expiration),
		RefBlockNum:      testTransaction.RefBlockNum,
		RefBlockPrefix:   testTransaction.RefBlockPrefix,
		MaxNetUsageWords: testTransaction.MaxNetUsageWords,
		MaxCpuUsageMs:    testTransaction.MaxCpuUsageMs,
		DelaySec:         testTransaction.DelaySec,
	},
	ContextFreeActions: []Action{
		{
			Account: uint64(testTransaction.ContextFreeActions[0].Account),
			Name:    uint64(testTransaction.ContextFreeActions[0].Name),
			Authorization: []PermissionLevel{
				{
					Actor:      uint64(testTransaction.ContextFreeActions[0].Authorization[0].Actor),
					Permission: uint64(testTransaction.ContextFreeActions[0].Authorization[0].Permission),
				},
			},
			Data: testTransaction.ContextFreeActions[0].Data,
		},
	},
	Actions: []Action{
		{
			Account: uint64(testTransaction.Actions[0].Account),
			Name:    uint64(testTransaction.Actions[0].Name),
			Authorization: []PermissionLevel{
				{
					Actor:      uint64(testTransaction.Actions[0].Authorization[0].Actor),
					Permission: uint64(testTransaction.Actions[0].Authorization[0].Permission),
				},
				{
					Actor:      uint64(testTransaction.Actions[0].Authorization[1].Actor),
					Permission: uint64(testTransaction.Actions[0].Authorization[1].Permission),
				},
				{
					Actor:      uint64(testTransaction.Actions[0].Authorization[2].Actor),
					Permission: uint64(testTransaction.Actions[0].Authorization[2].Permission),
				},
				{
					Actor:      uint64(testTransaction.Actions[0].Authorization[3].Actor),
					Permission: uint64(testTransaction.Actions[0].Authorization[3].Permission),
				},
				{
					Actor:      uint64(testTransaction.Actions[0].Authorization[4].Actor),
					Permission: uint64(testTransaction.Actions[0].Authorization[4].Permission),
				},
				{
					Actor:      uint64(testTransaction.Actions[0].Authorization[5].Actor),
					Permission: uint64(testTransaction.Actions[0].Authorization[5].Permission),
				},
			},
			Data: testTransaction.Actions[0].Data,
		},
		{
			Account: uint64(testTransaction.Actions[1].Account),
			Name:    uint64(testTransaction.Actions[1].Name),
			Authorization: []PermissionLevel{
				{
					Actor:      uint64(testTransaction.Actions[1].Authorization[0].Actor),
					Permission: uint64(testTransaction.Actions[1].Authorization[0].Permission),
				},
			},
			Data: testTransaction.Actions[1].Data,
		},
	},
	Extensions: []TransactionExtension{},
}

var testTransactionCanada = eoscanada.Transaction{
	TransactionHeader: eoscanada.TransactionHeader{
		Expiration: eoscanada.JSONTime{
			Time: chain.TimePointSec(1519257492).Time(),
		},
		RefBlockNum:      1,
		RefBlockPrefix:   2,
		MaxNetUsageWords: 3,
		MaxCPUUsageMS:    4,
		DelaySec:         5,
	},
	ContextFreeActions: []*eoscanada.Action{
		{
			Account: "greymass",
			Name:    "nonce",
			Authorization: []eoscanada.PermissionLevel{
				{Actor: "gm", Permission: "active"},
			},
			ActionData: eoscanada.NewActionDataFromHexData([]byte{
				0xbe, 0xef,
			}),
		},
	},
	Actions: []*eoscanada.Action{
		{
			Account: "greymass",
			Name:    "transform",
			Authorization: []eoscanada.PermissionLevel{
				{Actor: "aaron.gm", Permission: "pancake"},
				{Actor: "daniel.gm", Permission: "pancake"},
				{Actor: "johan.gm", Permission: "pancake"},
				{Actor: "max.gm", Permission: "pancake"},
				{Actor: "myles.gm", Permission: "pancake"},
				{Actor: "scott.gm", Permission: "pancake"},
			},
			ActionData: eoscanada.NewActionDataFromHexData([]byte{
				0xde, 0xad, 0xc0, 0xde,
			}),
		},
		{
			Account: "greymass",
			Name:    "execute",
			Authorization: []eoscanada.PermissionLevel{
				{Actor: "teamgreymass", Permission: "active"},
			},
			ActionData: eoscanada.NewActionDataFromHexData([]byte{
				0xba, 0xbe, 0x92, 0x3c, 0x59, 0xd5, 0x14, 0x5b, 0xc3, 0x13, 0x03, 0x93, 0x35, 0xf5, 0x9f, 0x3b,
				0xc7, 0x55, 0xfd, 0xe1, 0xde, 0xaf, 0xa1, 0x0e, 0x62, 0x43, 0xff, 0xf4, 0x23, 0x46, 0xbe, 0xb4,
				0xb1, 0xe7, 0x81, 0x88, 0x5f, 0x1b, 0x6c, 0x82, 0x42, 0x60, 0x79, 0xcc, 0xb2, 0x7d, 0x9e, 0x74,
				0x2f, 0x3f, 0x7f, 0x4e, 0x1b, 0x7b, 0xd0, 0xb9, 0x50, 0x82, 0x6d, 0x44, 0x3b, 0x50, 0xc2, 0xe6,
				0xde, 0x34, 0xc0, 0x84, 0x6f, 0xcd, 0x84, 0xfa, 0x73, 0x6e, 0x70, 0x0e, 0xc5, 0x0b, 0x6b, 0xce,
				0xbf, 0x36, 0x75, 0x41, 0x1d, 0x45, 0x48, 0x26, 0x07, 0xe1, 0x92, 0x2b, 0xcf, 0x8f, 0x9d, 0xf8,
				0x5b, 0xc9, 0x8c, 0xb7, 0x1e, 0xcf, 0xa1, 0x67, 0x05, 0x36, 0xe3, 0x34, 0x0e, 0xd9, 0xc5, 0x9a,
				0xe9, 0x54, 0xa6, 0x91, 0x6d, 0xed, 0x90, 0xa9, 0xe7, 0x88, 0x1e, 0xf1, 0xbb, 0x41, 0x1c, 0x05,
			}),
		},
	},
}

var testTransactionData = []byte{
	0x94, 0x07, 0x8e, 0x5a, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x04, 0x05, 0x01, 0x00, 0x00,
	0x00, 0x18, 0x1b, 0xe9, 0xd5, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x85, 0x26, 0x9d, 0x01, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x64, 0x00, 0x00, 0x00, 0x00, 0xa8, 0xed, 0x32, 0x32, 0x02,
	0xbe, 0xef, 0x02, 0x00, 0x00, 0x00, 0x18, 0x1b, 0xe9, 0xd5, 0x65, 0x00, 0x00, 0x90, 0x97, 0x2e,
	0x3c, 0xcd, 0xcd, 0x06, 0x00, 0x00, 0x00, 0x92, 0x81, 0x49, 0xaf, 0x31, 0x00, 0x00, 0x00, 0x40,
	0x41, 0x83, 0xa6, 0xa9, 0x00, 0x00, 0x90, 0x0c, 0x44, 0xe5, 0xa6, 0x49, 0x00, 0x00, 0x00, 0x40,
	0x41, 0x83, 0xa6, 0xa9, 0x00, 0x00, 0x00, 0x92, 0x81, 0x69, 0x1a, 0x7d, 0x00, 0x00, 0x00, 0x40,
	0x41, 0x83, 0xa6, 0xa9, 0x00, 0x00, 0x00, 0x00, 0x48, 0x06, 0xba, 0x91, 0x00, 0x00, 0x00, 0x40,
	0x41, 0x83, 0xa6, 0xa9, 0x00, 0x00, 0x00, 0x92, 0x01, 0xac, 0xa2, 0x97, 0x00, 0x00, 0x00, 0x40,
	0x41, 0x83, 0xa6, 0xa9, 0x00, 0x00, 0x00, 0x92, 0x81, 0x9c, 0x29, 0xc2, 0x00, 0x00, 0x00, 0x40,
	0x41, 0x83, 0xa6, 0xa9, 0x04, 0xde, 0xad, 0xc0, 0xde, 0x00, 0x00, 0x00, 0x18, 0x1b, 0xe9, 0xd5,
	0x65, 0x00, 0x00, 0x00, 0x40, 0x65, 0x8d, 0x54, 0x57, 0x01, 0x80, 0xb1, 0x91, 0x5e, 0x5d, 0x26,
	0x8d, 0xca, 0x00, 0x00, 0x00, 0x00, 0xa8, 0xed, 0x32, 0x32, 0x80, 0x01, 0xba, 0xbe, 0x92, 0x3c,
	0x59, 0xd5, 0x14, 0x5b, 0xc3, 0x13, 0x03, 0x93, 0x35, 0xf5, 0x9f, 0x3b, 0xc7, 0x55, 0xfd, 0xe1,
	0xde, 0xaf, 0xa1, 0x0e, 0x62, 0x43, 0xff, 0xf4, 0x23, 0x46, 0xbe, 0xb4, 0xb1, 0xe7, 0x81, 0x88,
	0x5f, 0x1b, 0x6c, 0x82, 0x42, 0x60, 0x79, 0xcc, 0xb2, 0x7d, 0x9e, 0x74, 0x2f, 0x3f, 0x7f, 0x4e,
	0x1b, 0x7b, 0xd0, 0xb9, 0x50, 0x82, 0x6d, 0x44, 0x3b, 0x50, 0xc2, 0xe6, 0xde, 0x34, 0xc0, 0x84,
	0x6f, 0xcd, 0x84, 0xfa, 0x73, 0x6e, 0x70, 0x0e, 0xc5, 0x0b, 0x6b, 0xce, 0xbf, 0x36, 0x75, 0x41,
	0x1d, 0x45, 0x48, 0x26, 0x07, 0xe1, 0x92, 0x2b, 0xcf, 0x8f, 0x9d, 0xf8, 0x5b, 0xc9, 0x8c, 0xb7,
	0x1e, 0xcf, 0xa1, 0x67, 0x05, 0x36, 0xe3, 0x34, 0x0e, 0xd9, 0xc5, 0x9a, 0xe9, 0x54, 0xa6, 0x91,
	0x6d, 0xed, 0x90, 0xa9, 0xe7, 0x88, 0x1e, 0xf1, 0xbb, 0x41, 0x1c, 0x05, 0x00,
}

var testTransferData = []byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x28, 0x5d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xae, 0x39,
	0x10, 0x27, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x45, 0x4f, 0x53, 0x00, 0x00, 0x00, 0x00,
	0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
}

var transferAbiJson = `
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
            "name": "close",
            "base": "",
            "fields": [
                {
                    "name": "owner",
                    "type": "name"
                },
                {
                    "name": "symbol",
                    "type": "symbol"
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
                    "type": "extra[]?"
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
`
