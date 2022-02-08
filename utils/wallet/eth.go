package wallet

import (
	"api.ethscrow/utils"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

var Network *ethclient.Client
var GweiToEth = int64(1000000000)

func SetupEthClient() error {
	cli, err := ethclient.Dial(utils.ETHNET_URL) // TODO:Switch to mainnet when ready
	if err != nil {
		log.Fatalln(err)
		return err
	}
	Network = cli
	return nil
}
