package clientgo

import (
	"testing"

	"github.com/drkaka/lg"
	"github.com/stretchr/testify/suite"
)

type SendClientTestSuite struct {
	suite.Suite
}

func (suite *SendClientTestSuite) SetupTest() {
	lg.InitLogger(true)
}

func (suite *SendClientTestSuite) TestUnavailable() {
	cli, avaiable, err := NewSendClient(1986, "127.0.0.1")
	suite.NoError(err, "error create new client")
	suite.Nil(cli, "client should be nil")
	suite.False(avaiable, "server should not available")
}

func (suite *SendClientTestSuite) TestPing() {
	cli, avaiable, err := NewSendClient(19869, "127.0.0.1")
	suite.NoError(err, "error create new client")
	suite.True(avaiable, "server should be available")

	err = cli.Ping()
	suite.NoError(err, "ping error")

	err = cli.conn.Close()
	suite.NoError(err, "conn closed")
}

func (suite *SendClientTestSuite) TestSendAction() {
	cli, avaiable, err := NewSendClient(19869, "127.0.0.1")
	suite.NoError(err, "error create new client")
	suite.True(avaiable, "server should be available")

	err = cli.SendAction("abc")
	suite.NoError(err, "send action error")

	err = cli.conn.Close()
	suite.NoError(err, "conn closed")
}

func TestSendClientTestSuite(t *testing.T) {
	suite.Run(t, new(SendClientTestSuite))
}

func BenchmarkSendActions(b *testing.B) {
	cli, avaiable, err := NewSendClient(19869, "127.0.0.1")
	if !avaiable || err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		cli.SendAction("abc")
	}
}
