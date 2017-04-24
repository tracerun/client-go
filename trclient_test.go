package trclient

import (
	"testing"

	"github.com/drkaka/lg"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
}

func (suite *ClientTestSuite) SetupTest() {
	lg.InitLogger(true)
}

func (suite *ClientTestSuite) TestPing() {
	cli, err := NewSendClient(19869, "127.0.0.1")
	suite.NoError(err, "error create new client")

	err = cli.Ping()
	suite.NoError(err, "ping error")

	err = cli.conn.Close()
	suite.NoError(err, "conn closed")
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
