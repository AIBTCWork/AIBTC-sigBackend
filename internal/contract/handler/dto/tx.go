package dto

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// 主要用于 gin 的 ShouldBind 做 参数绑定 + 校验（常见 binding:"required" / email / oneof 等）。
type TxReq struct {
	Amount int64  `json:"amount" binding:"required"`
	To     string `json:"to" binding:"required"`
}

type TxDTO struct {
	Amount   *big.Int       `json:"amount"`
	To       common.Address `json:"to"`
	Deadline uint64         `json:"deadline"`
	Nonce    uint64         `json:"nonce"`
}

type TxResp struct {
	Sig      string   `json:"sig"`
	Deadline uint64   `json:"deadline"`
	Nonce    uint64   `json:"nonce"`
	Amount   *big.Int `json:"amount"`
}
