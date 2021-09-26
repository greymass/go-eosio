package chain_test

import (
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/chain"
)

func TestTime(t *testing.T) {
	assert.Equal(t, chain.TimePoint(0).String(), "1970-01-01T00:00:00.000")
	tp, err := chain.NewTimePointFromString("2018-02-21T23:58:00.000")
	assert.NoError(t, err)
	assert.Equal(t, tp, chain.TimePoint(1519257480000000))
	assert.Equal(t, tp.Time().Year(), 2018)
	assert.ABICoding(t, chain.TimePoint(1519257492222000), []byte{0x30, 0x70, 0x25, 0xb3, 0xc1, 0x65, 0x05, 0x00})
	assert.JSONCoding(t, chain.TimePoint(1519257492222000), `"2018-02-21T23:58:12.222"`)
	assert.ABICoding(t, chain.TimePointSec(1235), []byte{0xd3, 0x04, 0x00, 0x00})
	assert.JSONCoding(t, chain.TimePointSec(1235), `"1970-01-01T00:20:35"`)
	assert.ABICoding(t, chain.BlockTimestamp(1145145361), []byte{0x11, 0x88, 0x41, 0x44})
	assert.JSONCoding(t, chain.BlockTimestamp(1145145361), `"2018-02-21T23:58:00.500"`)
	assert.JSONCoding(t, chain.BlockTimestamp(1145145362), `"2018-02-21T23:58:01.000"`)
}
