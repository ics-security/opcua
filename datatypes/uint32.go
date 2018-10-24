// Copyright 2018 gopcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"encoding/binary"

	"github.com/wmnsk/gopcua/errors"
	"github.com/wmnsk/gopcua/id"
)

// UInt32 represents datatype UInt32.
//
// This type exists for handling primitive types in Variant.Value, which should
// implement Data interface.
type UInt32 struct {
	Value uint32
}

// DecodeFromBytes decodes given bytes into UInt32.
func (i *UInt32) DecodeFromBytes(b []byte) error {
	if len(b) < 4 {
		return errors.NewErrTooShortToDecode(i, "should be longer")
	}

	i.Value = binary.LittleEndian.Uint32(b[:4])
	return nil
}

// Serialize serializes UInt32 into bytes.
func (i *UInt32) Serialize() ([]byte, error) {
	b := make([]byte, i.Len())
	if err := i.SerializeTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// SerializeTo serializes UInt32 into bytes.
func (i *UInt32) SerializeTo(b []byte) error {
	binary.LittleEndian.PutUint32(b[:4], uint32(i.Value))

	return nil
}

// Len returns the actual length of Variant in int.
func (i *UInt32) Len() int {
	return 4
}

// DataType returns type of Data.
func (i *UInt32) DataType() uint16 {
	return id.UInt32
}

// Uint32Array represents the array of uint32 type of data.
type Uint32Array struct {
	ArraySize int32
	Values    []uint32
}

// NewUint32Array creates a new NewUint32Array from multiple uint32 values.
func NewUint32Array(vals []uint32) *Uint32Array {
	if vals == nil {
		u := &Uint32Array{
			ArraySize: 0,
		}
		return u
	}

	u := &Uint32Array{
		ArraySize: int32(len(vals)),
	}
	u.Values = append(u.Values, vals...)

	return u
}

// DecodeUint32Array decodes given bytes into Uint32Array.
func DecodeUint32Array(b []byte) (*Uint32Array, error) {
	s := &Uint32Array{}
	if err := s.DecodeFromBytes(b); err != nil {
		return nil, err
	}

	return s, nil
}

// DecodeFromBytes decodes given bytes into Uint32Array.
// TODO: add validation to avoid crash.
func (u *Uint32Array) DecodeFromBytes(b []byte) error {
	u.ArraySize = int32(binary.LittleEndian.Uint32(b[:4]))
	if u.ArraySize <= 0 {
		return nil
	}

	var offset = 4
	for i := 1; i <= int(u.ArraySize); i++ {
		u.Values = append(u.Values, binary.LittleEndian.Uint32(b[offset:offset+4]))
		offset += 4
	}

	return nil
}

// Serialize serializes Uint32Array into bytes.
func (u *Uint32Array) Serialize() ([]byte, error) {
	b := make([]byte, u.Len())
	if err := u.SerializeTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// SerializeTo serializes Uint32Array into bytes.
func (u *Uint32Array) SerializeTo(b []byte) error {
	var offset = 4
	binary.LittleEndian.PutUint32(b[:4], uint32(u.ArraySize))

	for _, v := range u.Values {
		binary.LittleEndian.PutUint32(b[offset:offset+4], v)
		offset += 4
	}

	return nil
}

// Len returns the actual length in int.
func (u *Uint32Array) Len() int {
	return 4 + (len(u.Values) * 4)
}
