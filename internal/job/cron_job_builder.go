package job

import (
	"AI-BTC/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
)

type CronJobBuilder struct {
	l logger.LoggerV1
}

func NewCronJobBuilder(l logger.LoggerV1, opt prometheus.SummaryOpts) *CronJobBuilder {
	return &CronJobBuilder{
		l: l}
}

func (b *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobAdapterFunc(func() {
		b.l.Debug("开始运行",
			logger.String("name", name))
		err := job.Run()
		if err != nil {
			b.l.Error("执行失败",
				logger.Error(err),
				logger.String("name", name))
		}
		b.l.Debug("结束运行",
			logger.String("name", name))
	})
}

type cronJobAdapterFunc func()

func (c cronJobAdapterFunc) Run() {
	c()
}
