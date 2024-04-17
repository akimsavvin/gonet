package shutdown

import (
	"context"
	"os"
	"os/signal"
	"time"
)

type ShutdownFunc = func(ctx context.Context)

// OnGraceful block invoking goroutine and wait for SIGINT syscall to invoke a callback
func OnGraceful(callback ShutdownFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	callback(ctx)
}
