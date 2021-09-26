package chain_test

import (
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/chain"
)

func TestSymbol(t *testing.T) {
	symbol, err := chain.NewSymbolFromString("4,BEZOS")
	assert.NoError(t, err)
	assert.Equal(t, symbol, chain.Symbol(91600282010116))
	assert.Equal(t, symbol.Name(), "BEZOS")
	assert.Equal(t, symbol.Decimals(), 4)
	assert.Equal(t, symbol.Precision(), 10000)
	assert.Equal(t, symbol.Code(), chain.SymbolCode(357813601602))

	_, err = chain.NewSymbolFromString("MO")
	assert.HasError(t, &err)
	_, err = chain.NewSymbolFromString("??")
	assert.HasError(t, &err)
	_, err = chain.NewSymbolFromString("-1,XYZ")
	assert.HasError(t, &err)
	_, err = chain.NewSymbolFromString("JERFERY,BERSOS")
	assert.HasError(t, &err)

	assert.ABICoding(t, symbol, []byte{0x04, 0x42, 0x45, 0x5a, 0x4f, 0x53, 0x00, 0x00})
	assert.JSONCoding(t, symbol, `"4,BEZOS"`)

	assert.ABICoding(t, symbol.Code(), []byte{0x42, 0x45, 0x5a, 0x4f, 0x53, 0x00, 0x00, 0x00})
	assert.JSONCoding(t, symbol.Code(), `"BEZOS"`)
}
