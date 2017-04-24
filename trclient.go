package trclient

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/drkaka/lg"
	"go.uber.org/zap/zapcore"
)

func p(msg string, fields ...zapcore.Field) {
	if lg.L(nil) != nil {
		lg.L(nil).Debug(msg, fields...)
	}
}

// SendClient used to send data to TraceRun server
type SendClient struct {
	port     uint16
	address  string
	conn     net.Conn
	connLock *sync.Mutex
}

// NewSendClient to create a new send client
func NewSendClient(port uint16, address string) (*SendClient, error) {
	cli := &SendClient{
		port:     port,
		address:  address,
		connLock: new(sync.Mutex),
	}
	err := cli.getConn()
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (cli *SendClient) getConn() error {
	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	if cli.conn != nil {
		p("reuse conn")
		return nil
	}

	p("new conn")
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

	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	buf := make([]byte, 3)
	_, err = cli.conn.Write(buf)
	return err
}
