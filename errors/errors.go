// Copyright 2018 gopcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// ErrDecodeFailure indicates the attempt to decode is failed.
type ErrDecodeFailure struct {
	Type    interface{}
	Message string
}

// NewErrDecodeFailure creates a new ErrDecodeFailure with custom message as a reason.
func NewErrDecodeFailure(decodedType interface{}, msg string) *ErrDecodeFailure {
	return &ErrDecodeFailure{
		Type:    decodedType,
		Message: msg,
	}
}

// Error returns the error message for ErrDecodeFailure.
func (e *ErrDecodeFailure) Error() string {
	return fmt.Sprintf("failed to decode %T: %s", e.Type, e.Message)
}

// ErrSerializeFailure indicates the attempt to serialize is failed.
type ErrSerializeFailure struct {
	Type    interface{}
	Message string
}

// NewErrSerializeFailure creates a new ErrSerializeFailure with custom message as a reason.
func NewErrSerializeFailure(serializedType interface{}, msg string) *ErrSerializeFailure {
	return &ErrSerializeFailure{
		Type:    serializedType,
		Message: msg,
	}
}

// Error returns the error message for ErrSerializeFailure.
func (e *ErrSerializeFailure) Error() string {
	return fmt.Sprintf("failed to serialize %T: %s", e.Type, e.Message)
}

// ErrUnsupported indicates the value in Version field is invalid.
type ErrUnsupported struct {
	Type    interface{}
	Message string
}

// NewErrUnsupported creates a ErrUnsupported.
func NewErrUnsupported(unsupportedType interface{}, msg string) *ErrUnsupported {
	return &ErrUnsupported{
		Type:    unsupportedType,
		Message: msg,
	}
}

// Error returns the type of receiver and some additional message.
func (e *ErrUnsupported) Error() string {
	return fmt.Sprintf("unsupported %T: %s", e.Type, e.Message)
}

// ErrUnexpected indicates the given value is not expected.
type ErrUnexpected struct {
	Value   interface{}
	Message string
}

// NewErrUnexpected creates a ErrUnexpected.
func NewErrUnexpected(unexpected interface{}, msg string) *ErrUnexpected {
	return &ErrUnexpected{
		Value:   unexpected,
		Message: msg,
	}
}

// Error returns the type of receiver and some additional message.
func (e *ErrUnexpected) Error() string {
	return fmt.Sprintf("%s is unexpected: %s", e.Value, e.Message)
}

// ErrNetworkNotAvailable indicates a required network object is not available.
type ErrNetworkNotAvailable struct {
	Type    interface{}
	Message string
}

// NewErrNetworkNotAvailable creates a new ErrNetworkNotAvailable with custom message as a reason.
func NewErrNetworkNotAvailable(network interface{}, msg string) *ErrNetworkNotAvailable {
	return &ErrNetworkNotAvailable{
		Type:    network,
		Message: msg,
	}
}

// Error returns the type of network unavailable with custom message.
func (e *ErrNetworkNotAvailable) Error() string {
	return fmt.Sprintf("network %T is not available: %s", e.Type, e.Message)
}

/* XXX - obsoleted

// ErrTooShortToDecode indicates the length of user input is too short to be decoded.
type ErrTooShortToDecode struct {
	Type    interface{}
	Message string
}

// NewErrTooShortToDecode creates a ErrTooShortToDecode.
func NewErrTooShortToDecode(decodedType interface{}, msg string) *ErrTooShortToDecode {
	return &ErrTooShortToDecode{
		Type:    decodedType,
		Message: msg,
	}
}

// Error returns the type of receiver of decoder method and some additional message.
func (e *ErrTooShortToDecode) Error() string {
	return fmt.Sprintf("too short to decode as %T: %s", e.Type, e.Message)
}

// ErrInvalidLength indicates the value in Length field is invalid.
type ErrInvalidLength struct {
	Type    interface{}
	Message string
}

// NewErrInvalidLength creates a ErrInvalidLength.
func NewErrInvalidLength(rcvType interface{}, msg string) *ErrInvalidLength {
	return &ErrInvalidLength{
		Type:    rcvType,
		Message: msg,
	}
}

// Error returns the type of receiver and some additional message.
func (e *ErrInvalidLength) Error() string {
	return fmt.Sprintf("got invalid Length in %T: %s", e.Type, e.Message)
}

// ErrInvalidType indicates the value in Type/Code field is invalid.
type ErrInvalidType struct {
	Type    interface{}
	Action  string
	Message string
}

// NewErrInvalidType creates a ErrInvalidType.
//
// The parameter action is the action taken when this error is raised(e.g., "decode").
func NewErrInvalidType(invalidType interface{}, action, msg string) *ErrInvalidType {
	return &ErrInvalidType{
		Type:    invalidType,
		Action:  action,
		Message: msg,
	}
}

// Error returns the type of receiver and some additional message.
func (e *ErrInvalidType) Error() string {
	return fmt.Sprintf("cannot %s as %T: %s", e.Action, e.Type, e.Message)
}

// ErrReceiverNil indicates the receiver is nil.
type ErrReceiverNil struct {
	Type interface{}
}

// NewErrReceiverNil creates a ErrReceiverNil.
func NewErrReceiverNil(rcvType interface{}) *ErrReceiverNil {
	return &ErrReceiverNil{
		Type: rcvType,
	}
}

// Error returns the type of receiver.
func (e *ErrReceiverNil) Error() string {
	return fmt.Sprintf("Receiver %T is nil.", e.Type)
}
*/

// Cause is just a wrapper of Cause() in pkg/errors
//
// See https://godoc.org/github.com/pkg/errors#Cause for detail.
func Cause(err error) error {
	return errors.Cause(err)
}

// Errorf is just a wrapper of Errorf() in pkg/errors
//
// See https://godoc.org/github.com/pkg/errors#Errorf for detail.
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// New is just a wrapper of New() in pkg/errors
//
// See https://godoc.org/github.com/pkg/errors#New for detail.
func New(text string) error {
	return errors.New(text)
}

// WithMessage is just a wrapper of WithMessage() in pkg/errors
//
// See https://godoc.org/github.com/pkg/errors#WithMessage for detail.
func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

// WithStack is just a wrapper of WithStack() in pkg/errors
//
// See https://godoc.org/github.com/pkg/errors#WithStack for detail.
func WithStack(err error) error {
	return errors.WithStack(err)
}

// Wrap is just a wrapper of Wrap() in pkg/errors
//
// See https://godoc.org/github.com/pkg/errors#Wrap for detail.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf is just a wrapper of Wrapf() in pkg/errors
//
// See https://godoc.org/github.com/pkg/errors#Wrapf for detail.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// Frame is just a wrapper of Frame in pkg/errors
//
// See https://godoc.org/github.com/pkg/errors#Frame for detail.
type Frame = errors.Frame

// StackTrace is just a wrapper of StackTrace in pkg/errors
//
// See https://godoc.org/github.com/pkg/errors#StackTrace for detail.
type StackTrace = errors.StackTrace
