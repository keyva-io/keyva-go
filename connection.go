// Internal Keyva protocol codec.
//
// This file is an implementation detail of the Keyva client library.
// Do not use directly — use Client instead.
//
// Auto-generated from keyva protocol spec. Do not edit.

package keyva

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

const defaultPort = 6399

type connection struct {
	conn   net.Conn
	reader *bufio.Reader
}

func dial(host string, port int, useTLS bool) (*connection, error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	var c net.Conn
	var err error

	if useTLS {
		c, err = tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, &tls.Config{})
	} else {
		c, err = net.DialTimeout("tcp", addr, 10*time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("keyva: connect %s: %w", addr, err)
	}

	return &connection{
		conn:   c,
		reader: bufio.NewReaderSize(c, 64*1024),
	}, nil
}

func (c *connection) execute(args ...string) (any, error) {
	// Encode command
	buf := fmt.Sprintf("*%d\r\n", len(args))
	for _, arg := range args {
		buf += fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
	}
	if _, err := io.WriteString(c.conn, buf); err != nil {
		return nil, fmt.Errorf("keyva: write: %w", err)
	}
	return c.readFrame()
}

func (c *connection) readFrame() (any, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("keyva: read: %w", err)
	}
	if len(line) < 3 {
		return nil, fmt.Errorf("keyva: short response")
	}
	tag := line[0]
	payload := line[1 : len(line)-2] // strip \r\n

	switch tag {
	case '+':
		return payload, nil
	case '-':
		return nil, parseError(payload)
	case ':':
		n, err := strconv.ParseInt(payload, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("keyva: invalid integer: %s", payload)
		}
		return n, nil
	case '$':
		length, err := strconv.Atoi(payload)
		if err != nil {
			return nil, fmt.Errorf("keyva: invalid bulk length: %s", payload)
		}
		if length < 0 {
			return nil, nil
		}
		data := make([]byte, length+2)
		if _, err := io.ReadFull(c.reader, data); err != nil {
			return nil, fmt.Errorf("keyva: bulk read: %w", err)
		}
		return string(data[:length]), nil
	case '*':
		count, err := strconv.Atoi(payload)
		if err != nil {
			return nil, fmt.Errorf("keyva: invalid array length: %s", payload)
		}
		arr := make([]any, count)
		for i := range count {
			arr[i], err = c.readFrame()
			if err != nil {
				return nil, err
			}
		}
		return arr, nil
	case '%':
		count, err := strconv.Atoi(payload)
		if err != nil {
			return nil, fmt.Errorf("keyva: invalid map length: %s", payload)
		}
		m := make(map[string]any, count)
		for range count {
			key, err := c.readFrame()
			if err != nil {
				return nil, err
			}
			val, err := c.readFrame()
			if err != nil {
				return nil, err
			}
			m[fmt.Sprint(key)] = val
		}
		return m, nil
	case '_':
		return nil, nil
	default:
		return nil, fmt.Errorf("keyva: unknown response type: %c", tag)
	}
}

func (c *connection) close() error {
	return c.conn.Close()
}
