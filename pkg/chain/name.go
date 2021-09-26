package chain

import "github.com/greymass/go-eosio/pkg/abi"

// Type representing an EOSIO name.
type Name uint64

// Create a new name from a string, e.g. "teamgreymass".
func NewName(s string) Name {
	return Name(stringToName(s))
}

// Convenience for NewName.
func N(s string) Name {
	return Name(stringToName(s))
}

// Returns string representation of the name.
func (n Name) String() string {
	return nameToString(uint64(n))
}

// abi.Unmarshaler conformance

func (n *Name) UnmarshalABI(d *abi.Decoder) error {
	v, err := d.ReadUint64()
	if err == nil {
		*n = Name(v)
	}
	return err
}

// encoding.TextMarshaler conformance

func (n Name) MarshalText() (text []byte, err error) {
	return []byte(n.String()), nil
}

// encoding.TextUnmarshaler conformance

func (n *Name) UnmarshalText(text []byte) error {
	new := NewName(string(text))
	*n = new
	return nil
}

// from https://github.com/eoscanada/eos-go/blob/21697b8969f6446181086db27ae7e7a302bf166d/name.go

var base32Alphabet = []byte(".12345abcdefghijklmnopqrstuvwxyz")

func nameToString(value uint64) string {
	a := []byte{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.'}
	tmp := value
	i := uint32(0)
	for ; i <= 12; i++ {
		bit := 0x1f
		if i == 0 {
			bit = 0x0f
		}
		c := base32Alphabet[tmp&uint64(bit)]
		a[12-i] = c

		shift := uint(5)
		if i == 0 {
			shift = 4
		}
		tmp >>= shift
	}
	return trimRightDots(a)
}

func trimRightDots(bytes []byte) string {
	trimUpTo := -1
	for i := 12; i >= 0; i-- {
		if bytes[i] == '.' {
			trimUpTo = i
		} else {
			break
		}
	}
	if trimUpTo == -1 {
		return string(bytes)
	}
	return string(bytes[0:trimUpTo])
}

func stringToName(s string) uint64 {
	var rv uint64 = 0
	var i uint32
	sLen := uint32(len(s))
	for ; i <= 12; i++ {
		var c uint64
		if i < sLen {
			c = uint64(charToSymbol(s[i]))
		}
		if i < 12 {
			c &= 0x1f
			c <<= 64 - 5*(i+1)
		} else {
			c &= 0x0f
		}
		rv |= c
	}
	return rv
}

func charToSymbol(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 'a' + 6
	}
	if c >= '1' && c <= '5' {
		return c - '1' + 1
	}
	return 0
}
