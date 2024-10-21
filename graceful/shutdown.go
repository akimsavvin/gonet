// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var sigs = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

// WaitShutdown block invoking goroutine and wait for os.Interrupt or os.Kill signals
func WaitShutdown() {
	ctx, _ := Context()
	<-ctx.Done()
}

// Context creates a new context that will be cancelled on os.Interrupt or os.Kill signals
func Context(parent ...context.Context) (context.Context, context.CancelFunc) {
	if len(parent) > 0 {
		return signal.NotifyContext(parent[0], sigs...)
	}

	return signal.NotifyContext(context.Background(), sigs...)
}
