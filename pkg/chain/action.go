package chain

import "github.com/greymass/go-eosio/pkg/abi"

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
