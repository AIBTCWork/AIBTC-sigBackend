package handler

import (
	"AI-BTC/internal/contract/service"
	"math/big"
	"net/http"

	"AI-BTC/internal/contract/handler/dto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type ContractHandler struct {
	contractService service.ContractServiceI
}

func NewContractHandler(contractService service.ContractServiceI) *ContractHandler {
	return &ContractHandler{
		contractService: contractService,
	}
}

func (h *ContractHandler) RegisterRoutes(server *gin.Engine) {
	w := server.Group("/wallet")
	{
		w.POST("/sign-transaction", h.SignTransaction)
	}

	c := server.Group("/contract")
	{
		c.GET("/balance", h.BalanceOf)
		c.GET("/claim_amount", h.ClaimAmount)
		c.GET("/total_claim_amount", h.TotolClaimAmount)
	}
}

func (h *ContractHandler) SignTransaction(c *gin.Context) {
	var req dto.TxReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// 将 Amount 乘以 10^18 (将 ETH 单位转换为 Wei)
	amountInt := new(big.Int).Mul(big.NewInt(req.Amount), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	tx := dto.TxDTO{
		Amount: amountInt,
		To:     common.HexToAddress(req.To),
	}
	err, res := h.contractService.SignTransaction(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
func (h *ContractHandler) BalanceOf(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
	}
	account := common.HexToAddress(address)
	res, err := h.contractService.BalanceOf(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"balance": res})
}
func (h *ContractHandler) ClaimAmount(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
	}
	account := common.HexToAddress(address)
	res, err := h.contractService.ClaimAmount(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"balance": res})
}
func (h *ContractHandler) TotolClaimAmount(c *gin.Context) {
	res, err := h.contractService.TotolClaimAmount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"balance": res})
}
