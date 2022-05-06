package beefy

import (
	"os"
	"testing"

	"github.com/olegnn/go-substrate-rpc-client/v4/client"
	"github.com/olegnn/go-substrate-rpc-client/v4/config"
)

var beefy *Beefy

func TestMain(m *testing.M) {
	cl, err := client.Connect(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}
	beefy = NewBeefy(cl)
	os.Exit(m.Run())
}
