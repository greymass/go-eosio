package chain

import (
	"time"

	"github.com/greymass/go-eosio/pkg/abi"
)

type TimePoint int64

type TimePointSec uint32

type BlockTimestamp uint32

const (
	TimePointFormat    = "2006-01-02T15:04:05.000" // shared with BlockTimestamp
	TimePointSecFormat = "2006-01-02T15:04:05"
)

const blockTimestampEpochMilli = 946684800000

func NewTimePoint(t time.Time) TimePoint {
	return TimePoint(t.UnixMicro())
}

func NewTimePointSec(t time.Time) TimePointSec {
	return TimePointSec(t.Unix())
}

func NewBlockTimestamp(t time.Time) BlockTimestamp {
	// block timestamp is number of half seconds since start of epoch
	return BlockTimestamp((t.Round(time.Millisecond*500).UnixMilli() - blockTimestampEpochMilli) / 500)
}

func NewTimePointFromString(s string) (TimePoint, error) {
	if len(s) > 0 && s[len(s)-1] == 'Z' {
		s = s[:len(s)-1]
	}
	t, err := time.Parse(TimePointFormat, s)
	if err != nil {
		return 0, err
	}
	return NewTimePoint(t), nil
}

func NewTimePointSecFromString(s string) (TimePointSec, error) {
	if len(s) > 0 && s[len(s)-1] == 'Z' {
		s = s[:len(s)-1]
	}
	t, err := time.Parse(TimePointSecFormat, s)
	if err != nil {
		return 0, err
	}
	return NewTimePointSec(t), nil
}

func NewBlockTimestampFromString(s string) (BlockTimestamp, error) {
	if len(s) > 0 && s[len(s)-1] == 'Z' {
		s = s[:len(s)-1]
	}
	t, err := time.Parse(TimePointFormat, s)
	if err != nil {
		return 0, err
	}
	return NewBlockTimestamp(t), nil
}

func (tp TimePoint) Time() time.Time {
	return time.UnixMicro(int64(tp))
}

func (tps TimePointSec) Time() time.Time {
	return time.Unix(int64(tps), 0)
}

func (bts BlockTimestamp) Time() time.Time {
	return time.UnixMilli(int64(bts)*500 + blockTimestampEpochMilli)
}

func (tp TimePoint) String() string {
	return tp.Time().UTC().Format(TimePointFormat)
}

func (tps TimePointSec) String() string {
	return tps.Time().UTC().Format(TimePointSecFormat)
}

func (bts BlockTimestamp) String() string {
	return bts.Time().UTC().Format(TimePointFormat)
}

// abi.Marshaler conformance

func (tp TimePoint) MarshalABI(e *abi.Encoder) error {
	return e.WriteInt64(int64(tp))
}

func (tps TimePointSec) MarshalABI(e *abi.Encoder) error {
	return e.WriteUint32(uint32(tps))
}

func (bts BlockTimestamp) MarshalABI(e *abi.Encoder) error {
	return e.WriteUint32(uint32(bts))
}

// abi.Unmarshaler conformance

func (tp *TimePoint) UnmarshalABI(d *abi.Decoder) error {
	v, err := d.ReadInt64()
	if err == nil {
		*tp = TimePoint(v)
	}
	return err
}

func (tps *TimePointSec) UnmarshalABI(d *abi.Decoder) error {
	v, err := d.ReadUint32()
	if err == nil {
		*tps = TimePointSec(v)
	}
	return err
}

func (bts *BlockTimestamp) UnmarshalABI(d *abi.Decoder) error {
	v, err := d.ReadUint32()
	if err == nil {
		*bts = BlockTimestamp(v)
	}
	return err
}

// encoding.TextMarshaler conformance

func (tp TimePoint) MarshalText() (text []byte, err error) {
	return []byte(tp.String()), nil
}

func (tps TimePointSec) MarshalText() (text []byte, err error) {
	return []byte(tps.String()), nil
}

func (bts BlockTimestamp) MarshalText() (text []byte, err error) {
	return []byte(bts.String()), nil
}

// encoding.TextUnmarshaler conformance

func (tp *TimePoint) UnmarshalText(text []byte) error {
	new, err := NewTimePointFromString(string(text))
	if err == nil {
		*tp = new
	}
	return err
}

func (tps *TimePointSec) UnmarshalText(text []byte) error {
	new, err := NewTimePointSecFromString(string(text))
	if err == nil {
		*tps = new
	}
	return err
}

func (bts *BlockTimestamp) UnmarshalText(text []byte) error {
	new, err := NewBlockTimestampFromString(string(text))
	if err == nil {
		*bts = new
	}
	return err
}
