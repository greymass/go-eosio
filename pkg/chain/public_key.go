package chain

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/greymass/go-eosio/pkg/abi"
	"github.com/greymass/go-eosio/pkg/base58"
)

type PublicKey struct {
	Type KeyType
	Data []byte
}

func NewPublicKey(t KeyType, d []byte) *PublicKey {
	return &PublicKey{
		Type: t,
		Data: d,
	}
}

func NewPublicKeyFromString(s string) (*PublicKey, error) {
	if len(s) < 7 {
		return nil, errors.New("invalid public key string")
	}
	if s[0:4] == "PUB_" {
		// new format
		var t KeyType
		switch s[4:6] {
		case "K1":
			t = K1
		case "P1":
			t = P1
		case "WA":
			t = WA
		default:
			return nil, fmt.Errorf("unknown key type: %s", s[4:6])
		}
		d, err := base58.CheckDecodeEosio(s[7:], t.String())
		return &PublicKey{
			Type: t,
			Data: d,
		}, err
	}
	// legacy format
	d, err := base58.CheckDecode(s[len(s)-50:])
	return &PublicKey{
		Type: K1,
		Data: d,
	}, err
}

func (pk *PublicKey) String() string {
	return "PUB_" + pk.Type.String() + "_" + base58.CheckEncodeEosio(pk.Data, pk.Type.String())
}

// panics if key type isn't k1
func (pk *PublicKey) LegacyString(prefix string) string {
	if pk.Type != K1 {
		panic("only K1 keys can be converted to legacy format")
	}
	return prefix + base58.CheckEncode(pk.Data)
}

// abi.Unmarshaler conformance

func (pk *PublicKey) UnmarshalABI(d *abi.Decoder) error {
	t, err := d.ReadByte()
	if err != nil {
		return err
	}
	pk.Type = KeyType(t)
	switch pk.Type {
	case K1, P1:
		_, pk.Data, err = d.ReadBytes(33)
	case WA:
		_, data, err := d.ReadBytes(34) // key_data + user_presence
		if err != nil {
			return err
		}
		l, err := d.ReadVaruint32() // rpid length
		if err != nil {
			return err
		}
		_, rpid, err := d.ReadBytes(int(l))
		if err != nil {
			return err
		}
		tmp := make([]byte, binary.MaxVarintLen32)
		n := binary.PutUvarint(tmp, uint64(l))
		data = append(data, tmp[:n]...)
		data = append(data, rpid...)
		pk.Data = data
	default:
		return fmt.Errorf("unknown key type: %d", t)
	}

	return err
}

// encoding.TextMarshaler conformance

func (pk *PublicKey) MarshalText() (text []byte, err error) {
	return []byte(pk.String()), nil
}

// encoding.TextUnmarshaler conformance

func (pk *PublicKey) UnmarshalText(text []byte) error {
	new, err := NewPublicKeyFromString(string(text))
	if err == nil {
		*pk = *new
	}
	return err
}
