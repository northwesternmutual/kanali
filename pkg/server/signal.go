package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/northwesternmutual/kanali/pkg/log"
)

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() context.Context {
	close(onlyOneSignalHandler) // panics when called twice

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		logSignal(<-c)
		cancel()
		logSignal(<-c)
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}

func logSignal(sig interface{}) {
	logger := log.WithContext(nil)

	switch sig {
	case os.Interrupt:
		logger.Info("received SIGINT")
	case syscall.SIGTERM:
		logger.Info("received SIGTERM")
	}
}
