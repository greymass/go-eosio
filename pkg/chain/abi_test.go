package chain_test

import (
	"bytes"
	"encoding/hex"
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

var noopAbi = loadAbi(`{
	"version": "eosio::abi/1.1",
	"structs": [{"name": "noop", "base": "", "fields": []}],
	"actions": [{"name": "noop", "type": "noop", "ricardian_contract": ""}]
}`)

func TestAbiDecodeZeroFieldStruct(t *testing.T) {
	rv, err := noopAbi.DecodeAction(bytes.NewReader([]byte{}), chain.N("noop"))
	assert.NoError(t, err)
	assert.Equal(t, rv, map[string]interface{}{})
}

func TestAbiEncodeZeroFieldStruct(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	err := noopAbi.EncodeAction(buf, chain.N("noop"), map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, len(buf.Bytes()), 0)
}

// ABI 1.0 binary format does not include the Variants field.
// This tests that we can decode such ABIs without error.
// Uses real ABI data from EOS mainnet block 128 (eosio system contract).
func TestAbiDecodeVersion10(t *testing.T) {
	// Real ABI 1.0 data from EOS block 128, setabi action for eosio account
	// Fetched via: curl https://eos.greymass.com/v1/chain/get_block -d '{"block_num_or_id": 128}'
	abi10Hex := "0e656f73696f3a3a6162692f312e30050c6163636f756e745f6e616d65046e616d650f7065726d697373696f6e5f6e616d65046e616d650b616374696f6e5f6e616d65046e616d65137472616e73616374696f6e5f69645f747970650b636865636b73756d3235360b7765696768745f747970650675696e74313614107065726d697373696f6e5f6c6576656c0002056163746f720c6163636f756e745f6e616d650a7065726d697373696f6e0f7065726d697373696f6e5f6e616d650a6b65795f7765696768740002036b65790a7075626c69635f6b6579067765696768740b7765696768745f74797065177065726d697373696f6e5f6c6576656c5f77656967687400020a7065726d697373696f6e107065726d697373696f6e5f6c6576656c067765696768740b7765696768745f747970650b776169745f776569676874000208776169745f7365630675696e743332067765696768740b7765696768745f7479706509617574686f726974790004097468726573686f6c640675696e743332046b6579730c6b65795f7765696768745b5d086163636f756e7473197065726d697373696f6e5f6c6576656c5f7765696768745b5d0577616974730d776169745f7765696768745b5d0a6e65776163636f756e7400040763726561746f720c6163636f756e745f6e616d65046e616d650c6163636f756e745f6e616d65056f776e657209617574686f726974790661637469766509617574686f7269747907736574636f64650004076163636f756e740c6163636f756e745f6e616d6506766d747970650575696e743809766d76657273696f6e0575696e743804636f6465056279746573067365746162690002076163636f756e740c6163636f756e745f6e616d65036162690562797465730a757064617465617574680004076163636f756e740c6163636f756e745f6e616d650a7065726d697373696f6e0f7065726d697373696f6e5f6e616d6506706172656e740f7065726d697373696f6e5f6e616d65046175746809617574686f726974790a64656c657465617574680002076163636f756e740c6163636f756e745f6e616d650a7065726d697373696f6e0f7065726d697373696f6e5f6e616d65086c696e6b617574680004076163636f756e740c6163636f756e745f6e616d6504636f64650c6163636f756e745f6e616d6504747970650b616374696f6e5f6e616d650b726571756972656d656e740f7065726d697373696f6e5f6e616d650a756e6c696e6b617574680003076163636f756e740c6163636f756e745f6e616d6504636f64650c6163636f756e745f6e616d6504747970650b616374696f6e5f6e616d650b63616e63656c64656c617900020e63616e63656c696e675f61757468107065726d697373696f6e5f6c6576656c067472785f6964137472616e73616374696f6e5f69645f74797065076f6e6572726f7200020973656e6465725f69640775696e743132380873656e745f747278056279746573127365745f6163636f756e745f6c696d6974730004076163636f756e740c6163636f756e745f6e616d650972616d5f627974657305696e7436340a6e65745f77656967687405696e7436340a6370755f77656967687405696e74363407736574707269760002076163636f756e740c6163636f756e745f6e616d650769735f7072697604696e7438117365745f676c6f62616c5f6c696d6974730001136370755f757365635f7065725f706572696f6405696e7436340c70726f64756365725f6b657900020d70726f64756365725f6e616d650c6163636f756e745f6e616d6511626c6f636b5f7369676e696e675f6b65790a7075626c69635f6b65790d7365745f70726f6475636572730001087363686564756c650e70726f64756365725f6b65795b5d0c726571756972655f6175746800010466726f6d0c6163636f756e745f6e616d650e00409e9a2264b89a0a6e65776163636f756e740000000040258ab2c207736574636f64650000000000b863b2c206736574616269000040cbdaa86c52d50a75706461746561757468000040cbdaa8aca24a0a64656c65746561757468000000002d6b03a78b086c696e6b61757468000040cbdac0e9e2d40a756e6c696e6b617574680000bc892a4585a6410b63616e63656c64656c617900000000e0d27bd5a4076f6e6572726f72000000ce4eba68b2c2127365745f6163636f756e745f6c696d697473000000ce4ebac8b2c2117365745f676c6f62616c5f6c696d6974730000000060bb5bb3c207736574707269760000000038d15bb3c20d7365745f70726f64756365727300000000a0656dacba0c726571756972655f617574680000000000"

	abi10Data, err := hex.DecodeString(abi10Hex)
	assert.NoError(t, err)

	reader := bytes.NewReader(abi10Data)
	dec := chain.NewDecoder(reader)
	var decoded chain.Abi
	err = dec.Decode(&decoded)

	assert.NoError(t, err)
	assert.Equal(t, decoded.Version, "eosio::abi/1.0")
	assert.Equal(t, len(decoded.Types), 5)
	assert.Equal(t, len(decoded.Structs), 20)
	assert.Equal(t, len(decoded.Actions), 14)
	assert.Equal(t, len(decoded.Variants), 0)
}

var extensionAbi = loadAbi(`{
	"version": "eosio::abi/1.1",
	"structs": [{
		"name": "setabi",
		"base": "",
		"fields": [
			{"name": "account", "type": "name"},
			{"name": "abi", "type": "bytes"},
			{"name": "memo", "type": "string$"}
		]
	}],
	"actions": [{"name": "setabi", "type": "setabi", "ricardian_contract": ""}]
}`)

func TestAbiDecodeExtensionFieldAbsent(t *testing.T) {
	hexData := "0000000000ea305504deadbeef"
	data, err := hex.DecodeString(hexData)
	assert.NoError(t, err)

	reader := bytes.NewReader(data)
	decoded, err := extensionAbi.Decode(reader, "setabi")

	assert.NoError(t, err)
	result := decoded.(map[string]interface{})
	assert.Equal(t, result["account"], chain.N("eosio"))
	assert.Equal(t, result["abi"], chain.Bytes{0xde, 0xad, 0xbe, 0xef})
	_, hasMemo := result["memo"]
	if hasMemo && result["memo"] != nil && result["memo"] != "" {
		t.Errorf("expected memo to be absent or empty, got %v", result["memo"])
	}
}

func TestAbiDecodeExtensionFieldPresent(t *testing.T) {
	hexData := "0000000000ea305504deadbeef0568656c6c6f"
	data, err := hex.DecodeString(hexData)
	assert.NoError(t, err)

	reader := bytes.NewReader(data)
	decoded, err := extensionAbi.Decode(reader, "setabi")

	assert.NoError(t, err)
	result := decoded.(map[string]interface{})
	assert.Equal(t, result["account"], chain.N("eosio"))
	assert.Equal(t, result["abi"], chain.Bytes{0xde, 0xad, 0xbe, 0xef})
	assert.Equal(t, result["memo"], "hello")
}
