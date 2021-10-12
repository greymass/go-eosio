package abi

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"reflect"
)

// Fast-path decoding function can be implemented to handle additional types without reflection.
type EncodeFunc func(enc *Encoder, v interface{}) (done bool, err error)

type Encoder struct {
	w  io.Writer
	fn EncodeFunc
}

type Marshaler interface {
	MarshalABI(*Encoder) error
}

func NewEncoder(w io.Writer, fn EncodeFunc) *Encoder {
	return &Encoder{w: w, fn: fn}
}

// Encode given value.
func (enc *Encoder) Encode(v interface{}) error {
	var err error
	// fast path encoding for custom types
	done, err := enc.fn(enc, v)
	if done || err != nil {
		return err
	}
	// fast path encoding for built-in types
	switch v := v.(type) {
	case bool:
		err = enc.WriteBool(v)
	case string:
		err = enc.WriteString(v)

	case uint8:
		err = enc.WriteUint8(v)
	case uint16:
		err = enc.WriteUint16(v)
	case uint32:
		err = enc.WriteUint32(v)
	case uint64:
		err = enc.WriteUint64(v)
	case int8:
		err = enc.WriteInt8(v)
	case int16:
		err = enc.WriteInt16(v)
	case int32:
		err = enc.WriteInt32(v)
	case int64:
		err = enc.WriteInt64(v)
	case float32:
		err = enc.WriteFloat32(v)
	case float64:
		err = enc.WriteFloat64(v)
	case int:
		err = enc.WriteVarint(v)
	case uint:
		err = enc.WriteVaruint(v)
	case []byte:
		err = enc.WriteVaruint(uint(len(v)))
		if err == nil {
			err = enc.WriteBytes(v)
		}
	// reflection for the rest
	default:
		err = enc.EncodeValue(reflect.ValueOf(v))
	}

	return err
}

func (enc *Encoder) EncodeValue(v reflect.Value) error {
	// check if value conforms to Marshaler
	if v.CanInterface() {
		if m, ok := v.Interface().(Marshaler); ok {
			return m.MarshalABI(enc)
		}
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() || !v.Elem().CanInterface() {
			return errors.New("eosio encoder: encountered unexpected nil pointer")
		}
		return enc.Encode(v.Elem().Interface())

	case reflect.Slice:
		l := v.Len()
		err := enc.WriteVaruint(uint(l))
		if err != nil {
			return err
		}
		for i := 0; i < l; i++ {
			if err := enc.Encode(v.Index(i).Interface()); err != nil {
				return err
			}
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := enc.Encode(v.Index(i).Interface()); err != nil {
				return err
			}
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			tag := v.Type().Field(i).Tag.Get("eosio")
			if tag == "optional" {
				exists := !v.Field(i).IsZero()
				enc.WriteBool(exists)
				if !exists {
					continue
				}
			}
			err := enc.Encode(v.Field(i).Interface())
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		l := v.Len()
		err := enc.WriteVaruint(uint(l))
		if err != nil {
			return err
		}
		for _, key := range v.MapKeys() {
			if err := enc.Encode(key.Interface()); err != nil {
				return err
			}
			if err := enc.Encode(v.MapIndex(key).Interface()); err != nil {
				return err
			}
		}

	case reflect.Uint64:
		return enc.WriteUint64(v.Uint())
	case reflect.Uint32:
		return enc.WriteUint32(uint32(v.Uint()))
	case reflect.Int64:
		return enc.WriteInt64(v.Int())
	case reflect.Int32:
		return enc.WriteInt32(int32(v.Int()))

	default:
		return errors.New("eosio encoder: unexpected type " + v.Type().String())
	}

	return nil
}

// writing methods

func (enc *Encoder) WriteBytes(b []byte) error {
	_, err := enc.w.Write(b)
	return err
}

func (enc *Encoder) WriteByte(b byte) error {
	_, err := enc.w.Write([]byte{b})
	return err
}
func (enc *Encoder) WriteBool(v bool) error {
	if v {
		return enc.WriteByte(1)
	}
	return enc.WriteByte(0)
}

func (enc *Encoder) WriteUint8(v uint8) error {
	return enc.WriteByte(byte(v))
}

func (enc *Encoder) WriteUint16(v uint16) error {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], v)
	return enc.WriteBytes(b[:])
}

func (enc *Encoder) WriteUint32(v uint32) error {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	return enc.WriteBytes(b[:])
}

func (enc *Encoder) WriteUint64(v uint64) error {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], v)
	return enc.WriteBytes(b[:])
}

func (enc *Encoder) WriteInt8(v int8) error {
	return enc.WriteByte(byte(v))
}

func (enc *Encoder) WriteInt16(v int16) error {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], uint16(v))
	return enc.WriteBytes(b[:])
}

func (enc *Encoder) WriteInt32(v int32) error {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(v))
	return enc.WriteBytes(b[:])
}

func (enc *Encoder) WriteInt64(v int64) error {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	return enc.WriteBytes(b[:])
}

func (enc *Encoder) WriteVaruint(v uint) error {
	var b [4]byte
	l := binary.PutUvarint(b[:], uint64(v))
	return enc.WriteBytes(b[:l])
}

func (enc *Encoder) WriteVarint(v int) error {
	var b [4]byte
	l := binary.PutVarint(b[:], int64(v))
	return enc.WriteBytes(b[:l])
}

func (enc *Encoder) WriteString(v string) error {
	err := enc.WriteVaruint(uint(len(v)))
	if err != nil {
		return err
	}
	return enc.WriteBytes([]byte(v))
}

func (enc *Encoder) WriteFloat32(v float32) error {
	return enc.WriteUint32(math.Float32bits(v))
}

func (enc *Encoder) WriteFloat64(v float64) error {
	return enc.WriteUint64(math.Float64bits(v))
}
