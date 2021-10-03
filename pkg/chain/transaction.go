package chain

import "github.com/greymass/go-eosio/pkg/abi"

type TransactionHeader struct {
	Expiration       TimePointSec `json:"expiration"`
	RefBlockNum      uint16       `json:"ref_block_num"`
	RefBlockPrefix   uint32       `json:"ref_block_prefix"`
	MaxNetUsageWords uint         `json:"max_net_usage_words"`
	MaxCpuUsageMs    uint8        `json:"max_cpu_usage_ms"`
	DelaySec         uint         `json:"delay_sec"`
}

type TransactionExtension struct {
	Type uint16 `json:"type"`
	Data Bytes  `json:"data"`
}

type Transaction struct {
	TransactionHeader
	ContextFreeActions []Action               `json:"context_free_actions"`
	Actions            []Action               `json:"actions"`
	Extensions         []TransactionExtension `json:"transaction_extensions"`
}

// abi.Marshaler conformance

func (txh TransactionHeader) MarshalABI(e *abi.Encoder) error {
	var err error
	err = txh.Expiration.MarshalABI(e)
	if err != nil {
		return err
	}
	err = e.WriteUint16(txh.RefBlockNum)
	if err != nil {
		return err
	}
	err = e.WriteUint32(txh.RefBlockPrefix)
	if err != nil {
		return err
	}
	err = e.WriteVaruint(uint(txh.MaxNetUsageWords))
	if err != nil {
		return err
	}
	err = e.WriteUint8(txh.MaxCpuUsageMs)
	if err != nil {
		return err
	}
	err = e.WriteVaruint(uint(txh.DelaySec))
	return err
}

func (txe TransactionExtension) MarshalABI(e *abi.Encoder) error {
	var err error
	err = e.WriteUint16(txe.Type)
	if err != nil {
		return err
	}
	err = txe.Data.MarshalABI(e)
	return err
}

func (tx Transaction) MarshalABI(e *abi.Encoder) error {
	var err error
	err = tx.TransactionHeader.MarshalABI(e)
	if err != nil {
		return err
	}
	l := uint(len(tx.ContextFreeActions))
	err = e.WriteVaruint(l)
	if err != nil {
		return err
	}
	for i := uint(0); i < l; i++ {
		err = tx.ContextFreeActions[i].MarshalABI(e)
		if err != nil {
			return err
		}
	}
	l = uint(len(tx.Actions))
	err = e.WriteVaruint(l)
	if err != nil {
		return err
	}
	for i := uint(0); i < l; i++ {
		err = tx.Actions[i].MarshalABI(e)
		if err != nil {
			return err
		}
	}
	l = uint(len(tx.Extensions))
	err = e.WriteVaruint(l)
	if err != nil {
		return err
	}
	for i := uint(0); i < l; i++ {
		err = tx.Extensions[i].MarshalABI(e)
		if err != nil {
			return err
		}
	}
	return err
}

// abi.Unmarshaler conformance

func (txh *TransactionHeader) UnmarshalABI(d *abi.Decoder) error {
	var err error
	err = txh.Expiration.UnmarshalABI(d)
	if err != nil {
		return err
	}
	txh.RefBlockNum, err = d.ReadUint16()
	if err != nil {
		return err
	}
	txh.RefBlockPrefix, err = d.ReadUint32()
	if err != nil {
		return err
	}
	txh.MaxNetUsageWords, err = d.ReadVaruint()
	if err != nil {
		return err
	}
	txh.MaxCpuUsageMs, err = d.ReadUint8()
	if err != nil {
		return err
	}
	txh.DelaySec, err = d.ReadVaruint()
	return err
}

func (txe *TransactionExtension) UnmarshalABI(d *abi.Decoder) error {
	var err error
	txe.Type, err = d.ReadUint16()
	if err != nil {
		return err
	}
	err = txe.Data.UnmarshalABI(d)
	return err
}

func (tx *Transaction) UnmarshalABI(d *abi.Decoder) error {
	var err error
	err = tx.TransactionHeader.UnmarshalABI(d)
	if err != nil {
		return err
	}
	var len uint
	len, err = d.ReadVaruint()
	if err != nil {
		return err
	}
	tx.ContextFreeActions = make([]Action, len)
	for i := 0; i < int(len); i++ {
		err = tx.ContextFreeActions[i].UnmarshalABI(d)
		if err != nil {
			return err
		}
	}
	len, err = d.ReadVaruint()
	if err != nil {
		return err
	}
	tx.Actions = make([]Action, len)
	for i := 0; i < int(len); i++ {
		err = tx.Actions[i].UnmarshalABI(d)
		if err != nil {
			return err
		}
	}
	len, err = d.ReadVaruint()
	if err != nil {
		return err
	}
	tx.Extensions = make([]TransactionExtension, len)
	for i := 0; i < int(len); i++ {
		err = tx.Extensions[i].UnmarshalABI(d)
		if err != nil {
			return err
		}
	}
	return err
}
