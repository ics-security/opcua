package datatypes

import (
	"encoding/binary"

	"github.com/wmnsk/gopcua/errors"
	"github.com/wmnsk/gopcua/id"
)

// Int32 represents datatype Int32.
//
// This type exists for handling primitive types in Variant.Value, which should
// implement Data interface.
type Int32 struct {
	Value int32
}

// DecodeFromBytes decodes given bytes into Int32.
func (i *Int32) DecodeFromBytes(b []byte) error {
	if len(b) < 4 {
		return errors.NewErrTooShortToDecode(i, "should be longer")
	}

	i.Value = int32(binary.LittleEndian.Uint32(b[:4]))
	return nil
}

// Serialize serializes Int32 into bytes.
func (i *Int32) Serialize() ([]byte, error) {
	b := make([]byte, i.Len())
	if err := i.SerializeTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// SerializeTo serializes Int32 into bytes.
func (i *Int32) SerializeTo(b []byte) error {
	binary.LittleEndian.PutUint32(b[:4], uint32(i.Value))

	return nil
}

// Len returns the actual length of Variant in int.
func (i *Int32) Len() int {
	return 4
}

// DataType returns type of Data.
func (i *Int32) DataType() uint16 {
	return id.Int32
}
