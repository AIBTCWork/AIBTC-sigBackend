package ioc

import (
	contractHdl "AI-BTC/internal/contract/handler"
	"AI-BTC/pkg/logger"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitWebServer(mdls []gin.HandlerFunc,
	contractHdl *contractHdl.ContractHandler,
) *gin.Engine {
	switch viper.Get("server.mode") {
	case "prod", "production", "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	g := gin.New()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.Use(mdls...)
	contractHdl.RegisterRoutes(g)
	return g
}

func InitGinMiddlewares(
	l logger.LoggerV1) []gin.HandlerFunc {

	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			//AllowOrigins:     []string{"http://localhost:3000"},
			AllowCredentials: true,

			AllowHeaders: []string{"Content-Type", "AccessToken", "X-CSRF-Token", "Authorization", "X-Requested-With", "XMLHttpRequest", "Unique-Finger"},
			// 这个是允许前端访问你的后端响应中带的头部
			ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "Unique-Finger"},
			//AllowHeaders: []string{"content-type"},
			//AllowMethods: []string{"POST"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					//if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "aibtc.work")
			},
			MaxAge: 12 * time.Hour,
		}),
		func(ctx *gin.Context) {
			println("Middleware")
		},
	}
}
