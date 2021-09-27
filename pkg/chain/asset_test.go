package chain_test

import (
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/chain"
)

func TestAsset(t *testing.T) {
	asset, err := chain.NewAssetFromString("0.42 BONK")
	assert.NoError(t, err)
	assert.Equal(t, asset.Value, int64(42))
	assert.Equal(t, asset.Symbol, chain.Symbol(323436364290))
	assert.Equal(t, "BONK", asset.Symbol.Name())
	assert.Equal(t, 2, asset.Symbol.Decimals())
	assert.Equal(t, "0.42 BONK", asset.String())
	assert.Equal(t, "0.42 BONK", chain.NewAsset(42, chain.Symbol(323436364290)).String())

	asset2, err := chain.NewAssetFromString("-1.234567890 TEST")
	assert.NoError(t, err)
	assert.Equal(t, "-1.234567890 TEST", asset2.String())

	assert.JSONCoding(t, asset, `"0.42 BONK"`)
	assert.ABICoding(t, asset, []byte{0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x42, 0x4f, 0x4e, 0x4b, 0x00, 0x00, 0x00})

	_, err = chain.NewAssetFromString("1 muu")
	assert.NotNil(t, err)
	_, err = chain.NewAssetFromString(" BAR")
	assert.NotNil(t, err)
	_, err = chain.NewAssetFromString("1")
	assert.NotNil(t, err)

	ex := chain.ExtendedAsset{
		Quantity: *asset,
		Contract: chain.N("scout"),
	}
	assert.JSONCoding(t, &ex, `{"quantity":"0.42 BONK","contract":"scout"}`)
	assert.ABICoding(t, &ex, []byte{
		0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x42, 0x4f, 0x4e, 0x4b, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x80, 0xac, 0x29, 0xc2,
	})

	asset3, err := chain.NewAssetFromString("1.0000 EOS")
	assert.NoError(t, err)
	assert.Equal(t, asset3.Value, int64(10000))
	assert.Equal(t, asset3.Decimals(), 4)
	assert.Equal(t, asset3.Precision(), 10000)
	assert.Equal(t, asset3.String(), "1.0000 EOS")
}

func BenchmarkAssetFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chain.NewAssetFromString("1.0042 EOS")
	}
}

func BenchmarkAssetToString(b *testing.B) {
	a, _ := chain.NewAssetFromString("1.0042 EOS")
	for i := 0; i < b.N; i++ {
		_ = a.String()
	}
}
