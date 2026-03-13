package main

import (
	"AI-BTC/internal/bot/service"
	"AI-BTC/internal/contract/handler"
	service2 "AI-BTC/internal/contract/service"
	"AI-BTC/ioc"
	"AI-BTC/pkg/contract"
	"AI-BTC/pkg/eip712"
	"AI-BTC/utils"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViperV1()
	app := InitApp()
	// 启动 cron 定时任务
	app.cron.Start()
	err := app.server.Run(viper.GetString("server.port"))

	if err != nil {
		panic(err)
	}

}
func initViperV1() {
	cfile := pflag.String(
		"config",
		filepath.Join(utils.GetRootDir(), "conf/config.yaml"),
		"配置文件路径",
	)
	// 这一步之后，cfile 里面才有值
	pflag.Parse()
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

type App struct {
	server *gin.Engine
	cron   *cron.Cron
}

func InitApp() *App {
	loggerV1 := ioc.InitLogger()
	v := ioc.InitGinMiddlewares(loggerV1)

	botServiceI := service.NewBotService()
	token := contract.InitToken()
	signer := eip712.NewSigner()
	contractServiceI := service2.NewContractService(token, signer)
	contractHandler := handler.NewContractHandler(contractServiceI)
	engine := ioc.InitWebServer(v, contractHandler)
	checkJob := ioc.InitCheckJob(botServiceI, contractServiceI, loggerV1)
	cron := ioc.InitJobs(checkJob, loggerV1)
	app := &App{
		server: engine,
		cron:   cron,
	}
	return app
}

// wire.go:

var ContractSet = wire.NewSet(service2.NewContractService, eip712.NewSigner, contract.InitToken)

var BotSet = wire.NewSet(service.NewBotService)
