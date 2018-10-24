package datatypes

import (
	"encoding/binary"

	"github.com/wmnsk/gopcua/errors"
	"github.com/wmnsk/gopcua/id"
)

type Int64 struct {
	Value int64
}

func (i *Int64) DecodeFromBytes(b []byte) error {
	if len(b) < 8 {
		return errors.NewErrTooShortToDecode(i, "should be longer")
	}

	i.Value = int64(binary.LittleEndian.Uint64(b[:8]))
	return nil
}

func (i *Int64) Serialize() ([]byte, error) {
	b := make([]byte, i.Len())
	if err := i.SerializeTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (i *Int64) SerializeTo(b []byte) error {
	binary.LittleEndian.PutUint64(b[:8], uint64(i.Value))

	return nil
}

func (i *Int64) Len() int {
	return 8
}

func (i *Int64) DataType() uint16 {
	return id.Int64
}
