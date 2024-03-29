package main

import (
	"ddns/pkg/config"
	"ddns/pkg/service"
	"flag"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

var (
	logLevel = "debug"
	logFile  = "../log/ddns.log"
	confFile = "../conf/ddns.yml"
)

func parseFlag() {
	var c = flag.String("c", "", "set config path")
	var p = flag.String("p", "", "set log path")
	var l = flag.String("l", "", "set log level")
	flag.Parse()

	if *c != "" {
		confFile = *c
	}
	if *p != "" {
		logFile = *p
	}
	if *l != "" {
		logLevel = *l
	}
}

func main() {
	parseFlag()

	var err error
	err = config.LoadConfig(confFile)
	if err != nil {
		panic(err)
	}

	log.InitLogger(
		log.WithLogLevel(config.Cfg.Log.Level),
		log.WithLogFile(config.Cfg.Log.File),
		log.WithLogFileCompress(true))

	log.Info("start ddns", zap.String("version", config.Version),
		zap.Reflect("config", config.Cfg))
	err = service.DdnsStart()
	if err != nil {
		panic(err)
	}
}
