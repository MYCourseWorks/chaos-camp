package util

// ByteSize represents a constant for byte size.
type ByteSize uint64

const (
	Byte     ByteSize = 1
	KiloByte ByteSize = Byte << 10
	MegaByte ByteSize = KiloByte << 10
	GigaByte ByteSize = MegaByte << 10
	TeraByte ByteSize = GigaByte << 10
)
