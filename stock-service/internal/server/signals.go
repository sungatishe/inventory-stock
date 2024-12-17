package server

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalHandler() <-chan struct{} {
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		close(stopChan)
	}()

	return stopChan
}
