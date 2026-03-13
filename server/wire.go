//go:build wireinject

package main

import (
	botService "AI-BTC/internal/bot/service"
	contractHdl "AI-BTC/internal/contract/handler"
	contractService "AI-BTC/internal/contract/service"
	"AI-BTC/ioc"
	"AI-BTC/pkg/contract"
	"AI-BTC/pkg/eip712"

	"github.com/google/wire"
)

// InitApp 为 wire 注入入口：
// - configPath 由 main 传入
// - 返回 *http.Server，main 自行 ListenAndServe
func InitApp() *App {
	wire.Build(
		ioc.InitCheckJob,
		ioc.InitJobs,
		ioc.InitLogger,
		ContractSet,
		BotSet,
		ioc.InitWebServer,
		ioc.InitGinMiddlewares,
		contractHdl.NewContractHandler,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}

var ContractSet = wire.NewSet(
	contractService.NewContractService,
	eip712.NewSigner,
	contract.InitToken,
)

var BotSet = wire.NewSet(
	botService.NewBotService,
)
