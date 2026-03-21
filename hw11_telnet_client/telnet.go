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

type Client struct {
	timeout    time.Duration
	address    string
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
	// Add persistent readers
	inReader   *bufio.Reader
	connReader *bufio.Reader
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.connection = conn
	// Initialize the reader for the connection once it exists
	c.connReader = bufio.NewReader(conn)

	return nil
}

func (c *Client) Send() error {
	if c.connection == nil {
		return errors.New("connection is nil")
	}

	msg, err := c.inReader.ReadBytes('\n')
	if err != nil {
		return err
	}

	_, err = c.connection.Write(msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Receive() error {
	if c.connection == nil {
		return errors.New("connection is nil")
	}

	msg, err := c.connReader.ReadBytes('\n')
	if err != nil {
		return err
	}

	_, err = c.out.Write(msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Close() error {
	if c.connection != nil {
		err := c.connection.Close()
		if err != nil {
			return err
		}
	}

	err := c.in.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		timeout:    timeout,
		address:    address,
		connection: nil,
		in:         in,
		out:        out,
		// Initialize the reader for In (Stdin) immediately
		inReader: bufio.NewReader(in),
	}
}
