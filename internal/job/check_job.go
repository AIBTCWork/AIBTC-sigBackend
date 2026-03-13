package job

import (
	bot "AI-BTC/internal/bot/service"
	contarct "AI-BTC/internal/contract/service"
	"AI-BTC/pkg/logger"
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

type CheckJob struct {
	botSvc      bot.BotServiceI
	contractSvc contarct.ContractServiceI
	l           logger.LoggerV1
	timeout     time.Duration
	account     common.Address
}

func NewCheckJob(
	botSvc bot.BotServiceI,
	contractSvc contarct.ContractServiceI,
	l logger.LoggerV1,
	timeout time.Duration,

) *CheckJob {
	return &CheckJob{
		account:     common.HexToAddress(viper.GetString("contract.mint_address")),
		botSvc:      botSvc,
		contractSvc: contractSvc,
		l:           l,
		timeout:     timeout}
}

func (r *CheckJob) Name() string {
	return "check_job"
}

func (r *CheckJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	balance, err := r.contractSvc.Balance()
	if err != nil {
		r.l.Error("failed to get balance", logger.Error(err))
		return err
	}
	r.l.Info("balance", logger.Uint64("balance", balance))
	if balance < 10 {
		r.botSvc.SendMessage(ctx, "balance is less than 10", bot.ExtendFields{})
	}
	return nil
}

func (r *CheckJob) Close() error {
	r.l.Info("burn job closed")
	return nil
}
