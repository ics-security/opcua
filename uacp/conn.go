package uacp

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/gopcua/opcua/ua"
	"github.com/gopcua/opcua/utils"
)

const (
	KB = 1024
	MB = 1024 * KB

	DefaultReceiveBufSize = 0xffff
	DefaultSendBufSize    = 0xffff
	DefaultMaxChunkCount  = 512
	DefaultMaxMessageSize = 2 * MB
)

// connid stores the current connection id. updated with atomic.AddUint32
var connid uint32

// nextid returns the next connection id
func nextid() uint32 {
	return atomic.AddUint32(&connid, 1)
}

func Dial(ctx context.Context, endpoint string) (*Conn, error) {
	log.Printf("Connect to %s", endpoint)
	network, raddr, err := utils.ResolveEndpoint(endpoint)
	if err != nil {
		return nil, err
	}
	c, err := net.DialTCP(network, nil, raddr)
	if err != nil {
		return nil, err
	}

	conn := &Conn{
		id: nextid(),
		c:  c,
		ack: &Acknowledge{
			ReceiveBufSize: DefaultReceiveBufSize,
			SendBufSize:    DefaultSendBufSize,
			MaxChunkCount:  0, // use what the server wants
			MaxMessageSize: 0, // use what the server wants
		},
	}

	log.Printf("conn %d: start HEL/ACK handshake", conn.id)
	if err := conn.handshake(endpoint); err != nil {
		log.Printf("conn %d: HEL/ACK handshake failed: %s", conn.id, err)
		conn.Close()
		return nil, err
	}
	return conn, nil
}

// Listener is a OPC UA Connection Protocol network listener.
type Listener struct {
	l        net.Listener
	ack      *Acknowledge
	endpoint string
}

// Listen acts like net.Listen for OPC UA Connection Protocol networks.
//
// Currently the endpoint can only be specified in "opc.tcp://<addr[:port]>/path" format.
//
// If the IP field of laddr is nil or an unspecified IP address, Listen listens
// on all available unicast and anycast IP addresses of the local system.
// If the Port field of laddr is 0, a port number is automatically chosen.
func Listen(endpoint string, ack *Acknowledge) (*Listener, error) {
	if ack == nil {
		ack = &Acknowledge{
			ReceiveBufSize: DefaultReceiveBufSize,
			SendBufSize:    DefaultSendBufSize,
			MaxChunkCount:  DefaultMaxChunkCount,
			MaxMessageSize: DefaultMaxMessageSize,
		}
	}

	network, laddr, err := utils.ResolveEndpoint(endpoint)
	if err != nil {
		return nil, err
	}
	l, err := net.Listen(network, laddr.String())
	if err != nil {
		return nil, err
	}
	return &Listener{
		l:        l,
		ack:      ack,
		endpoint: endpoint,
	}, nil
}

// Accept accepts the next incoming call and returns the new connection.
//
// The first param ctx is to be passed to monitor(), which monitors and handles
// incoming messages automatically in another goroutine.
func (l *Listener) Accept(ctx context.Context) (*Conn, error) {
	c, err := l.l.Accept()
	if err != nil {
		return nil, err
	}
	conn := &Conn{nextid(), c, l.ack}
	if err := conn.srvhandshake(l.endpoint); err != nil {
		c.Close()
		return nil, err
	}
	return conn, nil
}

// Close closes the Listener.
func (l *Listener) Close() error {
	return l.l.Close()
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.l.Addr()
}

// Endpoint returns the listener's EndpointURL.
func (l *Listener) Endpoint() string {
	return l.endpoint
}

type Conn struct {
	id  uint32
	c   net.Conn
	ack *Acknowledge
}

func (c *Conn) ID() uint32 {
	return c.id
}

func (c *Conn) ReceiveBufSize() uint32 {
	return c.ack.ReceiveBufSize
}

func (c *Conn) SendBufSize() uint32 {
	return c.ack.SendBufSize
}

func (c *Conn) MaxMessageSize() uint32 {
	return c.ack.MaxMessageSize
}

func (c *Conn) MaxChunkCount() uint32 {
	return c.ack.MaxChunkCount
}

func (c *Conn) Close() error {
	log.Printf("conn %d: close", c.id)
	return c.c.Close()
}

func (c *Conn) Read(b []byte) (int, error) {
	return c.c.Read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
	return c.c.Write(b)
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}

func (c *Conn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

func (c *Conn) handshake(endpoint string) error {
	hel := &Hello{
		Version:        c.ack.Version,
		ReceiveBufSize: c.ack.ReceiveBufSize,
		SendBufSize:    c.ack.SendBufSize,
		MaxMessageSize: c.ack.MaxMessageSize,
		MaxChunkCount:  c.ack.MaxChunkCount,
		EndPointURL:    endpoint,
	}

	if err := c.send("HELF", hel); err != nil {
		return err
	}

	b, err := c.recv()
	if err != nil {
		return err
	}

	msgtyp := string(b[:4])
	if msgtyp != "ACKF" {
		return fmt.Errorf("got %s want ACK", msgtyp)
	}

	ack := new(Acknowledge)
	if _, err := ua.Decode(b[hdrlen:], ack); err != nil {
		return fmt.Errorf("decode ACK failed: %s", err)
	}

	if ack.Version != 0 {
		return fmt.Errorf("invalid version %d", ack.Version)
	}
	if ack.MaxChunkCount == 0 {
		ack.MaxChunkCount = DefaultMaxChunkCount
		log.Printf("conn %d: server has no chunk limit. Using %d", c.id, ack.MaxChunkCount)
	}
	if ack.MaxMessageSize == 0 {
		ack.MaxMessageSize = DefaultMaxMessageSize
		log.Printf("conn %d: server has no message size limit. Using %d", c.id, ack.MaxMessageSize)
	}
	c.ack = ack
	log.Printf("conn %d: recv ACK:%v", c.id, ack)
	return nil
}

func (c *Conn) srvhandshake(endpoint string) error {
	b, err := c.recv()
	if err != nil {
		c.sendError(BadTCPInternalError)
		return err
	}

	// HEL or RHE?
	msgtyp := string(b[:4])
	msg := b[hdrlen:]
	switch msgtyp {
	case "HELF":
		hel := new(Hello)
		if _, err := ua.Decode(msg, hel); err != nil {
			c.sendError(BadTCPInternalError)
			return err
		}
		if hel.EndPointURL != endpoint {
			c.sendError(BadTCPEndpointURLInvalid)
			return fmt.Errorf("invalid endpoint url %s", hel.EndPointURL)
		}
		if err := c.send("ACKF", c.ack); err != nil {
			c.sendError(BadTCPInternalError)
			return err
		}
		return nil

	case "RHEF":
		rhe := new(ReverseHello)
		if _, err := ua.Decode(msg, rhe); err != nil {
			c.sendError(BadTCPInternalError)
			return err
		}
		if rhe.EndPointURL != endpoint {
			c.sendError(BadTCPEndpointURLInvalid)
			return fmt.Errorf("invalid endpoint url %s", rhe.EndPointURL)
		}
		log.Printf("conn %d: connecting to %s", c.id, rhe.ServerURI)
		c.c.Close()
		c, err := Dial(context.Background(), rhe.ServerURI)
		if err != nil {
			return err
		}
		c.c = c
		return nil

	default:
		c.sendError(BadTCPInternalError)
		return fmt.Errorf("invalid handshake packet %q", msgtyp)
	}
}

func (c *Conn) sendError(code uint32) {
	// we swallow the error to silence complaints from the linter
	// since sending an error will close the connection and we
	// want to bubble a different error up.
	_ = c.send("ERRF", &Error{Error: code})
}

// hdrlen is the size of the uacp header
const hdrlen = 8

// recv receives a message from the stream and returns it without the header.
func (c *Conn) recv() ([]byte, error) {
	hdr := make([]byte, hdrlen)
	_, err := io.ReadFull(c.c, hdr)
	if err != nil {
		return nil, fmt.Errorf("hdr read faled: %s", err)
	}

	var h Header
	if _, err := ua.Decode(hdr, &h); err != nil {
		return nil, fmt.Errorf("hdr decode failed: %s", err)
	}

	if h.MessageSize > c.ack.ReceiveBufSize {
		return nil, fmt.Errorf("packet too large: %d > %d bytes", h.MessageSize, c.ack.ReceiveBufSize)
	}

	b := make([]byte, h.MessageSize-hdrlen)
	if _, err := io.ReadFull(c.c, b); err != nil {
		return nil, fmt.Errorf("read msg failed: %s", err)
	}

	log.Printf("conn %d: recv %s%c with %d bytes", c.id, h.MessageType, h.ChunkType, len(b))
	return append(hdr, b...), nil
}

func (c *Conn) send(typ string, msg interface{}) error {
	if len(typ) != 4 {
		return fmt.Errorf("invalid msg type: %s", typ)
	}

	body, err := ua.Encode(msg)
	if err != nil {
		return fmt.Errorf("encode msg failed: %s", err)
	}

	h := Header{
		MessageType: typ[:3],
		ChunkType:   typ[3],
		MessageSize: uint32(len(body) + 8),
	}

	if h.MessageSize > c.ack.SendBufSize {
		return fmt.Errorf("send packet too large: %d > %d bytes", h.MessageSize, c.ack.SendBufSize)
	}

	hdr, err := h.Encode()
	if err != nil {
		return fmt.Errorf("encode hdr failed: %s", err)
	}

	b := append(hdr, body...)
	if _, err := c.c.Write(b); err != nil {
		return fmt.Errorf("write failed: %s", err)
	}
	log.Printf("conn %d: sent %s with %d bytes", c.id, typ, len(b))

	return nil
}
