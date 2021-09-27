package benchmarks

// copy of chain types with no UnmarshalABI methods so they are forced to use reflection
type TransactionHeader struct {
	Expiration       uint32
	RefBlockNum      uint16
	RefBlockPrefix   uint32
	MaxNetUsageWords uint
	MaxCpuUsageMs    uint8
	DelaySec         uint
}
type TransactionExtension struct {
	Type uint16
	Data []byte
}
type Transaction struct {
	TransactionHeader
	ContextFreeActions []Action
	Actions            []Action
	Extensions         []TransactionExtension
}
type PermissionLevel struct {
	Actor      uint64
	Permission uint64
}
type Action struct {
	Account       uint64
	Name          uint64
	Authorization []PermissionLevel
	Data          []byte
}
