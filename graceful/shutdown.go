// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package graceful

import (
	"context"
	"os"
	"os/signal"
)

// OnShutdown block invoking goroutine and wait for os.Interrupt or os.Kill signals to invoke a callback
func OnShutdown(callback func()) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	callback()
}

// Context creates a new context cancelled on os.Interrupt or os.Kill signals
func Context() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
}
