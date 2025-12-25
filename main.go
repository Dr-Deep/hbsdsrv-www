/*
hbsdsrv WorldWideWeb Server
*/
package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/Dr-Deep/hbsdsrv-www/config"
	"github.com/Dr-Deep/hbsdsrv-www/handler"
	"github.com/Dr-Deep/hbsdsrv-www/srv"
	"github.com/Dr-Deep/logging-go"
)

const (
	defaultConfigFilePath = "./config.yml"
	rwForOwnerOnlyPerm    = 0o600
)

func initConfig(configFilePath string) (cfg *config.Configuration) {
	// #nosec G304 -- Zugriff nur auf bekannte Log- und Config-Dateien
	cfgFile, err := os.OpenFile(
		configFilePath,
		os.O_RDONLY,
		rwForOwnerOnlyPerm,
	)

	if errors.Is(err, os.ErrNotExist) {
		panic("config file not found")
	} else if err != nil {
		panic(err)
	}

	cfg, err = config.UnmarshalConfigFile(cfgFile)
	if err != nil {
		panic(err)
	}

	return cfg
}

func initLogger(logFilePath, logLevel string) (logger *logging.Logger) {
	if logFilePath != "" {
		// #nosec G304 -- Zugriff nur auf bekannte Log- und Config-Dateien
		logFile, err := os.OpenFile(
			logFilePath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			rwForOwnerOnlyPerm,
		)
		if err != nil {
			panic(err)
		}

		logger = logging.NewLogger(logFile)
	} else {
		logger = logging.NewLogger(os.Stdout)
	}

	switch logLevel {
	case "debug":
		logger.Level = logging.LogDebug

	case "info":
		logger.Level = logging.LogInfo

	case "error":
		logger.Level = logging.LogError

	case "fatal":
		logger.Level = logging.LogFatal

	default:
		logger.Level = logging.Level(0)
	}

	return logger
}

func main() {
	var (
		cfg    = initConfig(defaultConfigFilePath)
		logger = initLogger(cfg.Logging.File, cfg.Logging.Level)
	)

	var _handlers = []srv.Handler{
		handler.NewHandlerAssets(logger, cfg),
		handler.NewHandlerIndex(),
		handler.NewHandlerLogin(),
		&handler.HandlerContent{},
		&handler.HandlerTroll{},
	}

	www := srv.New(
		http.NewServeMux(),
		_handlers,
		logger,
		cfg,
	)

	if err := www.Start(); err != nil {
		logger.Fatal(err.Error())
	}
}
