package chain

type KeyType byte

const (
	K1 KeyType = 0
	P1 KeyType = 1
	WA KeyType = 2
)

func (t KeyType) String() string {
	switch t {
	case K1:
		return "K1"
	case P1:
		return "P1"
	case WA:
		return "WA"
	default:
		return "XX"
	}
}
