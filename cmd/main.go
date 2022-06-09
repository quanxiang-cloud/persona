package main

import (
	"flag"
	"fmt"
	"git.internal.yunify.com/qxp/persona/api/restful"
	"git.internal.yunify.com/qxp/persona/pkg/config"
	"git.internal.yunify.com/qxp/persona/pkg/db/elasticsearch"
	"git.internal.yunify.com/qxp/persona/pkg/misc/logger"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var (
	configPath = flag.String("config", "../configs/config.yml", "-config 配置文件地址")
)

var (
	log logr.Logger
)

func main() {
	flag.Parse()

	err := config.Init(*configPath)
	if err != nil {
		panic(err)
	}

	// init es index
	err = elasticsearch.InitEsIndex(config.Config)
	if err != nil {
		panic(fmt.Sprintf("Create es index error: %s", err))
	}

	// err = logger.New(&config.Config.Log)
	// if err != nil {
	// 	panic(err)
	// }

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	// 启动路由
	router, err := restful.NewRouter(config.Config, log)
	if err != nil {
		panic(err)
	}
	go router.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			router.Close()
			logger.Sync()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
