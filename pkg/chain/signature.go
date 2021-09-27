package chain

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/greymass/go-eosio/pkg/abi"
	"github.com/greymass/go-eosio/pkg/base58"
)

type Signature struct {
	Type KeyType
	Data []byte
}

func NewSignature(t KeyType, d []byte) *Signature {
	return &Signature{
		Type: t,
		Data: d,
	}
}

func NewSignatureString(s string) (*Signature, error) {
	if len(s) < 7 || s[:4] != "SIG_" {
		return nil, errors.New("invalid signature string")
	}
	var t KeyType
	switch s[4:6] {
	case "K1":
		t = K1
	case "P1":
		t = P1
	case "WA":
		t = WA
	default:
		return nil, fmt.Errorf("unknown signature type: %s", s[4:6])
	}
	d, err := base58.CheckDecodeEosio(s[7:], t.String())
	return &Signature{
		Type: t,
		Data: d,
	}, err

}

func (pk *Signature) String() string {
	return "SIG_" + pk.Type.String() + "_" + base58.CheckEncodeEosio(pk.Data, pk.Type.String())
}

// abi.Marshaler conformance

func (s Signature) MarshalABI(e *abi.Encoder) error {
	err := e.WriteByte(byte(s.Type))
	if err != nil {
		return err
	}
	return e.WriteBytes(s.Data)
}

// abi.Unmarshaler conformance

func (s *Signature) UnmarshalABI(d *abi.Decoder) error {
	t, err := d.ReadByte()
	if err != nil {
		return err
	}
	s.Type = KeyType(t)
	switch s.Type {
	case K1, P1:
		_, s.Data, err = d.ReadBytes(65)
	case WA:
		_, data, err := d.ReadBytes(65) // sig data
		if err != nil {
			return err
		}
		al, err := d.ReadVaruint32() // auth_data len
		if err != nil {
			return err
		}
		_, ad, err := d.ReadBytes(int(al)) // auth_data
		if err != nil {
			return err
		}
		cl, err := d.ReadVaruint32() // client_json len
		if err != nil {
			return err
		}
		_, cd, err := d.ReadBytes(int(cl)) // client_json
		if err != nil {
			return err
		}
		tmp := make([]byte, binary.MaxVarintLen32)
		n := binary.PutUvarint(tmp, uint64(al))
		data = append(data, tmp[:n]...)
		data = append(data, ad...)
		n = binary.PutUvarint(tmp, uint64(cl))
		data = append(data, tmp[:n]...)
		data = append(data, cd...)
		s.Data = data
	default:
		return fmt.Errorf("unknown key type: %d", t)
	}

	return err
}

// encoding.TextMarshaler conformance

func (s *Signature) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

// encoding.TextUnmarshaler conformance

func (s *Signature) UnmarshalText(text []byte) error {
	new, err := NewSignatureString(string(text))
	if err == nil {
		*s = *new
	}
	return err
}
