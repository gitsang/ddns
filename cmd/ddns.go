package main

import (
	"ddns/pkg/config"
	"ddns/pkg/service"
	"flag"

	log "github.com/gitsang/golog"
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
		log.WithLogLevel(log.LevelInfo),
		log.WithLogFile(logFile),
		log.WithLogFileCompress(true))

	err = service.DdnsStart()
	if err != nil {
		panic(err)
	}
}
