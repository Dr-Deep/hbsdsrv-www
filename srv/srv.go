package srv

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Dr-Deep/hbsdsrv-www/config"
	"github.com/Dr-Deep/logging-go"
)

type WWWServer struct {
	// handler, paths, logging

	mux      *http.ServeMux
	handlers []Handler

	logger *logging.Logger
	cfg    *config.Configuration

	interuptSigs chan os.Signal
	sync.Mutex
}

func New(mux *http.ServeMux, handlers []Handler, logger *logging.Logger, cfg *config.Configuration) *WWWServer {
	return &WWWServer{
		mux:          mux,
		handlers:     handlers,
		logger:       logger,
		cfg:          cfg,
		interuptSigs: make(chan os.Signal, 1),
	}
}

func (www *WWWServer) Start() error {
	www.Lock()

	// OS Signals
	signal.Notify(
		www.interuptSigs,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	www.mux.HandleFunc("/", www.Handle)

	www.Unlock()

	return www.run()
}

func (www *WWWServer) run() error {
	defer www.handlePanic()

	go func() {
		for {
			select {
			case <-www.interuptSigs:
				www.Stop()
			}
		}
	}()

	www.logger.Info("listening on", www.cfg.Server.Address)
	return http.ListenAndServe(www.cfg.Server.Address, www.mux)
}

func (www *WWWServer) handlePanic() {
	if r := recover(); r != nil {
		www.logger.Error("PANIC", fmt.Sprintf("%#v", r))
		www.Stop()
	}
}

func (www *WWWServer) Stop() {
	www.logger.Info("received stop signal")
	www.Lock()
	www.logger.Info("stopping...")
	os.Exit(0)
}
