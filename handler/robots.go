/*
 * Serve /robots.txt
 */
package handler

import (
	"github.com/Dr-Deep/hbsdsrv-www/config"
	"github.com/Dr-Deep/logging-go"
)

const text = `User-agent: *
Allow: /`

type HandlerRobots struct {
	cfg    *config.Configuration
	logger *logging.Logger
}
