// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"context"
	"github.com/akimsavvin/gonet/shutdown"
	"log"
	"net/http"
)

const DefaultHTTPAddress = ":80"

var globalPrefix = ""

func validatePrefix(prefix string) string {
	if prefix == "/" {
		return ""
	}

	runes := []rune(prefix)

	const slashSymbol = '/'
	if runes[0] != slashSymbol {
		runes = append([]rune{slashSymbol}, runes...)
	}

	if runes[len(runes)-1] == slashSymbol {
		runes = runes[:len(runes)-1]
	}

	return string(runes)
}

func AddGlobalPrefix(prefix string) {
	prefix = validatePrefix(prefix)
	globalPrefix = prefix
}

func StartHttp(addresses ...string) {
	r := newRouter()
	controllers := getControllers()
	for _, controller := range controllers {
		controller.Register(r)
	}

	address := DefaultHTTPAddress

	if len(addresses) > 0 {
		address = addresses[0]
	}

	server := &http.Server{
		Addr:    address,
		Handler: r.Handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Could not run http server due to error: %s\n", err.Error())
		}
	}()

	log.Printf("Started http server on \"%s\"\n", address)

	shutdown.OnGraceful(func(ctx context.Context) {
		log.Printf(server.Shutdown(ctx).Error())
	})
}
