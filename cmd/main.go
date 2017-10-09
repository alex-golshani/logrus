package main

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/xitonix/logrus"
)

func main() {
	logger := log.New(log.InfoLevel)
	logger.SetFormatter(&log.JSONFormatter{})
	entry := logger.WithField("sample", "nah")

	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j <= 1000; j++ {
				logger.Infof("Hello %s", "XitoniX")
				log.Info("A message from the exported functions")
				entry.AsError().WithError(errors.New("something wrong happened")).Writef("Hi...")
				entry.AsDebug().WithField("test", "test").Write("Do not log me")
			}
		}()
	}

	wg.Wait()
	fmt.Println("Done")
}
