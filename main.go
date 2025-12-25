/*
hbsdsrv WorldWideWeb Server
*/
package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/Dr-Deep/hbsdsrv-www/config"
	"github.com/Dr-Deep/hbsdsrv-www/srv"
	"github.com/Dr-Deep/logging-go"
)

/*
server_name hbsdsrv.1337.cx vmd171781.contaboserver.net;
location /www {
	proxy_pass         https://localhost:8082;
	rewrite ^/www(/.*) $1 break;
	proxy_http_version 1.1;
	proxy_set_header   Host $host;
	proxy_set_header   X-Real-IP $remote_addr;
	proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
	proxy_set_header   X-Forwarded-Proto $scheme;
	proxy_connect_timeout   1m;
	proxy_send_timeout      1m;
   	proxy_read_timeout      1m;
}
*/

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

	www := srv.New(
		http.NewServeMux(),
		logger,
		cfg,
	)

	if err := www.Start(); err != nil {
		logger.Fatal(err.Error())
	}
}
