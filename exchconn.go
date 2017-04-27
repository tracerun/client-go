package clientgo

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
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

// GetMeta to get meta information
func (cli *ExchClient) GetMeta() (*service.Meta, error) {
	err := cli.getConn()
	if err != nil {
		return nil, err
	}

	thisRoute := uint8(2)
	headerBuf := service.GenerateHeaderBuf(uint16(0), thisRoute)

	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	_, err = cli.conn.Write(headerBuf)
	if err != nil {
		return nil, err
	}

	data, route, err := service.ReadOne(cli.conn)
	if err != nil {
		return nil, err
	} else if route != thisRoute {
		return nil, ErrRoute
	}

	var meta service.Meta
	err = proto.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}

// GetActions to get all actions
func (cli *ExchClient) GetActions() ([]string, []uint32, []uint32, error) {
	err := cli.getConn()
	if err != nil {
		return nil, nil, nil, err
	}

	thisRoute := uint8(11)
	headerBuf := service.GenerateHeaderBuf(uint16(0), thisRoute)

	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	_, err = cli.conn.Write(headerBuf)
	if err != nil {
		return nil, nil, nil, err
	}

	data, route, err := service.ReadOne(cli.conn)
	if err != nil {
		return nil, nil, nil, err
	} else if route != thisRoute {
		return nil, nil, nil, ErrRoute
	}

	var allActions service.AllActions
	err = proto.Unmarshal(data, &allActions)
	if err != nil {
		return nil, nil, nil, err
	}

	var targets []string
	var starts, lasts []uint32
	for i := 0; i < len(allActions.Actions); i++ {
		targets = append(targets, allActions.Actions[i].Target)
		starts = append(starts, allActions.Actions[i].Start)
		lasts = append(lasts, allActions.Actions[i].Last)
	}

	return targets, starts, lasts, nil
}

// GetTargets to get all the targets
func (cli *ExchClient) GetTargets() ([]string, error) {
	err := cli.getConn()
	if err != nil {
		return nil, err
	}

	thisRoute := uint8(20)
	headerBuf := service.GenerateHeaderBuf(uint16(0), thisRoute)

	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	_, err = cli.conn.Write(headerBuf)
	if err != nil {
		return nil, err
	}

	data, route, err := service.ReadOne(cli.conn)
	if err != nil {
		return nil, err
	} else if route != thisRoute {
		return nil, ErrRoute
	}

	var allTargets service.Targets
	err = proto.Unmarshal(data, &allTargets)
	if err != nil {
		return nil, err
	}

	return allTargets.Target, nil
}

// GetSlots to get slots of a target
func (cli *ExchClient) GetSlots(target string, from, to uint32) ([]uint32, []uint32, error) {
	err := cli.getConn()
	if err != nil {
		return nil, nil, err
	}

	thisRoute := uint8(21)
	buf, err := proto.Marshal(&service.SlotRange{
		Target: target,
		Start:  from,
		End:    to,
	})
	if err != nil {
		return nil, nil, err
	}
	headerBuf := service.GenerateHeaderBuf(uint16(len(buf)), thisRoute)

	cli.connLock.Lock()
	defer cli.connLock.Unlock()

	_, err = cli.conn.Write(append(headerBuf, buf...))
	if err != nil {
		return nil, nil, err
	}

	data, route, err := service.ReadOne(cli.conn)
	if err != nil {
		return nil, nil, err
	} else if route != thisRoute {
		return nil, nil, ErrRoute
	}

	var allSlots service.Slots
	if err := proto.Unmarshal(data, &allSlots); err != nil {
		return nil, nil, err
	}

	var starts, slots []uint32
	for i := 0; i < len(allSlots.Slots); i++ {
		starts = append(starts, allSlots.Slots[i].Start)
		slots = append(slots, allSlots.Slots[i].Slot)
	}

	return starts, slots, nil
}
