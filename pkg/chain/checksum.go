package chain

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	//lint:ignore SA1019 ripemd160 is deprecated because google knows better than you, and it will be slow to punish you for your sins
	"golang.org/x/crypto/ripemd160"

	"github.com/greymass/go-eosio/pkg/abi"
)

// ripemd160 checksum type
type Checksum160 [20]byte

// sha256 checksum type
type Checksum256 [32]byte

// sha512 checksum type
type Checksum512 [64]byte

// Return ripemd160 hash of message
func Checksum160Digest(message []byte) Checksum160 {
	h := ripemd160.New()
	h.Write(message)
	d := h.Sum(nil)
	var rv Checksum160
	copy(rv[:], d[:20])
	return rv
}

// Return sha256 hash of message
func Checksum256Digest(message []byte) Checksum256 {
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	var rv Checksum256
	copy(rv[:], d[:32])
	return rv
}

// Return sha512 hash of message
func Checksum512Digest(message []byte) Checksum512 {
	h := sha512.New()
	h.Write(message)
	d := h.Sum(nil)
	var rv Checksum512
	copy(rv[:], d[:64])
	return rv
}

func (c Checksum160) String() string {
	return hex.EncodeToString(c[:])
}

func (c Checksum256) String() string {
	return hex.EncodeToString(c[:])
}

func (c Checksum512) String() string {
	return hex.EncodeToString(c[:])
}

// abi.Unmarshaler conformance

func (c160 *Checksum160) UnmarshalABI(d *abi.Decoder) error {
	_, data, err := d.ReadBytes(20)
	if err == nil {
		copy((*c160)[:], data[:20])
	}
	return err
}

func (c256 *Checksum256) UnmarshalABI(d *abi.Decoder) error {
	_, data, err := d.ReadBytes(32)
	if err == nil {
		copy((*c256)[:], data[:32])
	}
	return err
}

func (c512 *Checksum512) UnmarshalABI(d *abi.Decoder) error {
	_, data, err := d.ReadBytes(64)
	if err == nil {
		copy((*c512)[:], data[:64])
	}
	return err
}

// encoding.TextMarshaler conformance

func (c160 Checksum160) MarshalText() (text []byte, err error) {
	return []byte(c160.String()), nil
}

func (c256 Checksum256) MarshalText() (text []byte, err error) {
	return []byte(c256.String()), nil
}

func (c512 Checksum512) MarshalText() (text []byte, err error) {
	return []byte(c512.String()), nil
}

// encoding.TextUnmarshaler conformance

func (c160 *Checksum160) UnmarshalText(text []byte) error {
	data := make([]byte, hex.DecodedLen(len(text)))
	_, err := hex.Decode(data, text)
	if err == nil {
		copy((*c160)[:], data[:20])
	}
	return err
}

func (c256 *Checksum256) UnmarshalText(text []byte) error {
	data := make([]byte, hex.DecodedLen(len(text)))
	_, err := hex.Decode(data, text)
	if err == nil {
		copy((*c256)[:], data[:32])
	}
	return err
}

func (c512 *Checksum512) UnmarshalText(text []byte) error {
	data := make([]byte, hex.DecodedLen(len(text)))
	_, err := hex.Decode(data, text)
	if err == nil {
		copy((*c512)[:], data[:64])
	}
	return err
}
