package eip712

import (
	"AI-BTC/internal/contract/handler/dto"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/spf13/viper"
)

type Signer struct {
	chainID         int64
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
	typedData       apitypes.TypedData
}

func NewSigner() *Signer {
	privateKeyHex := viper.GetString("contract.private_key")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic(err)
	}
	return &Signer{
		chainID:         viper.GetInt64("contract.chain_id"),
		contractAddress: common.HexToAddress(viper.GetString("contract.mint_address")),
		privateKey:      privateKey,
	}
}

func (s *Signer) SignTx() ([]byte, error) {

	hash, _, err := apitypes.TypedDataAndHash(s.typedData)
	if err != nil {
		return nil, err
	}

	sig, err := crypto.Sign(hash, s.privateKey)
	if err != nil {
		return nil, err
	}

	sig[64] += 27 // 转换为以太坊签名格式
	return sig, nil
}

func (s *Signer) BuildTypedData(tx dto.TxDTO) {
	s.typedData = apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"ClaimRequest": {
				{Name: "user", Type: "address"},
				{Name: "amount", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		},
		PrimaryType: "ClaimRequest",
		Domain: apitypes.TypedDataDomain{
			Name:              "Airdrop",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(s.chainID),
			VerifyingContract: s.contractAddress.Hex(),
		},
		Message: apitypes.TypedDataMessage{
			"user":     tx.To.Hex(),
			"amount":   tx.Amount,
			"nonce":    new(big.Int).SetUint64(tx.Nonce),
			"deadline": new(big.Int).SetUint64(tx.Deadline),
		},
	}
}
