package chain

import (
	"fmt"
	"io"
	"reflect"

	"github.com/greymass/go-eosio/pkg/abi"
)

// EOSIO ABI definition, describes the binary representation of a collection of types.
type Abi struct {
	Version          string       `json:"version"`
	Types            []AbiType    `json:"types"`
	Variants         []AbiVariant `json:"variants"`
	Structs          []AbiStruct  `json:"structs"`
	Actions          []AbiAction  `json:"actions"`
	Tables           []AbiTable   `json:"tables"`
	RicardianClauses []AbiClause  `json:"ricardian_clauses"`
}

type AbiType struct {
	NewTypeName string `json:"new_type_name"`
	Type        string `json:"type"`
}

type AbiVariant struct {
	Name  string   `json:"name"`
	Types []string `json:"types"`
}

type AbiStruct struct {
	Name   string     `json:"name"`
	Base   string     `json:"base"`
	Fields []AbiField `json:"fields"`
}

type AbiField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type AbiAction struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	RicardianContract string `json:"ricardian_contract"`
}

type AbiTable struct {
	Name      string   `json:"name"`
	IndexType string   `json:"index_type"`
	KeyNames  []string `json:"key_names"`
	KeyTypes  []string `json:"key_types"`
	Type      string   `json:"type"`
}

type AbiClause struct {
	Id   string `json:"id"`
	Body string `json:"body"`
}

func (a Abi) GetTable(name string) *AbiTable {
	for _, t := range a.Tables {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

func (a Abi) GetAction(name string) *AbiAction {
	for _, t := range a.Actions {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

func (a Abi) GetStruct(name string) *AbiStruct {
	for _, t := range a.Structs {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

func (a Abi) GetType(name string) *AbiType {
	for _, t := range a.Types {
		if t.NewTypeName == name {
			return &t
		}
	}
	return nil
}

func (a Abi) GetVariant(name string) *AbiVariant {
	for _, t := range a.Variants {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

func (a Abi) Decode(r io.Reader, name string) (interface{}, error) {
	res := resolver{&a, make(map[string]*resolvedType)}
	t := res.resolve(name)
	dec := NewDecoder(r)
	var rv interface{}
	err := a.decodeType(dec, t, &rv)
	return rv, err
}

func (a Abi) Encode(w io.Writer, name string, v interface{}) error {
	res := resolver{&a, make(map[string]*resolvedType)}
	t := res.resolve(name)
	enc := NewEncoder(w)
	return a.encodeType(enc, t, v)
}

func (a Abi) encodeType(enc *abi.Encoder, t *resolvedType, v interface{}) error {
	var err error
	exists := v != nil
	if t.isOptional {
		err = enc.WriteBool(exists)
		if !exists || err != nil {
			return nil
		}
	} else if !exists {
		return fmt.Errorf("found nil for non optional %v", t.baseName)
	}
	if t.isArray {
		va, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("expected slice, found %v", reflect.TypeOf(v))
		}
		err := enc.WriteVaruint32(uint32(len(va)))
		if err != nil {
			return err
		}
		for _, e := range va {
			err = a.encodeInner(enc, t, e)
			if err != nil {
				return err
			}
		}
	} else {
		err = a.encodeInner(enc, t, v)
	}
	return err
}

func (a Abi) encodeInner(enc *abi.Encoder, t *resolvedType, v interface{}) error {
	var err error
	if ref := t.ref; ref != nil {
		return a.encodeType(enc, ref, v)
	} else if fields := t.allFields(); fields != nil {
		vs, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map, found %v", reflect.TypeOf(v))
		}
		for _, f := range fields {
			if err = a.encodeType(enc, f.typ, vs[f.name]); err != nil {
				return err
			}
		}
	} else if variant := t.variant; variant != nil {
		va, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("expected slice, found %v", reflect.TypeOf(v))
		}
		if len(va) != 2 {
			return fmt.Errorf("expected two elements, found %v", len(va))
		}
		vn, ok := va[0].(string)
		if !ok {
			return fmt.Errorf("expected string variant name, found %v", reflect.TypeOf(va[0]))
		}
		// iterate over variants and find first that matches vn
		var tv *resolvedType
		var ti int
		for i, tt := range *variant {
			if tt.name == vn {
				tv = tt
				ti = i
				break
			}
		}
		if tv == nil {
			return fmt.Errorf("unknown variant %v", vn)
		}
		err = enc.WriteVaruint32(uint32(ti))
		if err == nil {
			err = a.encodeType(enc, tv, va[1])
		}
	} else {
		var ok bool
		switch t.baseName {
		// encoder builtins
		case "bool":
			var vv bool
			if vv, ok = v.(bool); ok {
				err = enc.WriteBool(vv)
			}
		case "string":
			var vv string
			if vv, ok = v.(string); ok {
				err = enc.WriteString(vv)
			}
		case "uint8":
			var vv uint8
			if vv, ok = v.(uint8); ok {
				err = enc.WriteUint8(vv)
			}
		case "uint16":
			var vv uint16
			if vv, ok = v.(uint16); ok {
				err = enc.WriteUint16(vv)
			}
		case "uint32":
			var vv uint32
			if vv, ok = v.(uint32); ok {
				err = enc.WriteUint32(vv)
			}
		case "uint64":
			var vv uint64
			var vv2 Uint64
			if vv, ok = v.(uint64); ok {
				err = enc.WriteUint64(vv)
			} else if vv2, ok = v.(Uint64); ok {
				err = enc.WriteInt64(int64(vv2))
			}
		case "uint128":
			var vv Uint128
			if vv, ok = v.(Uint128); ok {
				err = vv.MarshalABI(enc)
			}
		case "int8":
			var vv int8
			if vv, ok = v.(int8); ok {
				err = enc.WriteInt8(vv)
			}
		case "int16":
			var vv int16
			if vv, ok = v.(int16); ok {
				err = enc.WriteInt16(vv)
			}
		case "int32":
			var vv int32
			if vv, ok = v.(int32); ok {
				err = enc.WriteInt32(vv)
			}
		case "int64":
			var vv int64
			if vv, ok = v.(int64); ok {
				err = enc.WriteInt64(vv)
			}
		case "int128":
			var vv Int128
			if vv, ok = v.(Int128); ok {
				err = vv.MarshalABI(enc)
			}
		case "float32":
			var vv float32
			if vv, ok = v.(float32); ok {
				err = enc.WriteFloat32(vv)
			}
		case "float64":
			var vv float64
			if vv, ok = v.(float64); ok {
				err = enc.WriteFloat64(vv)
			}
		case "float128":
			var vv Float128
			if vv, ok = v.(Float128); ok {
				err = vv.MarshalABI(enc)
			}
		case "varuint32":
			var vv uint32
			if vv, ok = v.(uint32); ok {
				err = enc.WriteVaruint32(vv)
			}
		case "varint32":
			var vv int32
			if vv, ok = v.(int32); ok {
				err = enc.WriteVarint32(vv)
			}
		case "bytes":
			var vv []byte
			if vv, ok = v.([]byte); ok {
				err = enc.WriteBytes(vv)
			}
		// chain builtins
		case "asset":
			var vv Asset
			if vv, ok = v.(Asset); ok {
				err = vv.MarshalABI(enc)
			}
		case "block_timestamp_type":
			var vv BlockTimestamp
			if vv, ok = v.(BlockTimestamp); ok {
				err = vv.MarshalABI(enc)
			}
		case "checksum160":
			var vv Checksum160
			if vv, ok = v.(Checksum160); ok {
				err = vv.MarshalABI(enc)
			}
		case "checksum256":
			var vv Checksum256
			if vv, ok = v.(Checksum256); ok {
				err = vv.MarshalABI(enc)
			}
		case "checksum512":
			var vv Checksum512
			if vv, ok = v.(Checksum512); ok {
				err = vv.MarshalABI(enc)
			}
		case "eosio::name":
			var vv Name
			if vv, ok = v.(Name); ok {
				err = vv.MarshalABI(enc)
			}
		case "extended_asset":
			var vv ExtendedAsset
			if vv, ok = v.(ExtendedAsset); ok {
				err = vv.MarshalABI(enc)
			}
		case "name":
			var vv Name
			if vv, ok = v.(Name); ok {
				err = vv.MarshalABI(enc)
			}
		case "publickey":
			var vv PublicKey
			if vv, ok = v.(PublicKey); ok {
				err = vv.MarshalABI(enc)
			}
		case "signature":
			var vv Signature
			if vv, ok = v.(Signature); ok {
				err = vv.MarshalABI(enc)
			}
		case "symbol_code":
			var vv SymbolCode
			if vv, ok = v.(SymbolCode); ok {
				err = vv.MarshalABI(enc)
			}
		case "symbol":
			var vv Symbol
			if vv, ok = v.(Symbol); ok {
				err = vv.MarshalABI(enc)
			}
		case "time_point_sec":
			var vv TimePointSec
			if vv, ok = v.(TimePointSec); ok {
				err = vv.MarshalABI(enc)
			}
		case "time_point":
			var vv TimePoint
			if vv, ok = v.(TimePoint); ok {
				err = vv.MarshalABI(enc)
			}
		}
		if !ok && err == nil {
			err = fmt.Errorf("expected %v found %v", t.baseName, reflect.TypeOf(v))
		}
	}

	return err
}

func (a Abi) decodeType(dec *abi.Decoder, t *resolvedType, v *interface{}) error {
	var err error
	if t.isOptional {
		var exists bool
		exists, err = dec.ReadBool()
		if err != nil || !exists {
			return err
		}
	}
	if t.isArray {
		var l uint32
		l, err = dec.ReadVaruint32()
		if err == nil {
			va := make([]interface{}, l)
			for i := uint32(0); i < l; i++ {
				err = a.decodeInner(dec, t, &va[i])
				if err != nil {
					return err // can't recover from this
				}
			}
		}
	} else {
		err = a.decodeInner(dec, t, v)
	}
	if err == io.EOF && t.isExtension {
		return nil
	}
	return err
}

func (a Abi) decodeInner(dec *abi.Decoder, t *resolvedType, v *interface{}) error {
	var err error

	if ref := t.ref; ref != nil {
		return a.decodeType(dec, ref, v)
	} else if fields := t.allFields(); fields != nil {
		vs := make(map[string]interface{})
		for _, field := range fields {
			var fv interface{}
			err := a.decodeType(dec, field.typ, &fv)
			if err != nil {
				return err
			}
			vs[field.name] = fv
		}
		*v = vs

	} else if variant := t.variant; variant != nil {
		var idx uint32
		idx, err = dec.ReadVaruint32()
		if err != nil {
			return err
		}
		if idx >= uint32(len(*variant)) {
			return fmt.Errorf("invalid variant index %d, expected max %d", idx, len(*variant))
		}
		tv := (*variant)[idx]
		var vv []interface{} = make([]interface{}, 2)
		vv[0] = tv.name
		err = a.decodeType(dec, tv, &vv[1])
		if err != nil {
			return err
		}
		*v = vv
	} else {
		switch t.baseName {
		// decoder builtins
		case "string":
			*v, err = dec.ReadString()
		case "bool":
			*v, err = dec.ReadBool()
		case "int8":
			*v, err = dec.ReadInt8()
		case "uint8":
			*v, err = dec.ReadUint8()
		case "int16":
			*v, err = dec.ReadInt16()
		case "uint16":
			*v, err = dec.ReadUint16()
		case "int32":
			*v, err = dec.ReadInt32()
		case "uint32":
			*v, err = dec.ReadUint32()
		case "int64":
			*v, err = dec.ReadInt64()
		case "uint64":
			var uv uint64
			uv, err = dec.ReadUint64()
			*v = Uint64(uv)
		case "int128":
			var rv Int128
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "uint128":
			var rv Uint128
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "float32":
			*v, err = dec.ReadFloat32()
		case "float64":
			*v, err = dec.ReadFloat64()
		case "float128":
			var rv Float128
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "varint32":
			*v, err = dec.ReadVarint32()
		case "varuint32":
			*v, err = dec.ReadVaruint32()
		// chain builtins
		case "asset":
			var rv Asset
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "block_timestamp_type":
			var rv BlockTimestamp
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "bytes":
			var rv Bytes
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "checksum160":
			var rv Checksum160
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "checksum256":
			var rv Checksum256
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "checksum512":
			var rv Checksum512
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "extended_asset":
			var rv ExtendedAsset
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "name":
			var rv Name
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "publickey":
			var rv PublicKey
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "signature":
			var rv Signature
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "symbol_code":
			var rv SymbolCode
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "symbol":
			var rv Symbol
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "time_point_sec":
			var rv TimePointSec
			err = rv.UnmarshalABI(dec)
			*v = rv
		case "time_point":
			var rv TimePoint
			err = rv.UnmarshalABI(dec)
			*v = rv
		default:
			return fmt.Errorf("unknown type %v", t.baseName)
		}
	}
	return err
}

// create a tree of types that's easier to traverse

type resolver struct {
	abi   *Abi
	types map[string]*resolvedType
}

func (r *resolver) resolve(name string) *resolvedType {
	if r.types[name] != nil {
		return r.types[name]
	}
	var isOptional bool
	baseName := name
	if name[len(name)-1] == '?' {
		isOptional = true
		baseName = name[:len(name)-1]
	}
	var isExtension bool
	if baseName[len(baseName)-1] == '$' {
		isExtension = true
		baseName = baseName[:len(baseName)-1]
	}
	var isArray bool
	if baseName[len(baseName)-2] == '[' && baseName[len(baseName)-1] == ']' {
		isArray = true
		baseName = baseName[:len(baseName)-2]
	}

	t := resolvedType{
		name:        name,
		baseName:    baseName,
		isArray:     isArray,
		isOptional:  isOptional,
		isExtension: isExtension,
	}
	r.types[name] = &t

	if as := r.abi.GetStruct(name); as != nil {
		t.fields = &[]*struct {
			name string
			typ  *resolvedType
		}{}
		if as.Base != "" {
			t.base = r.resolve(as.Base)
		}
		for _, f := range as.Fields {
			*t.fields = append(*t.fields, &struct {
				name string
				typ  *resolvedType
			}{
				name: f.Name,
				typ:  r.resolve(f.Type),
			})
		}
	} else if av := r.abi.GetVariant(name); av != nil {
		t.variant = &[]*resolvedType{}
		for _, v := range av.Types {
			vt := r.resolve(v)
			*t.variant = append(*t.variant, vt)
		}
	} else if at := r.abi.GetType(name); at != nil {
		t.ref = r.resolve(at.Type)
	}

	return &t
}

type resolvedType struct {
	name        string
	baseName    string
	isArray     bool
	isOptional  bool
	isExtension bool

	base   *resolvedType
	fields *[]*struct {
		name string
		typ  *resolvedType
	}
	variant *[]*resolvedType
	ref     *resolvedType
}

func (t *resolvedType) String() string {
	return "type<" + t.name + ">"
}

func (rt *resolvedType) allFields() []*struct {
	name string
	typ  *resolvedType
} {
	if rt.fields == nil {
		return nil
	}
	var rv []*struct {
		name string
		typ  *resolvedType
	}
	var seen map[string]bool = make(map[string]bool)
	var cur *resolvedType = rt
	for {
		if cur.fields == nil {
			return nil // invalid struct
		}
		if seen[cur.name] {
			return nil // circular reference
		}
		for i := len(*cur.fields) - 1; i >= 0; i-- {
			// append to front
			rv = append([]*struct {
				name string
				typ  *resolvedType
			}{
				{
					name: (*cur.fields)[i].name,
					typ:  (*cur.fields)[i].typ,
				},
			}, rv...)
		}
		seen[cur.name] = true
		cur = cur.base
		if cur == nil {
			break
		}
	}

	return rv
}
