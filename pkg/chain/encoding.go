package chain

import (
	"io"

	"github.com/greymass/go-eosio/pkg/abi"
)

// this is dumb but speeds up encoding and decoding by about 5x
// you'd think the compiler would be smart enough to optimize this
// but it isn't, so we do it manually

func chainDecoder(dec *abi.Decoder, v interface{}) (done bool, err error) {
	done = true
	switch v := v.(type) {
	case *Action:
		err = v.UnmarshalABI(dec)
	case *Asset:
		err = v.UnmarshalABI(dec)
	case *Blob:
		err = v.UnmarshalABI(dec)
	case *BlockNum:
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
	case *Float128:
		err = v.UnmarshalABI(dec)
	case *Int128:
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
	case *Uint128:
		err = v.UnmarshalABI(dec)
	case *Uint64:
		err = v.UnmarshalABI(dec)
	default:
		done = false
	}
	return done, err
}

func chainEncoder(enc *abi.Encoder, v interface{}) (done bool, err error) {
	done = true
	switch v := v.(type) {
	case Action:
		err = v.MarshalABI(enc)
	case Asset:
		err = v.MarshalABI(enc)
	case Blob:
		err = v.MarshalABI(enc)
	case BlockNum:
		err = v.MarshalABI(enc)
	case BlockTimestamp:
		err = v.MarshalABI(enc)
	case Bytes:
		err = v.MarshalABI(enc)
	case Checksum160:
		err = v.MarshalABI(enc)
	case Checksum256:
		err = v.MarshalABI(enc)
	case Checksum512:
		err = v.MarshalABI(enc)
	case Float128:
		err = v.MarshalABI(enc)
	case Int128:
		err = v.MarshalABI(enc)
	case Name:
		err = v.MarshalABI(enc)
	case PermissionLevel:
		err = v.MarshalABI(enc)
	case PublicKey:
		err = v.MarshalABI(enc)
	case Signature:
		err = v.MarshalABI(enc)
	case Symbol:
		err = v.MarshalABI(enc)
	case SymbolCode:
		err = v.MarshalABI(enc)
	case TimePoint:
		err = v.MarshalABI(enc)
	case TimePointSec:
		err = v.MarshalABI(enc)
	case Transaction:
		err = v.MarshalABI(enc)
	case TransactionExtension:
		err = v.MarshalABI(enc)
	case TransactionHeader:
		err = v.MarshalABI(enc)
	case Uint128:
		err = v.MarshalABI(enc)
	case Uint64:
		err = v.MarshalABI(enc)
	default:
		done = false
	}
	return done, err
}

func NewDecoder(r io.Reader) *abi.Decoder {
	return abi.NewDecoder(r, chainDecoder)
}

func NewEncoder(w io.Writer) *abi.Encoder {
	return abi.NewEncoder(w, chainEncoder)
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

func NewCustomEncoder(w io.Writer, fn abi.EncodeFunc) *abi.Encoder {
	return abi.NewEncoder(w, func(enc *abi.Encoder, v interface{}) (done bool, err error) {
		done, err = fn(enc, v)
		if !done && err == nil {
			done, err = chainEncoder(enc, v)
		}
		return done, err
	})
}
