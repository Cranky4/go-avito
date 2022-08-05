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

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &ZTelnetClient{
		timeout: timeout,
		address: address,
		in:      in,
		out:     out,
	}
}

type ZTelnetClient struct {
	connection        net.Conn
	connectionScanner *bufio.Scanner
	inScanner         *bufio.Scanner
	timeout           time.Duration
	address           string
	in                io.ReadCloser
	out               io.Writer
}

func (c *ZTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.connectionScanner = bufio.NewScanner(conn)
	c.connection = conn
	c.inScanner = bufio.NewScanner(c.in)

	return nil
}

func (c *ZTelnetClient) Send() error {
	if !c.inScanner.Scan() {
		return errors.New("end of sending")
	}
	if c.inScanner.Err() != nil {
		return c.inScanner.Err()
	}

	bytes := c.inScanner.Bytes()
	bytes = append(bytes, '\n')
	if _, err := c.connection.Write(bytes); err != nil {
		return err
	}
	return nil
}

func (c *ZTelnetClient) Receive() error {
	if !c.connectionScanner.Scan() {
		return errors.New("end of receiving")
	}
	if c.connectionScanner.Err() != nil {
		return c.connectionScanner.Err()
	}

	bytes := c.connectionScanner.Bytes()
	bytes = append(bytes, '\n')
	if _, err := c.out.Write(bytes); err != nil {
		return err
	}
	return nil
}

func (c *ZTelnetClient) Close() error {
	err := c.connection.Close()
	if err != nil {
		return err
	}
	err = c.in.Close()
	if err != nil {
		return err
	}
	return nil
}
