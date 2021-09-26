package chain

import (
	"io"

	"github.com/greymass/go-eosio/pkg/abi"
)

func chainDecoder(dec *abi.Decoder, v interface{}) (done bool, err error) {
	// this is dumb but speeds up decoding a lot
	// you'd think the compiler would be smart enough to optimize this
	done = true
	switch v := v.(type) {
	case *Action:
		err = v.UnmarshalABI(dec)
	case *Asset:
		err = v.UnmarshalABI(dec)
	case *Blob:
		err = v.UnmarshalABI(dec)
	case *BlockTimestamp:
		err = v.UnmarshalABI(dec)
	case *Bytes:
		err = v.UnmarshalABI(dec)
	case *Checksum160:
		err = v.UnmarshalABI(dec)
	case *Checksum256:
		err = v.UnmarshalABI(dec)
	case *Checksum512:
		err = v.UnmarshalABI(dec)
	case *Name:
		err = v.UnmarshalABI(dec)
	case *PermissionLevel:
		err = v.UnmarshalABI(dec)
	case *PublicKey:
		err = v.UnmarshalABI(dec)
	case *Signature:
		err = v.UnmarshalABI(dec)
	case *Symbol:
		err = v.UnmarshalABI(dec)
	case *SymbolCode:
		err = v.UnmarshalABI(dec)
	case *TimePoint:
		err = v.UnmarshalABI(dec)
	case *TimePointSec:
		err = v.UnmarshalABI(dec)
	case *Transaction:
		err = v.UnmarshalABI(dec)
	case *TransactionExtension:
		err = v.UnmarshalABI(dec)
	case *TransactionHeader:
		err = v.UnmarshalABI(dec)
	default:
		done = false
	}
	return done, err
}

func NewDecoder(r io.Reader) *abi.Decoder {
	return abi.NewDecoder(r, chainDecoder)
}

func NewCustomDecoder(r io.Reader, fn abi.DecodeFunc) *abi.Decoder {
	return abi.NewDecoder(r, func(dec *abi.Decoder, v interface{}) (done bool, err error) {
		done, err = fn(dec, v)
		if !done && err == nil {
			done, err = chainDecoder(dec, v)
		}
		return done, err
	})
}
