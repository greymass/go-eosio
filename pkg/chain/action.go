package chain

import (
	"bytes"
	"errors"

	"github.com/greymass/go-eosio/pkg/abi"
)

type PermissionLevel struct {
	Actor      Name `json:"actor"`
	Permission Name `json:"permission"`
}

type Action struct {
	Account       Name              `json:"account"`
	Name          Name              `json:"name"`
	Authorization []PermissionLevel `json:"authorization"`
	Data          Bytes             `json:"data"`
}

func NewAction(account Name, name Name, authorization []PermissionLevel, data Bytes) *Action {
	return &Action{account, name, authorization, data}
}

func (a Action) Decode(abi *Abi) (map[string]interface{}, error) {
	decoded, err := abi.DecodeAction(bytes.NewReader(a.Data), a.Name)
	if err != nil {
		return nil, err
	}
	rv, ok := decoded.(map[string]interface{})
	if !ok {
		return nil, errors.New("action data is not a map[string]interface{}")
	}
	return rv, nil
}

func (a Action) DecodeInto(v interface{}) error {
	return NewDecoder(bytes.NewReader(a.Data)).Decode(v)
}

func (a Action) Digest() Checksum256 {
	b := bytes.NewBuffer(nil)
	err := a.MarshalABI(NewEncoder(b))
	if err != nil {
		panic(err)
	}
	return Checksum256Digest(b.Bytes())
}

// abi.Marshaler conformance

func (pl PermissionLevel) MarshalABI(e *abi.Encoder) error {
	pl.Actor.MarshalABI(e)
	return pl.Permission.MarshalABI(e)
}

func (a Action) MarshalABI(e *abi.Encoder) error {
	var err error
	a.Account.MarshalABI(e)
	a.Name.MarshalABI(e)
	l := uint32(len(a.Authorization))
	err = e.WriteVaruint32(l)
	if err != nil {
		return err
	}
	for i := 0; i < int(l); i++ {
		err = a.Authorization[i].MarshalABI(e)
		if err != nil {
			return err
		}
	}
	return a.Data.MarshalABI(e)
}

// abi.Unmarshaler conformance

func (pl *PermissionLevel) UnmarshalABI(d *abi.Decoder) error {
	pl.Actor.UnmarshalABI(d)
	return pl.Permission.UnmarshalABI(d)
}

func (a *Action) UnmarshalABI(d *abi.Decoder) error {
	a.Account.UnmarshalABI(d)
	a.Name.UnmarshalABI(d)
	l, err := d.ReadVaruint32()
	if err != nil {
		return err
	}
	for i := 0; i < int(l); i++ {
		var pl PermissionLevel
		err = pl.UnmarshalABI(d)
		if err != nil {
			return err
		}
		a.Authorization = append(a.Authorization, pl)
	}
	return a.Data.UnmarshalABI(d)
}
