package chain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/greymass/go-eosio/pkg/abi"
)

type Symbol uint64

type SymbolCode uint64

// Create new symbol from precision and name.
func NewSymbol(precision uint8, name string) (Symbol, error) {
	v, err := rawSymbolValue(precision, name)
	if err != nil {
		return 0, err
	}
	return Symbol(v), nil
}

// Create new asset symbol from string, e.g. "4,EOS"
func NewSymbolFromString(s string) (Symbol, error) {
	p := strings.Split(s, ",")
	if len(p) != 2 {
		return 0, errors.New("invalid asset symbol string")
	}
	precision, err := strconv.ParseInt(p[0], 10, 8)
	if err != nil {
		return 0, err
	}
	symbolValue, err := rawSymbolValue(uint8(precision), p[1])
	if err != nil {
		return 0, err
	}

	return Symbol(symbolValue), nil
}

// Asset symbol name, e.g. "EOS"
func (s Symbol) Name() string {
	v := s
	v >>= 8
	var rv string
	for v > 0 {
		rv += string(byte(v))
		v >>= 8
	}
	return rv
}

// Asset symbol code.
func (s Symbol) Code() SymbolCode {
	return SymbolCode(s >> 8)
}

// Number of decimals in symbol, e.g. 4 for EOS
func (s Symbol) Decimals() int {
	return int(s & 0xFF)
}

// Precision of asset symbol, 10^Decimals.
func (s Symbol) Precision() int {
	var p10 int = 1
	p := s.Decimals()
	for p > 0 {
		p10 *= 10
		p--
	}
	return p10
}

// String representation of asset symbol, e.g. "4,EOS"
func (s Symbol) String() string {
	return fmt.Sprint(s.Decimals()) + "," + s.Name()
}

// abi.Unmarshaler conformance

func (s *Symbol) UnmarshalABI(d *abi.Decoder) error {
	v, err := d.ReadUint64()
	if err == nil {
		*s = Symbol(v)
	}
	return err
}

func (s *SymbolCode) UnmarshalABI(d *abi.Decoder) error {
	v, err := d.ReadUint64()
	if err == nil {
		*s = SymbolCode(v)
	}
	return err
}

// encoding.TextMarshaler conformance

func (s Symbol) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

func (sc SymbolCode) MarshalText() (text []byte, err error) {
	s := Symbol(uint64(sc) << 8)
	return []byte(s.Name()), nil
}

// encoding.TextUnmarshaler conformance

func (s *Symbol) UnmarshalText(text []byte) error {
	new, err := NewSymbolFromString(string(text))
	if err == nil {
		*s = new
	}
	return err
}

func (sc *SymbolCode) UnmarshalText(text []byte) error {
	new, err := NewSymbolFromString("0," + string(text))
	if err == nil {
		*sc = new.Code()
	}
	return err
}

// helpers

func rawSymbolValue(precision uint8, s string) (uint64, error) {
	if precision > 18 {
		return 0, errors.New("invalid symbol precision")
	}
	var rv uint64 = 0
	for i := 0; i < len(s); i++ {
		if !(s[i] >= 'A' && s[i] <= 'Z') {
			return 0, errors.New("invalid character in symbol name")
		}
		rv |= (uint64(s[i]) << (8 * (i + 1)))
	}
	rv |= uint64(precision)
	return rv, nil
}
