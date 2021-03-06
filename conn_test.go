package opcua

import (
	"context"
	"testing"
	"time"

	"github.com/gopcua/opcua/ua"
	"github.com/gopcua/opcua/uacp"
	"github.com/gopcua/opcua/uasc"
)

func TestDial(t *testing.T) {
	t.Skip()
	ctx := context.Background()
	conn, err := uacp.Dial(ctx, "opc.tcp://localhost:4840")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)
	defer conn.Close()
}

func TestSecureChannel(t *testing.T) {
	t.Skip()
	ctx := context.Background()
	conn, err := uacp.Dial(ctx, "opc.tcp://localhost:4840")
	if err != nil {
		t.Fatal(err)
	}
	s := uasc.NewSecureChannel(conn, nil)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}
	defer s.Close()
}

func TestClientRead(t *testing.T) {
	t.Skip()
	c := NewClient("opc.tcp://localhost:4840", nil)
	if err := c.Open(); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	v, err := c.Node(ua.NewNumericNodeID(0, 2258)).Value()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("timex: %v", v)
}
