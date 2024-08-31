// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package graceful

import (
	"context"
	"os"
	"os/signal"
	"time"
)

// ShutdownFunc describes a callback function that will be invoked on shutdown
type ShutdownFunc = func(ctx context.Context)

// OnShutdown block invoking goroutine and wait for os.Interrupt or os.Kill signals to invoke a callback
func OnShutdown(callback ShutdownFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	callback(ctx)
}

// Context creates a new context cancelled on os.Interrupt or os.Kill signals
func Context() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
}
