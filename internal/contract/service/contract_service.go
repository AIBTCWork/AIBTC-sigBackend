package service

import (
	"AI-BTC/internal/contract/handler/dto"
	"AI-BTC/pkg/contract"
	"AI-BTC/pkg/eip712"
	"AI-BTC/utils"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

type ContractServiceI interface {
	SignTransaction(tx dto.TxDTO) (e error, txResp dto.TxResp)
	Balance() (uint64, error)
	BalanceOf(account common.Address) (balance uint64, e error)
	Nonce(account common.Address) (uint64, error)
	ClaimAmount(account common.Address) (uint64, error)
	TotolClaimAmount() (uint64, error)
}

type ContractService struct {
	token *contract.Token
	// logger logger.LoggerV1
	eip712Signer eip712.Signer
	auth         *bind.TransactOpts
}

func NewContractService(token *contract.Token, eip712Signer *eip712.Signer) ContractServiceI {
	privateKey, e := crypto.HexToECDSA(viper.GetString("contract.private_key"))
	if e != nil {
		panic("invalid private key")
	}
	auth, e := bind.NewKeyedTransactorWithChainID(
		privateKey,
		big.NewInt(viper.GetInt64("contract.chain_id")),
	)
	if e != nil {
		panic("failed to create transactor")
	}
	return &ContractService{token: token, auth: auth, eip712Signer: *eip712Signer}
}

func (s *ContractService) SignTransaction(tx dto.TxDTO) (e error, txResp dto.TxResp) {
	tx.Nonce, e = s.Nonce(tx.To)
	tx.Deadline = uint64(time.Now().Add(time.Minute * 10).Unix())
	s.eip712Signer.BuildTypedData(tx)
	sig, e := s.eip712Signer.SignTx()
	sigHex := "0x" + common.Bytes2Hex(sig)
	txResp = dto.TxResp{
		Amount:   tx.Amount,
		Deadline: tx.Deadline,
		Nonce:    tx.Nonce,
		Sig:      sigHex,
	}
	return
}
func (s *ContractService) Balance() (uint64, error) {
	balanceBigInt, e := s.token.Balance(nil)
	if e != nil {
		return 0, e
	}
	balance, e := utils.BigIntDivideBy1e18ToUint64(balanceBigInt)
	return balance, e
}
func (s *ContractService) BalanceOf(account common.Address) (uint64, error) {
	balanceBigInt, e := s.token.BalanceOf(nil, account)
	if e != nil {
		return 0, e
	}
	balance, e := utils.BigIntDivideBy1e18ToUint64(balanceBigInt)
	return balance, e
}

func (s *ContractService) Nonce(account common.Address) (uint64, error) {
	nonce, e := s.token.Nonces(nil, account)

	return nonce.Uint64(), e
}

func (s *ContractService) ClaimAmount(to common.Address) (uint64, error) {
	amountBigInt, e := s.token.ClaimAmount(nil, to)
	if e != nil {
		return 0, e
	}
	amount, e := utils.BigIntDivideBy1e18ToUint64(amountBigInt)
	return amount, e
}
func (s *ContractService) TotolClaimAmount() (uint64, error) {
	amountBigInt, e := s.token.TotalAmount(nil)
	if e != nil {
		return 0, e
	}
	amount, e := utils.BigIntDivideBy1e18ToUint64(amountBigInt)
	return amount, e
}
