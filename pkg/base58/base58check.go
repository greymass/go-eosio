package base58

import (
	"crypto/sha256"
	"errors"

	//lint:ignore SA1019 ok google, y u no deprecate sha1
	"golang.org/x/crypto/ripemd160"
)

// ErrChecksum indicates that the checksum of a check-encoded string does not verify against the checksum.
var ErrChecksum = errors.New("checksum error")

// ErrInvalidFormat indicates that the check-encoded string has an invalid format.
var ErrInvalidFormat = errors.New("invalid format: checksum bytes missing")

// checksum: first four bytes of sha256^2
func checksum(input []byte) (cksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(cksum[:], h2[:4])
	return
}

// first four bytes of ripemd160(data + suffix)
func checksumEosio(input []byte, suffix []byte) (cksum [4]byte) {
	h := ripemd160.New()
	h.Write(input)
	if len(suffix) > 0 {
		h.Write(suffix)
	}
	sum := h.Sum(nil)
	copy(cksum[:], sum[:4])
	return
}

// Encode and append a four byte checksum.
func CheckEncode(input []byte) string {
	b := make([]byte, 0, len(input)+4)
	b = append(b, input[:]...)
	cksum := checksum(b)
	b = append(b, cksum[:]...)
	return Encode(b)
}

// Encode and append a four byte ripemd160 checksum.
func CheckEncodeEosio(input []byte, suffix string) string {
	b := make([]byte, 0, len(input)+4)
	b = append(b, input[:]...)
	cksum := checksumEosio(b, []byte(suffix))
	b = append(b, cksum[:]...)
	return Encode(b)
}

// Decode and verify checksum.
func CheckDecode(input string) (result []byte, err error) {
	decoded := Decode(input)
	if len(decoded) < 4 {
		return nil, ErrInvalidFormat
	}
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if checksum(decoded[:len(decoded)-4]) != cksum {
		return nil, ErrChecksum
	}
	payload := decoded[:len(decoded)-4]
	result = append(result, payload...)
	return
}

// Decode and verify ripemd160 checksum.
func CheckDecodeEosio(input string, suffix string) (result []byte, err error) {
	decoded := Decode(input)
	if len(decoded) < 4 {
		return nil, ErrInvalidFormat
	}
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if checksumEosio(decoded[:len(decoded)-4], []byte(suffix)) != cksum {
		return nil, ErrChecksum
	}
	payload := decoded[:len(decoded)-4]
	result = append(result, payload...)
	return
}
