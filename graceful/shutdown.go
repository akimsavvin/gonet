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
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)
	<-ch
}

// Context creates a new context that will be cancelled on os.Interrupt or os.Kill signals
func Context(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, sigs...)
}
