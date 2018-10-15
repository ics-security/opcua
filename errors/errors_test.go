// Copyright 2018 gopcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// XXX - Implement!
package errors

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var errorCases = []struct {
	description string
	errorobject error
	errorstring string
}{
	{
		"decode-failure",
		NewErrDecodeFailure(
			[]string{},
			"given bytes too short",
		),
		"failed to decode []string: given bytes too short",
	},
	{
		"serialize-failure",
		NewErrSerializeFailure(
			[]string{},
			"given bytes too short",
		),
		"failed to serialize []string: given bytes too short",
	},
	{
		"unsupported",
		NewErrUnsupported(
			[]string{},
			"not implemented",
		),
		"unsupported []string: not implemented",
	},
	{
		"unexpected",
		NewErrUnexpected(
			"xxx",
			"invalid type",
		),
		"xxx is unexpected: invalid type",
	},
	{
		"network",
		NewErrNetworkNotAvailable(
			&net.TCPConn{},
			"not established",
		),
		"network *net.TCPConn is not available: not established",
	},
}

func TestErrors(t *testing.T) {
	for _, c := range errorCases {
		got := c.errorobject.Error()
		if diff := cmp.Diff(got, c.errorstring); diff != "" {
			t.Error(diff)
		}
	}
}
