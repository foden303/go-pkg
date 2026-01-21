package kcp

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/xtaci/kcp-go/v5"
)

var (
	ErrExpectedPong = fmt.Errorf("expected pong message")
)

type Conn struct {
	mu   sync.Mutex
	conn *kcp.UDPSession
	cfg  *Config
}

type DialFunc func(addr string) (net.Conn, error)

func Dial(addr string, cfg *Config) (*Conn, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return DialWithContext(context.Background(), addr, cfg)
}

func DialWithContext(ctx context.Context, addr string, cfg *Config) (*Conn, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	var (
		conn *kcp.UDPSession
		err  error
	)
	for i := 0; i <= cfg.MaxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:

		}
		conn, err = kcp.DialWithOptions(addr, nil, 10, 3)
		if err == nil {
			break
		}
		// Wait before retrying
		if i < cfg.MaxRetries-1 {
			time.Sleep(cfg.RetryDelay)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect after %d retries %w", cfg.MaxRetries, err)
	}
	kcpConn := &Conn{
		conn: conn,
		cfg:  cfg,
	}
	kcpConn.applyConfig()
	return kcpConn, nil

}

// applyConfig applies the configuration settings to the KCP connection.
func (c *Conn) applyConfig() {
	c.conn.SetNoDelay(c.cfg.NoDelay, c.cfg.Interval, c.cfg.Resend, c.cfg.NoCongestion)
	c.conn.SetWindowSize(c.cfg.SendWindowSize, c.cfg.RecvWindowSize)
	c.conn.SetACKNoDelay(c.cfg.ACKNoDelay)
	c.conn.SetMtu(c.cfg.MTU)
	if c.cfg.RateLimit > 0 {
		c.conn.SetWriteBuffer(c.cfg.RateLimit)
	}
}

// Close closes the KCP connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// RemoteAddr returns the remote address of the KCP connection.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// LocalAddr returns the local address of the KCP connection.
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// Raw returns the underlying kcp.UDPSession.
func (c *Conn) Raw() *kcp.UDPSession {
	return c.conn
}

// SetRateLimit sets the rate limit for the KCP connection in bytes per second.
func (c *Conn) SetRateLimit(bytesPerSecond uint32) {
	c.conn.SetRateLimit(bytesPerSecond)
}

// SendMsg sends a message with the specified type and payload using the default write timeout.
func (c *Conn) SendMsg(msgType uint32, payload []byte) error {
	return c.SendMsgWithTimeout(msgType, payload, c.cfg.WriteTimeout)
}

// SendMsgWithTimeout sends a message with the specified type and payload, applying the given timeout.
func (c *Conn) SendMsgWithTimeout(msgType uint32, payload []byte, timeout time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Prepare the message buffer [length: 4 bytes][type: 4 bytes][payload]
	msg := make([]byte, 8+len(payload))
	binary.BigEndian.PutUint32(msg[0:4], uint32(len(payload)))
	binary.BigEndian.PutUint32(msg[4:8], msgType)
	copy(msg[8:], payload)
	if timeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(timeout))
	}
	_, err := c.conn.Write(msg)
	return err
}

// RecvMsg receives a message using the default read timeout.
func (c *Conn) RecvMsg() (msgType uint32, payload []byte, err error) {
	return c.RecvMsgWithTimeout(c.cfg.ReadTimeout)
}

// RecvMsgWithTimeout receives a message, applying the given timeout.
func (c *Conn) RecvMsgWithTimeout(timeout time.Duration) (msgType uint32, payload []byte, err error) {
	if timeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(timeout))
	}
	// Read header
	var header [8]byte
	// read full 8 bytes for header
	if _, err := io.ReadFull(c.conn, header[:]); err != nil {
		return 0, nil, err
	}
	msgLen := binary.BigEndian.Uint32(header[0:4])
	msgType = binary.BigEndian.Uint32(header[4:8])

	// Read payload
	payload = make([]byte, msgLen)
	// read full payload
	_, err = io.ReadFull(c.conn, payload)
	if err != nil {
		return msgType, nil, err
	}
	return msgType, payload, nil
}

// SendJSON marshals the given value to JSON and sends it with the specified message type.
func (c *Conn) SendJSON(msgType uint32, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed json marshal: %w", err)
	}
	return c.SendMsg(msgType, data)
}

// RecvJSON receives a message, un-marshals the JSON payload into the given value.
func (c *Conn) RecvJSON(v interface{}) (msgType uint32, err error) {
	msgType, payload, err := c.RecvMsg()
	if err != nil {
		return 0, err
	}
	if err := json.Unmarshal(payload, v); err != nil {
		return msgType, fmt.Errorf("failed json unmarshal: %w", err)
	}
	return msgType, nil
}

// WriteInt64 writes an int64 value to the connection using big-endian encoding.
func (c *Conn) WriteInt64(v int64) error {
	if c.cfg.WriteTimeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
	}
	return binary.Write(c.conn, binary.BigEndian, v)
}

// ReadInt64 reads an int64 value from the connection using big-endian encoding.
func (c *Conn) ReadInt64() (int64, error) {
	if c.cfg.ReadTimeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.cfg.ReadTimeout))
	}
	var v int64
	err := binary.Read(c.conn, binary.BigEndian, &v)
	return v, err
}

// Write writes data to the KCP connection.
func (c *Conn) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

// Read reads data from the KCP connection.
func (c *Conn) Read(data []byte) (int, error) {
	return c.conn.Read(data)
}

// CopyFrom copies data from the given reader to the KCP connection, reporting progress via the onProgress callback.
func (c *Conn) CopyFrom(r io.Reader, size int64, onProgress func(send int64)) error {
	buf := make([]byte, 32*1024) // 32KB buffer
	var sent int64

	for sent < size {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if c.cfg.WriteTimeout > 0 {
			c.conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
		}
		w, err := c.conn.Write(buf[:n])
		if err != nil {
			return err
		}
		sent += int64(w)
		if onProgress != nil {
			onProgress(sent)
		}
	}
	return nil
}

// CopyTo copies data from the KCP connection to the given writer, reporting progress via the onProgress callback.
func (c *Conn) CopyTo(w io.Writer, size int64, onProgress func(receive int64)) error {
	buf := make([]byte, 32*1024) // 32KB buffer
	var recv int64

	for recv < size {
		remaining := size - recv
		readSize := int64(len(buf))
		if readSize > remaining {
			readSize = remaining
		}
		if c.cfg.ReadTimeout > 0 {
			c.conn.SetReadDeadline(time.Now().Add(c.cfg.ReadTimeout))
		}
		n, err := c.conn.Read(buf[:readSize])
		if err != nil {
			if err != io.EOF {
				break
			}
			return err
		}
		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
		recv += int64(n)
		if onProgress != nil {
			onProgress(recv)
		}
	}
	return nil
}

// Ping sends a ping message and waits for a pong response, returning the round-trip time.
func (c *Conn) Ping(timeout time.Duration) (time.Duration, error) {
	start := time.Now()
	if err := c.SendMsgWithTimeout(MsgTypePing, nil, timeout); err != nil {
		return 0, err
	}
	msgType, _, err := c.RecvMsgWithTimeout(timeout)
	if err != nil {
		return 0, err
	}
	if msgType != MsgTypePong {
		return 0, ErrExpectedPong
	}
	return time.Since(start), nil
}
