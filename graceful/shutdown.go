package graceful

import (
	"context"
	"os"
	"os/signal"
	"time"
)

type ShutdownFunc = func(ctx context.Context)

// OnShutdown block invoking goroutine and wait for SIGINT syscall to invoke a callback
func OnShutdown(callback ShutdownFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	callback(ctx)
}
