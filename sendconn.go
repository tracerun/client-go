package trclient

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/tracerun/tracerun/service"
)

// SendClient used to send data to TraceRun server
type SendClient struct {
	port     uint16
	address  string
	conn     net.Conn
	connLock *sync.Mutex
}

// NewSendClient to create a new send client,
// returning an instance, whether you can reach the server, error
func NewSendClient(port uint16, address string) (*SendClient, bool, error) {
	cli := &SendClient{
		port:     port,
		address:  address,
		connLock: new(sync.Mutex),
	}
	err := cli.getConn()
	if err != nil {
		p(err.Error())

		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			return nil, false, nil
		}
		if operr, ok := err.(*net.OpError); ok && operr.Op == "dial" {
			return nil, false, nil
		}
		return nil, false, err
	}
	return cli, true, nil
}

func (cli *SendClient) getConn() error {
	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	if cli.conn != nil {
		p("reuse send conn")
		return nil
	}

	p("new send conn")
	url := fmt.Sprintf("%s:%d", cli.address, cli.port)
	conn, err := net.DialTimeout("tcp", url, time.Second)
	if err != nil {
		return err
	}
	cli.conn = conn
	go cli.waitForEOF()

	return nil
}

func (cli *SendClient) waitForEOF() {
	cli.connLock.Lock()
	conn := cli.conn
	cli.connLock.Unlock()

	buf := make([]byte, 256)
	_, err := conn.Read(buf)
	if err != nil {
		cli.connLock.Lock()
		conn.Close()
		cli.conn = nil
		cli.connLock.Unlock()
	}
}

// Ping to ping server
func (cli *SendClient) Ping() error {
	err := cli.getConn()
	if err != nil {
		return err
	}

	buf := make([]byte, 3)
	buf[2] = byte(1)

	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	_, err = cli.conn.Write(buf)
	return err
}

// SendAction to send an action
func (cli *SendClient) SendAction(target string, active bool) error {
	err := cli.getConn()
	if err != nil {
		return err
	}

	buf, err := proto.Marshal(&service.Action{
		Target: target,
		Active: active,
	})

	p(fmt.Sprintf("buffer length: %d", len(buf)))
	headerBuf := service.GenerateHeaderBuf(uint16(len(buf)), uint8(2))
	buf = append(headerBuf, buf...)

	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	_, err = cli.conn.Write(buf)
	return err
}

// StopServer to stop the server
func (cli *SendClient) StopServer() error {
	err := cli.getConn()
	if err != nil {
		return err
	}

	headerBuf := service.GenerateHeaderBuf(uint16(0), uint8(0))

	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	_, err = cli.conn.Write(headerBuf)
	return err
}
