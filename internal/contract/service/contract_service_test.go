package service

import (
	"AI-BTC/internal/contract/handler/dto"
	"AI-BTC/pkg/contract"
	"AI-BTC/pkg/eip712"
	"AI-BTC/utils"
	"log"
	"math/big"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// 测试主流程，创建推事件，创建拉事件
// 运行的时候要注意工作目录要定位到当前目录
type ContractTestSuite struct {
	suite.Suite
	ContractServiceI ContractServiceI
}

func (f *ContractTestSuite) SetupSuite() {
	initViperV1()
	f.ContractServiceI = NewContractService(
		contract.InitToken(),
		eip712.NewSigner(),
	)
}

func (f *ContractTestSuite) Test_Nonce() {
	nonce, err := f.ContractServiceI.Nonce(common.HexToAddress("0x3cd34cb4933bb0572E26837D7DCd929019CA8d44"))
	log.Println(nonce)
	f.Require().Nil(err)
}

func (f *ContractTestSuite) Test_Balance() {
	balance, err := f.ContractServiceI.Balance()
	log.Println(balance)
	// f.Assert().Equal(uint64(0), balance)
	f.Require().Nil(err)
}

func (f *ContractTestSuite) Test_BalanceOf() {
	balance, err := f.ContractServiceI.BalanceOf(common.HexToAddress("0xF321240E4E2880910f24d67B8C9B9299767988fe"))
	log.Println(balance)
	// f.Assert().Equal(uint64(0), balance)
	f.Require().Nil(err)
}
func (f *ContractTestSuite) Test_SignTransaction() {
	deadline := uint64(time.Now().Add(time.Hour).Unix())
	amountInt := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	e, sig := f.ContractServiceI.SignTransaction(dto.TxDTO{
		Amount:   amountInt,
		Deadline: deadline,
		Nonce:    0,
		To:       common.HexToAddress("0xF321240E4E2880910f24d67B8C9B9299767988fe"),
	})
	f.Assert().Nil(e)
	log.Println(sig)
	// []byte 转十六进制字符串

	// 十六进制字符串转 []byte（注意：Hex2Bytes 不处理 0x 前缀）
	// 如果需要从十六进制字符串转回 []byte：
	// sigBytes := common.Hex2Bytes(sigHex[2:])
	// log.Println("转换回 []byte 长度:", len(sigBytes))
	// log.Println(sigBytes)
}
func initViperV1() {
	cfile := pflag.String(
		"config",
		filepath.Join(utils.GetRootDir(), "conf/config.yaml"),
		"配置文件路径",
	)
	// 这一步之后，cfile 里面才有值
	// pflag.Parse()
	// 所有的默认值放好s
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")

	log.Println(val)
}

func TestFeedTestSuite(t *testing.T) {
	suite.Run(t, new(ContractTestSuite))
}
