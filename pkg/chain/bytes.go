package chain

import (
	"encoding/hex"

	"github.com/greymass/go-eosio/pkg/abi"
)

// Binary data type, encodes to hex string in JSON.
// Use the Blob type instead where possible since it encodes to base64.
type Bytes []byte

func (b Bytes) Hex() string {
	return hex.EncodeToString(b)
}

// abi.Marshaler conformance

func (b Bytes) MarshalABI(e *abi.Encoder) error {
	var err error
	l := uint32(len(b))
	err = e.WriteVaruint32(l)
	if err == nil {
		err = e.WriteBytes(b)
	}
	return err
}

// abi.Unmarshaler conformance

func (b *Bytes) UnmarshalABI(d *abi.Decoder) error {
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

func (b Bytes) MarshalText() (text []byte, err error) {
	return []byte(b.Hex()), nil
}

// encoding.TextUnmarshaler conformance

func (b *Bytes) UnmarshalText(text []byte) error {
	data := make([]byte, hex.DecodedLen(len(text)))
	_, err := hex.Decode(data, text)
	if err == nil {
		*b = data
	}
	return nil
}
