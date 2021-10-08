package chain

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/greymass/go-eosio/pkg/abi"
)

type Uint128 struct {
	Lo uint64
	Hi uint64
}

type Int128 Uint128

type Float128 struct {
	Data [16]byte // not sure what the best way to represent this is, abieos seems to use a hex string ¯\_(ツ)_/¯
}

// uint64 alias that encodes to string for values above 32bit instead of scientific notation in JSON
type Uint64 uint64

// Type representing a block number, EOSIO chains are only expected to live for 68 years, sorry kids!
type BlockNum uint32

func (bn BlockNum) String() string {
	return fmt.Sprintf("%010d", bn)
}

// Create new Uint128 from big.Int, panics if big.Int is too large.
func NewUint128(i *big.Int) Uint128 {
	var b [16]byte
	i.FillBytes(b[:])
	return Uint128{
		binary.BigEndian.Uint64(b[8:]),
		binary.BigEndian.Uint64(b[:8]),
	}
}

// Create new Int128 from big.Int, panics if big.Int is too large.
func NewInt128(i *big.Int) Int128 {
	var b []byte
	if i.BitLen() > 128 {
		panic("int128 only supports bitlen <= 128")
	}
	switch i.Sign() {
	case 0:
		return Int128{0, 0}
	case 1:
		b = i.Bytes()
		for len(b) < 16 {
			b = append([]byte{0}, b...)
		}
	case -1:
		length := uint(i.BitLen()/8+1) * 8
		b = new(big.Int).Add(i, new(big.Int).Lsh(big.NewInt(1), length)).Bytes()
		for len(b) < 16 {
			b = append([]byte{0xff}, b...)
		}
	}
	return Int128{
		binary.BigEndian.Uint64(b[8:]),
		binary.BigEndian.Uint64(b[:8]),
	}
}

func NewUint128FromString(s string) (Uint128, error) {
	i, ok := (&big.Int{}).SetString(s, 10)
	if !ok || i.BitLen() > 128 || i.Sign() == -1 {
		return Uint128{}, errors.New("invalid unsigned integer")
	}
	return NewUint128(i), nil
}

func NewInt128FromString(s string) (Int128, error) {
	i, ok := (&big.Int{}).SetString(s, 10)
	if !ok || i.BitLen() > 128 {
		return Int128{}, errors.New("invalid signed integer")
	}
	return NewInt128(i), nil
}

func (u128 Uint128) Bytes(o binary.ByteOrder) []byte {
	var b [16]byte
	if o == binary.BigEndian {
		o.PutUint64(b[:8], u128.Hi)
		o.PutUint64(b[8:], u128.Lo)
	} else {
		o.PutUint64(b[8:], u128.Hi)
		o.PutUint64(b[:8], u128.Lo)
	}
	return b[:]
}

func (i128 Int128) Bytes(o binary.ByteOrder) []byte {
	return Uint128(i128).Bytes(o)
}

func (u128 Uint128) BigInt() *big.Int {
	rv := big.NewInt(0)
	rv.SetBytes(u128.Bytes(binary.BigEndian))
	return rv
}

func (i128 Int128) BigInt() *big.Int {
	b := i128.Bytes(binary.BigEndian)
	rv := big.NewInt(0)
	rv.SetBytes(b)
	if len(b) > 0 && b[0]&0x80 > 0 {
		rv.Sub(rv, new(big.Int).Lsh(big.NewInt(1), uint(len(b))*8))
	}
	return rv
}

func (u128 Uint128) String() string {
	return u128.BigInt().String()
}

func (i128 Int128) String() string {
	return i128.BigInt().String()
}

// abi.Marshaler conformance

func (u128 Uint128) MarshalABI(e *abi.Encoder) error {
	return e.WriteBytes(u128.Bytes(binary.LittleEndian))
}

func (i128 Int128) MarshalABI(e *abi.Encoder) error {
	return e.WriteBytes(i128.Bytes(binary.LittleEndian))
}

func (f128 Float128) MarshalABI(e *abi.Encoder) error {
	return e.WriteBytes(f128.Data[:])
}

func (bn BlockNum) MarshalABI(e *abi.Encoder) error {
	return e.WriteUint32(uint32(bn))
}

// abi.Unmarshaler conformance

func (u128 *Uint128) UnmarshalABI(d *abi.Decoder) error {
	lo, err := d.ReadUint64()
	if err != nil {
		return err
	}
	hi, err := d.ReadUint64()
	if err == nil {
		*u128 = Uint128{Lo: lo, Hi: hi}
	}
	return err
}

func (i128 *Int128) UnmarshalABI(d *abi.Decoder) error {
	return (*Uint128)(i128).UnmarshalABI(d)
}

func (f128 *Float128) UnmarshalABI(d *abi.Decoder) error {
	_, b, err := d.ReadBytes(16)
	if err == nil {
		copy(f128.Data[:], b)
	}
	return err
}

func (bn *BlockNum) UnmarshalABI(d *abi.Decoder) error {
	v, err := d.ReadUint32()
	if err == nil {
		*bn = BlockNum(v)
	}
	return err
}

// encoding.TextMarshaler conformance

func (u128 Uint128) MarshalText() (text []byte, err error) {
	return []byte(u128.String()), nil
}

func (i128 Int128) MarshalText() (text []byte, err error) {
	return []byte(i128.String()), nil
}

func (f128 Float128) MarshalText() (text []byte, err error) {
	return []byte(hex.EncodeToString(f128.Data[:])), nil
}

// encoding.TextUnmarshaler conformance

func (u128 *Uint128) UnmarshalText(text []byte) error {
	var err error
	*u128, err = NewUint128FromString(string(text))
	return err
}

func (i128 *Int128) UnmarshalText(text []byte) error {
	var err error
	*i128, err = NewInt128FromString(string(text))
	return err
}

func (f128 *Float128) UnmarshalText(text []byte) error {
	b, err := hex.DecodeString(string(text))
	if err == nil {
		copy(f128.Data[:], b)
	}
	return err
}

// json.Marshaler conformance

func (u64 Uint64) MarshalJSON() ([]byte, error) {
	return writeUintJSON(uint64(u64)), nil
}

func (bn BlockNum) MarshalJSON() ([]byte, error) {
	return writeUintJSON(uint64(bn)), nil
}

// json.Unmarshaler conformance

func (u64 *Uint64) UnmarshalJSON(b []byte) error {
	v, err := readUintJSON(b)
	if err == nil {
		*u64 = Uint64(v)
	}
	return err
}

func (bn *BlockNum) UnmarshalJSON(b []byte) error {
	v, err := readUintJSON(b)
	if v > math.MaxUint32 {
		return fmt.Errorf("block number %d is too large", v)
	}
	if err == nil {
		*bn = BlockNum(v)
	}
	return err
}

// json helpers

func readUintJSON(b []byte) (uint64, error) {
	s := string(b)
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	return strconv.ParseUint(s, 10, 64)
}

func writeUintJSON(v uint64) []byte {
	s := strconv.FormatUint(v, 10)
	if v > math.MaxUint32 {
		s = `"` + s + `"`
	}
	return []byte(s)
}
