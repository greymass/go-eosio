package chain

import (
	"encoding/base64"

	"github.com/greymass/go-eosio/pkg/abi"
)

// Exactly like Bytes but serializes as Base64 instead of Base16 (hex) when encoded as JSON
// Worth noting is that this version encodes to valid Base64 while the eosio variant does not
// see https://github.com/EOSIO/eos/issues/8161
type Blob []byte

func (b Blob) Base64() string {
	return base64.StdEncoding.EncodeToString(b)
}

// abi.Marshaler conformance

func (b Blob) MarshalABI(e *abi.Encoder) error {
	var err error
	l := uint32(len(b))
	err = e.WriteVaruint32(l)
	if err == nil {
		err = e.WriteBytes(b)
	}
	return err
}

// abi.Unmarshaler conformance

func (b *Blob) UnmarshalABI(d *abi.Decoder) error {
	l, err := d.ReadVaruint32()
	if err != nil {
		return err
	}
	_, data, err := d.ReadBytes(int(l))
	if err == nil {
		*b = data
	}
	return err
}

// encoding.TextMarshaler conformance

func (b Blob) MarshalText() (text []byte, err error) {
	return []byte(b.Base64()), nil
}

// encoding.TextUnmarshaler conformance

func (b *Blob) UnmarshalText(text []byte) error {
	// fix up the base64 encoding padding
	switch len(text) % 4 {
	case 2:
		text = append(text, "=="...)
	case 3:
		text = append(text, "="...)
	}
	var err error
	*b, err = base64.StdEncoding.DecodeString(string(text))
	return err
}
