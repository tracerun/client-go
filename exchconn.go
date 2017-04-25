package trclient

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/tracerun/tracerun/service"
)

// ExchClient used to exchange data with TraceRun server
// an one time connection used to get data from server
type ExchClient struct {
	port     uint16
	address  string
	conn     net.Conn
	connLock *sync.Mutex
}

// NewExchClient to create an exchange client.
func NewExchClient(port uint16, address string) (*ExchClient, bool, error) {
	cli := &ExchClient{
		port:     port,
		address:  address,
		connLock: new(sync.Mutex),
	}
	err := cli.getConn()
	if err != nil {
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			return nil, false, nil
		}
		return nil, false, err
	}
	return cli, true, nil
}

func (cli *ExchClient) getConn() error {
	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	if cli.conn != nil {
		p("reuse exchange conn")
		return nil
	}

	p("new exchange conn")
	url := fmt.Sprintf("%s:%d", cli.address, cli.port)
	conn, err := net.DialTimeout("tcp", url, time.Second)
	if err != nil {
		return err
	}
	cli.conn = conn
	return nil
}

// GetActions to get all actions
func (cli *ExchClient) GetActions() error {
	err := cli.getConn()
	if err != nil {
		return err
	}

	headerBuf := service.GenerateHeaderBuf(uint16(0), uint8(3))

	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	_, err = cli.conn.Write(headerBuf)

	count, route, err := service.ReadHeader(cli.conn)
	if err != nil {
		return err
	} else if route != uint8(3) {
		return fmt.Errorf("get wrong response")
	}

	bs, err := service.ReadData(cli.conn, count)
	if err != nil {
		return err
	}

	return nil
}
