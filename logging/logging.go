package logging

import (
	"github.com/akimsavvin/gonet/di"
	"log/slog"
)

func AddLogging(h slog.Handler) {
	log := slog.New(newMiddleware(h))
	di.AddValue(log)
}
