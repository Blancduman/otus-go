package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		conn:    nil,
	}
}

func (t *telnetClient) Connect() (err error) {
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}

	return nil
}

func (t *telnetClient) Close() (err error) {
	if t.conn != nil {
		err = t.conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return errors.New("no stream connection")
	}

	scanner := bufio.NewScanner(t.in)

	if !scanner.Scan() {
		return errors.New("nothing to read")
	}

	_, err := t.conn.Write(append(scanner.Bytes(), '\n'))
	if err != nil {
		return err
	}

	return nil
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return errors.New("no stream connection")
	}

	scanner := bufio.NewScanner(t.conn)

	if !scanner.Scan() {
		return errors.New("connection closed")
	}

	_, err := t.out.Write(append(scanner.Bytes(), '\n'))
	if err != nil {
		return err
	}

	return nil
}
