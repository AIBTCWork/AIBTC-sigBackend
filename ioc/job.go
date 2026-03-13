package ioc

import (
	bot "AI-BTC/internal/bot/service"
	contract "AI-BTC/internal/contract/service"
	"AI-BTC/internal/job"
	"AI-BTC/pkg/logger"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
)

func InitCheckJob(botSvc bot.BotServiceI, contractSvc contract.ContractServiceI, l logger.LoggerV1) *job.CheckJob {
	return job.NewCheckJob(botSvc, contractSvc, l, time.Second*10)
}

func InitJobs(bj *job.CheckJob, l logger.LoggerV1) *cron.Cron {
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	//build task
	builder := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "geekbang_daming",
		Subsystem: "webook",
		Name:      "cron_job",
		Help:      "定时任务执行",
	})

	cronJob := builder.Build(bj)
	_, err := c.AddJob("@every 5m", cronJob)
	if err != nil {
		l.Error("添加定时任务失败", logger.Error(err))
	}
	return c
}
