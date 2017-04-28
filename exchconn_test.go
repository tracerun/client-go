package clientgo

import (
	"fmt"
	"testing"

	"github.com/drkaka/lg"
	"github.com/stretchr/testify/suite"
)

type ExchClientTestSuite struct {
	suite.Suite
}

func (suite *ExchClientTestSuite) SetupTest() {
	lg.InitLogger(true)
}

func (suite *ExchClientTestSuite) TestUnavailable() {
	cli, avaiable, err := NewSendClient(1986, "127.0.0.1")
	suite.NoError(err, "error create new client")
	suite.Nil(cli, "client should be nil")
	suite.False(avaiable, "server should not available")
}

func (suite *ExchClientTestSuite) TestGetMeta() {
	cli, avaiable, err := NewExchClient(19869, "127.0.0.1")
	suite.NoError(err, "error create new client")
	suite.True(avaiable, "server should be available")

	meta, err := cli.GetMeta()
	suite.NoError(err, "get all actions error")

	fmt.Println("meta:", meta)
}

func (suite *ExchClientTestSuite) TestGetActions() {
	cli, avaiable, err := NewExchClient(19869, "127.0.0.1")
	suite.NoError(err, "error create new client")
	suite.True(avaiable, "server should be available")

	targets, starts, lasts, err := cli.GetActions()
	suite.NoError(err, "get all actions error")

	fmt.Println("targets:", targets)
	fmt.Println("starts:", starts)
	fmt.Println("lasts:", lasts)
}

func (suite *ExchClientTestSuite) TestGetTargets() {
	cli, avaiable, err := NewExchClient(19869, "127.0.0.1")
	suite.NoError(err, "error create new client")
	suite.True(avaiable, "server should be available")

	targets, err := cli.GetTargets()
	suite.NoError(err, "get targets error")

	fmt.Println("targets:", targets)
}

func (suite *ExchClientTestSuite) TestGetSlots() {
	cli, avaiable, err := NewExchClient(19869, "127.0.0.1")
	suite.NoError(err, "error create new client")
	suite.True(avaiable, "server should be available")

	// get all slots
	starts, slots, err := cli.GetSlots("abc", 0, 0)
	suite.NoError(err, "get targets error")

	fmt.Println("starts:", starts)
	fmt.Println("slots:", slots)
}

func TestExchClientTestSuite(t *testing.T) {
	suite.Run(t, new(ExchClientTestSuite))
}
