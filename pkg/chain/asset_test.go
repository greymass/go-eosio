package chain_test

import (
	"strconv"
	"strings"
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

	asset4, err := chain.NewAssetFromString("34.0303 EOS")
	assert.NoError(t, err)
	assert.Equal(t, asset4.Value, int64(340303))
	assert.Equal(t, asset4.Decimals(), 4)
	assert.Equal(t, asset4.Precision(), 10000)
	assert.Equal(t, asset4.String(), "34.0303 EOS")

	asset5, err := chain.NewAssetFromString("1.123456789012345678 ABCDEFG")
	assert.NoError(t, err)
	assert.Equal(t, asset5.Value, int64(1123456789012345678))
	assert.Equal(t, asset5.Decimals(), 18)
	assert.Equal(t, asset5.Precision(), 1000000000000000000)
	assert.Equal(t, asset5.String(), "1.123456789012345678 ABCDEFG")
}

func FuzzAsset(f *testing.F) {
	f.Add("1.0000 EOS")
	f.Add("34.0303 EOS")
	f.Add("-1.234567890 TEST")
	f.Add("0.42 BONK")
	f.Add("1.0042 EOS")
	f.Add("1.00000042 LONG")
	f.Add("1.123456789012345678 ABCDEFG")
	f.Fuzz(func(t *testing.T, orig string) {
		asset, error := chain.NewAssetFromString(orig)
		if error != nil {
			return
		}
		t.Logf("%s: %q float=%f symbol=%s", orig, asset.String(), asset.FloatValue(), asset.Symbol.String())
		parts := strings.Split(orig, " ")
		if asset.Symbol.Name() != parts[1] {
			t.Errorf("symbol name mismatch, expected %s, got %s", parts[1], asset.Symbol.Name())
		}
		numParts := strings.Split(parts[0], ".")
		if len(numParts) == 2 {
			if len(numParts[1]) != asset.Symbol.Decimals() {
				t.Errorf("decimals mismatch, expected %d, got %d", asset.Symbol.Decimals(), len(numParts[1]))
			}
		} else if asset.Symbol.Decimals() != 0 {
			t.Errorf("decimals mismatch, expected 0 got %d", asset.Symbol.Decimals())
		}
		floatVal, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			t.Error(err)
		}
		if floatVal != asset.FloatValue() {
			t.Errorf("%.18f != %.18f", floatVal, asset.FloatValue())
		}
	})
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
