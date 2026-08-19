package main

import (
	"context"
	"log"
	"sync"
	"time"
)

const httpShutdownTimeout = 10 * time.Second

type httpServer interface {
	ShutdownWithTimeout(timeout time.Duration) error
}

type mongoClient interface {
	Disconnect(executionContext context.Context) error
}

type shutdownLogger interface {
	Printf(format string, values ...any)
}

func gracefulShutdown(
	executionContext context.Context,
	server httpServer,
	cancelJob context.CancelFunc,
	jobDone *sync.WaitGroup,
	database mongoClient,
	logger shutdownLogger,
) error {
	if logger == nil {
		logger = log.Default()
	}

	logger.Printf("stopping HTTP (at most %s for in-flight requests; returns now if none)", httpShutdownTimeout)
	var httpErr error
	if server != nil {
		httpErr = server.ShutdownWithTimeout(httpShutdownTimeout)
	}
	logger.Printf("HTTP stopped")

	logger.Printf("stopping user count job")
	if cancelJob != nil {
		cancelJob()
	}
	if jobDone != nil {
		jobDone.Wait()
	}
	logger.Printf("user count job stopped")

	logger.Printf("disconnecting Mongo")
	var mongoErr error
	if database != nil {
		mongoErr = database.Disconnect(executionContext)
	}
	logger.Printf("Mongo disconnected")

	if httpErr != nil {
		return httpErr
	}

	return mongoErr
}
