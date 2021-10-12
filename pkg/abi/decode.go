package abi

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

// Fast-path decoding function can be implemented to handle additional types without reflection.
type DecodeFunc func(dec *Decoder, v interface{}) (done bool, err error)

type Decoder struct {
	r  io.Reader
	fn DecodeFunc
}

type Unmarshaler interface {
	UnmarshalABI(*Decoder) error
}

// Create a new EOSIO ABI decoder. Unless you know what you're doing, you should use the chain.NewDecoder() function instead.
func NewDecoder(r io.Reader, fn DecodeFunc) *Decoder {
	return &Decoder{r: r, fn: fn}
}

// Decode into given value.
func (dec *Decoder) Decode(v interface{}) error {
	var err error
	// fast path decoding for custom types
	done, err := dec.fn(dec, v)
	if done || err != nil {
		return err
	}
	// fast path decoding for built-in types
	switch ptr := v.(type) {
	case *bool:
		*ptr, err = dec.ReadBool()
	case *string:
		*ptr, err = dec.ReadString()

	case *uint8:
		*ptr, err = dec.ReadUint8()
	case *uint16:
		*ptr, err = dec.ReadUint16()
	case *uint32:
		*ptr, err = dec.ReadUint32()
	case *uint64:
		*ptr, err = dec.ReadUint64()
	// TODO: uint128

	case *int8:
		*ptr, err = dec.ReadInt8()
	case *int16:
		*ptr, err = dec.ReadInt16()
	case *int32:
		*ptr, err = dec.ReadInt32()
	case *int64:
		*ptr, err = dec.ReadInt64()
	// TODO: int128

	case *float32:
		*ptr, err = dec.ReadFloat32()
	case *float64:
		*ptr, err = dec.ReadFloat64()

	// varuints represented using golang's undetermined int and uint types
	// we can do this because they are distinct types not just aliases to the underlying int type
	case *int:
		*ptr, err = dec.ReadVarint()
	case *uint:
		*ptr, err = dec.ReadVaruint()

	// variable length bytes, use chain.Bytes or chain.Blob instead to get correct json representations
	case *[]byte:
		var len uint
		len, err = dec.ReadVaruint()
		if err == nil {
			_, *ptr, err = dec.ReadBytes(int(len))
		}

	// reflection for the rest
	default:
		val := reflect.ValueOf(v)
		if val.Kind() != reflect.Ptr || val.IsNil() {
			return fmt.Errorf("abi: invalid type, unable to decode into %s", val.Type())
		}
		err = dec.DecodeValue(val.Elem())
	}

	return err
}

// Decode variant, which must be a struct with all pointer fields.
func (dec *Decoder) DecodeVariant(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("abi: invalid variant type, unable to decode into %s", val.Type())
	}
	// read variant index
	vIdx, err := dec.ReadByte()
	if err != nil {
		return err
	}
	u, ptrValue := indirect(reflect.ValueOf(v), false)
	if u != nil {
		return u.UnmarshalABI(dec)
	}
	// make sure the pointer is a pointer to a struct
	if ptrValue.Kind() != reflect.Struct {
		return fmt.Errorf("abi: invalid variant: expected struct, got %s", ptrValue.Kind())
	}
	// make sure variant index is not out of bounds
	if int(vIdx) >= ptrValue.NumField() {
		return fmt.Errorf("abi: variant index out of bounds: %d", vIdx)
	}
	// enumerate fields and set to nil where needed
	for i := 0; i < ptrValue.NumField(); i++ {
		if i != int(vIdx) && !ptrValue.Field(i).IsNil() {
			ptrValue.Field(i).Set(reflect.Zero(ptrValue.Field(i).Type()))
		}
	}
	variantValuePtr := ptrValue.Field(int(vIdx))
	// make sure the variant field is a pointer
	if variantValuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("abi: invalid variant: expected field pointer, got %s", variantValuePtr.Kind())
	}
	variantType := variantValuePtr.Type().Elem()
	variantValue := reflect.New(variantType).Elem()

	err = dec.Decode(variantValue.Addr().Interface())
	if err != nil {
		return err
	}
	variantValuePtr.Set(variantValue.Addr())

	return nil
}

// Decode into reflected value, you should generally not call this directly.
func (dec *Decoder) DecodeValue(v reflect.Value) error {
	var err error

	u, pv := indirect(v, false)
	if u != nil {
		return u.UnmarshalABI(dec)
	}

	switch pv.Kind() {

	case reflect.Ptr:
		if pv.IsNil() {
			pv.Set(reflect.New(pv.Type().Elem()))
		}
		err = dec.Decode(pv.Interface())

	case reflect.Array:
		l := pv.Len()
		for i := 0; i < l; i++ {
			pv := pv.Index(i)
			if pv.Kind() == reflect.Ptr {
				if pv.IsNil() {
					pv.Set(reflect.New(pv.Type().Elem()))
				}
				err = dec.Decode(pv.Interface())
			} else {
				err = dec.Decode(pv.Addr().Interface())
			}
			if err != nil {
				break
			}
		}

	// maps are packed <varuint32 len>[<key><value>, ..]
	case reflect.Map:
		var l uint
		l, err = dec.ReadVaruint()
		if err == nil {
			// get type of key and value
			keyType := pv.Type().Key()
			valueType := pv.Type().Elem()
			// allocate the map if needed
			if pv.IsNil() || pv.Len() != 0 {
				pv.Set(reflect.MakeMap(pv.Type()))
			}
			for i := 0; i < int(l); i++ {
				// read key
				key := reflect.New(keyType).Elem()
				err = dec.Decode(key.Addr().Interface())
				if err != nil {
					break
				}
				// read value
				value := reflect.New(valueType).Elem()
				err = dec.Decode(value.Addr().Interface())
				if err != nil {
					break
				}
				// set value
				pv.SetMapIndex(key, value)
			}
		}

	case reflect.Struct:
		t := pv.Type()
		l := pv.NumField()
		for i := 0; i < l; i++ {
			if pv := pv.Field(i); pv.CanSet() || t.Field(i).Name != "_" {
				t := t.Field(i)

				var vi interface{}
				if pv.Kind() == reflect.Ptr {
					if pv.IsNil() {
						pv.Set(reflect.New(pv.Type().Elem()))
					}
					vi = pv.Interface()
				} else {
					vi = pv.Addr().Interface()
				}

				tag := t.Tag.Get("eosio")
				if tag == "optional" {
					exists, err := dec.ReadBool()
					if err != nil {
						return err
					}
					if !exists {
						pv.Set(reflect.Zero(pv.Type()))
						continue
					}
				}

				if tag == "variant" {
					err = dec.DecodeVariant(vi)
				} else {
					err = dec.Decode(vi)
				}

				if tag == "extension" && err == io.EOF {
					// TODO: make sure extensions are only last field in a top-level struct
					pv.Set(reflect.Zero(pv.Type()))
					continue
				}

				if err != nil {
					return err
				}
			}
		}

	case reflect.Slice:
		var l uint
		l, err = dec.ReadVaruint()
		if err == nil {
			if pv.Len() != int(l) {
				pv.Set(reflect.MakeSlice(pv.Type(), int(l), int(l)))
			}
			for i := 0; i < int(l); i++ {
				pv := pv.Index(i)
				if pv.Kind() == reflect.Ptr {
					if pv.IsNil() {
						pv.Set(reflect.New(pv.Type().Elem()))
					}
					err = dec.Decode(pv.Interface())
				} else {
					err = dec.Decode(pv.Addr().Interface())
				}
				if err != nil {
					break
				}
			}
		}

	case reflect.Uint64:
		var rv uint64
		rv, err = dec.ReadUint64()
		if err == nil {
			pv.SetUint(rv)
		}

	case reflect.Uint32:
		var rv uint32
		rv, err = dec.ReadUint32()
		if err == nil {
			pv.SetUint(uint64(rv))
		}

	case reflect.Int64:
		var rv int64
		rv, err = dec.ReadInt64()
		if err == nil {
			pv.SetInt(rv)
		}

	case reflect.Int32:
		var rv int32
		rv, err = dec.ReadInt32()
		if err == nil {
			pv.SetInt(int64(rv))
		}

	default:
		return fmt.Errorf("abi: unsupported type: %s %s", pv.Type(), pv.Kind())
	}
	return err
}

// reading methods

func (dec *Decoder) ReadBytes(n int) (an int, b []byte, err error) {
	if n < 0 {
		return 0, nil, errors.New("abi: read with negative count")
	}
	if n == 0 {
		return 0, []byte{}, nil
	}
	b = make([]byte, n)
	an, err = io.ReadFull(dec.r, b)
	if err != nil {
		return an, b, err
	}
	return an, b, nil
}

func (dec *Decoder) ReadByte() (byte, error) {
	_, b, err := dec.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (dec *Decoder) ReadUint64() (uint64, error) {
	_, b, err := dec.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b), nil
}

func (dec *Decoder) ReadUint32() (uint32, error) {
	_, b, err := dec.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}

func (dec *Decoder) ReadUint16() (uint16, error) {
	_, b, err := dec.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func (dec *Decoder) ReadUint8() (uint8, error) {
	b, err := dec.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint8(b), nil
}

func (dec *Decoder) ReadInt64() (int64, error) {
	v, err := dec.ReadUint64()
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

func (dec *Decoder) ReadInt32() (int32, error) {
	v, err := dec.ReadUint32()
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func (dec *Decoder) ReadInt16() (int16, error) {
	v, err := dec.ReadUint16()
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

func (dec *Decoder) ReadInt8() (int8, error) {
	b, err := dec.ReadByte()
	if err != nil {
		return 0, err
	}
	return int8(b), nil
}

func (dec *Decoder) ReadString() (string, error) {
	len, err := dec.ReadVaruint()
	if err != nil {
		return "", err
	}
	_, utf8, err := dec.ReadBytes(int(len))
	if err != nil {
		return "", err
	}
	return string(utf8), nil
}

func (dec *Decoder) ReadBool() (bool, error) {
	b, err := dec.ReadByte()
	return b != 0, err
}

func (dec *Decoder) ReadVaruint() (uint, error) {
	v, err := binary.ReadUvarint(dec)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func (dec *Decoder) ReadVarint() (int, error) {
	v, err := binary.ReadVarint(dec)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func (dec *Decoder) ReadFloat32() (float32, error) {
	b, err := dec.ReadUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(b), nil
}

func (dec *Decoder) ReadFloat64() (float64, error) {
	b, err := dec.ReadUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(b), nil
}

// taken from encoding/json:
// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// If it encounters an Unmarshaler, indirect stops and returns that.
// If decodingNull is true, indirect stops at the first settable pointer so it
// can be set to nil.
func indirect(v reflect.Value, decodingNull bool) (Unmarshaler, reflect.Value) {
	// Issue #24153 indicates that it is generally not a guaranteed property
	// that you may round-trip a reflect.Value by calling Value.Addr().Elem()
	// and expect the value to still be settable for values derived from
	// unexported embedded struct fields.
	//
	// The logic below effectively does this when it first addresses the value
	// (to satisfy possible pointer methods) and continues to dereference
	// subsequent pointers as necessary.
	//
	// After the first round-trip, we set v back to the original value to
	// preserve the original RW flags contained in reflect.Value.
	v0 := v
	haveAddr := false

	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if decodingNull && v.CanSet() {
			break
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		//     var v interface{}
		//     v = &v
		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 && v.CanInterface() {
			if u, ok := v.Interface().(Unmarshaler); ok {
				return u, reflect.Value{}
			}
		}

		if haveAddr {
			v = v0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}
	return nil, v
}
