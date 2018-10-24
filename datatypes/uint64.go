package datatypes

import (
	"encoding/binary"

	"github.com/wmnsk/gopcua/errors"
	"github.com/wmnsk/gopcua/id"
)

// UInt64 represents datatype UInt64.
//
// This type exists for handling primitive types in Variant.Value, which should
// implement Data interface.
type UInt64 struct {
	Value int64
}

// DecodeFromBytes decodes given bytes into UInt64.
func (i *UInt64) DecodeFromBytes(b []byte) error {
	if len(b) < 8 {
		return errors.NewErrTooShortToDecode(i, "should be longer")
	}

	i.Value = int64(binary.LittleEndian.Uint64(b[:8]))
	return nil
}

// Serialize serializes UInt64 into bytes.
func (i *UInt64) Serialize() ([]byte, error) {
	b := make([]byte, i.Len())
	if err := i.SerializeTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// SerializeTo serializes UInt64 into bytes.
func (i *UInt64) SerializeTo(b []byte) error {
	binary.LittleEndian.PutUint64(b[:8], uint64(i.Value))

	return nil
}

// Len returns the actual length of Variant in int.
func (i *UInt64) Len() int {
	return 8
}

// DataType returns type of Data.
func (i *UInt64) DataType() uint16 {
	return id.UInt64
}
