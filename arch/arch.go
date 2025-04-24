package arch

import (
	"encoding/binary"
)

var ByteOrder binary.ByteOrder = binary.LittleEndian

type (
	Byte uint8
	Word uint16
)

