package chain

import (
	"bytes"
	"testing"
)

// lebVaruint encodes n as an unsigned LEB128 varuint (the EOSIO wire format).
func lebVaruint(n uint64) []byte {
	var b []byte
	for n >= 0x80 {
		b = append(b, byte(n)|0x80)
		n >>= 7
	}
	return append(b, byte(n))
}

// The generic reflect Slice decode path reads a varuint length then
// reflect.MakeSlice(type, l, l). A huge length must be rejected against the
// remaining input instead of driving a giant allocation (OOM / makeslice panic).
func TestDecodeSliceRejectsHugeLength(t *testing.T) {
	blob := lebVaruint(1 << 50) // ~1.1e15 declared elements, no data follows

	var panicked interface{}
	var err error
	func() {
		defer func() { panicked = recover() }()
		var s []uint32
		err = NewDecoder(bytes.NewReader(blob)).Decode(&s)
	}()
	t.Logf("len(blob)=%d err=%v panicked=%v", len(blob), err, panicked)
	if panicked != nil {
		t.Fatalf("Decode panicked on a huge slice length (DoS): %v", panicked)
	}
	if err == nil {
		t.Fatalf("expected an error for a slice length beyond the input, got nil")
	}
}

// Transaction.UnmarshalABI reads a varuint action count then make([]Action, len).
// Transactions are the most common untrusted blob (get_block / get_transaction
// RPC, P2P, signing requests). A huge count must be rejected before the make.
func TestDecodeTransactionRejectsHugeActionCount(t *testing.T) {
	// 13-byte all-zero TransactionHeader, then a huge context-free-action count.
	blob := append(make([]byte, 13), lebVaruint(1<<50)...)

	var panicked interface{}
	var err error
	func() {
		defer func() { panicked = recover() }()
		var tx Transaction
		err = NewDecoder(bytes.NewReader(blob)).Decode(&tx)
	}()
	t.Logf("len(blob)=%d err=%v panicked=%v", len(blob), err, panicked)
	if panicked != nil {
		t.Fatalf("Decode panicked on a huge action count (DoS): %v", panicked)
	}
	if err == nil {
		t.Fatalf("expected an error for an action count beyond the input, got nil")
	}
}

// A well-formed (empty) transaction still decodes after the fix.
func TestDecodeTransactionValidRoundTrip(t *testing.T) {
	// 13-byte header + three empty vectors (ctx-free actions, actions, extensions).
	blob := append(make([]byte, 13), 0x00, 0x00, 0x00)
	var tx Transaction
	if err := NewDecoder(bytes.NewReader(blob)).Decode(&tx); err != nil {
		t.Fatalf("valid transaction failed to decode: %v", err)
	}
	if len(tx.ContextFreeActions) != 0 || len(tx.Actions) != 0 || len(tx.Extensions) != 0 {
		t.Fatalf("expected empty vectors, got %d/%d/%d",
			len(tx.ContextFreeActions), len(tx.Actions), len(tx.Extensions))
	}
}
