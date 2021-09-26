package chain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/greymass/go-eosio/pkg/abi"
)

type Asset struct {
	Value int64
	Symbol
}

type ExtendedAsset struct {
	Quantity Asset `json:"quantity"`
	Contract Name  `json:"contract"`
}

// Create new asset from value (number of symbol units) and symbol.
func NewAsset(value int64, symbol Symbol) *Asset {
	return &Asset{value, symbol}
}

// Create new asset from string, e.g. "1.0000 EOS"
func NewAssetFromString(s string) (*Asset, error) {
	p := strings.Split(s, " ")
	if len(p) != 2 {
		return nil, errors.New("invalid asset string")
	}
	vp := strings.Split(p[0], ".")
	var precision uint8 = 0
	if len(vp) == 2 {
		precision = uint8(len(vp[1]))
	}
	units, err := strconv.ParseInt(strings.Replace(p[0], ".", "", 1), 10, 64)
	if err != nil {
		return nil, err
	}

	symbol, err := NewSymbol(precision, p[1])

	if err != nil {
		return nil, err
	}

	return &Asset{units, symbol}, nil
}

// String representation of asset, e.g. "1.0000 EOS"
func (a *Asset) String() string {
	s := ""
	if a.Value < 0 {
		s = "-"
	}
	v := int(a.Value)
	if v < 0 {
		v = -v
	}
	s += fmt.Sprint(v / a.Symbol.Precision())
	if a.Symbol.Decimals() > 0 {
		f := v % a.Symbol.Precision()
		s += "." + fmt.Sprint(f)
	}
	return s + " " + a.Symbol.Name()
}

// abi.Unmarshaler conformance

func (a *Asset) UnmarshalABI(d *abi.Decoder) error {
	v, err := d.ReadInt64()
	if err != nil {
		return err
	}
	*a = Asset{Value: v, Symbol: 0}
	return a.Symbol.UnmarshalABI(d)
}

func (ea *ExtendedAsset) UnmarshalABI(d *abi.Decoder) error {
	var err error
	err = ea.Quantity.UnmarshalABI(d)
	if err != nil {
		return err
	}
	err = ea.Contract.UnmarshalABI(d)
	return err
}

// encoding.TextMarshaler conformance

func (a *Asset) MarshalText() (text []byte, err error) {
	return []byte(a.String()), nil
}

// encoding.TextUnmarshaler conformance

func (a *Asset) UnmarshalText(text []byte) error {
	new, err := NewAssetFromString(string(text))
	if err == nil {
		*a = *new
	}
	return err
}
