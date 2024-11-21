package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gitsang/ddns/pkg/configer"
	"github.com/gitsang/ddns/pkg/logi"

	"github.com/spf13/cobra"
)

type LogConfig struct {
	Format    string `json:"format" default:"json" usage:"log format (json|console)"`
	Level     string `json:"level" default:"info" usage:"log level (debug|info|warn|error)"`
	Verbosity int    `json:"verbosity" default:"0" usage:"log verbosity (0-4)"`
	Output    struct {
		Stdout struct {
			Enable bool `json:"enable" default:"true" usage:"enable stdout log"`
		} `json:"stdout"`
		Stderr struct {
			Enable bool `json:"enable" default:"false" usage:"enable stderr log"`
		} `json:"stderr"`
		File struct {
			Enable     bool   `json:"enable" default:"false" usage:"enable file log"`
			Path       string `json:"path" default:"/var/log/ddns/ddns.log" usage:"log file path"`
			MaxSize    string `json:"maxSize" default:"10mb" usage:"log file max size using SI(decimal) standard (K|mb|Gb...)"`
			MaxAge     string `json:"maxAge" default:"7d" usage:"log file max age (d|h|m|s)"`
			MaxBackups int    `json:"maxBackups" default:"3" usage:"log file max backups"`
			Compress   bool   `json:"compress" default:"true" usage:"enable log file compress"`
		} `json:"file"`
	} `json:"output"`
}

type DdnsConfig struct {
	Enable    bool   `json:"enable" default:"false"`
	Type      string `json:"type" default:"A"`
	RR        string `json:"rr" default:"example.com"`
	Interface string `json:"interface" default:"eth0"`
	Prefix    string `json:"prefix" default:"192.168"`
}

type Config struct {
	Log struct {
		Default LogConfig   `json:"default"`
		Fanouts []LogConfig `json:"fanouts"`
	} `json:"log"`
	AccessKeyId     string       `json:"accessKeyId" default:"changit" usage:"aliyun access key id"`
	AccessKeySecret string       `json:"accessKeySecret" default:"changeit" usage:"aliyun access key secret"`
	Domain          string       `json:"domain" default:"example.com" usage:"your domain"`
	UpdateInterval  string       `json:"updateInterval" default:"1h" usage:"the interval to check and update dns record in duration format"`
	Ddnss           []DdnsConfig `json:"ddnss"`
}

var rootCmd = &cobra.Command{
	Use: "ddns",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var rootFlags = struct {
	ConfigPaths []string
}{}

var cfger *configer.Configer

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&rootFlags.ConfigPaths, "config", "c", nil, "config file path")

	cfger = configer.New(
		configer.WithTemplate(new(Config)),
		configer.WithEnvBind(
			configer.WithEnvPrefix("DDNS"),
			configer.WithEnvDelim("_"),
		),
		configer.WithFlagBind(
			configer.WithCommand(rootCmd),
			configer.WithFlagPrefix("ddns"),
			configer.WithFlagDelim("."),
		),
	)
}

func run() {
	// ctx
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// config
	var c Config
	err := cfger.Load(&c, rootFlags.ConfigPaths...)
	if err != nil {
		panic(err)
	}

	// logger
	logh := logi.NewFanOutHandler(
		NewLogHandler(c.Log.Default, c.Log.Fanouts...)...,
	)
	logger := slog.New(logh)
	logger.Info("starting...",
		slog.Any("pid", os.Getpid()),
		slog.Any("flags", rootFlags),
		slog.Any("config", c),
	)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
