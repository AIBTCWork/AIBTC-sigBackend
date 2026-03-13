package contract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

// abigen --abi abi.json  --type FToken --pkg anchor --out FToken.go
func InitToken() *Token {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/OjJ0OQoKA6YundCinKb2S")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// 加载合约实例 传入合约的地址
	contractAddress := common.HexToAddress(viper.GetString("contract.mint_address"))
	Token, err := NewToken(contractAddress, client)
	if err != nil {
		panic(err)
	}
	return Token
}
