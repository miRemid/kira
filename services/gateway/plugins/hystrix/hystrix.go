package hystrix

import (
	"context"
	"log"
	"net/http"

	"github.com/micro/micro/v2/plugin"

	"github.com/afex/hystrix-go/hystrix"
)

func NewPlugin() plugin.Plugin {
	return plugin.NewPlugin(
		plugin.WithName("hystrix"),
		plugin.WithHandler(hystrixHandler),
	)
}

func hystrixHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Method + "-" + r.RequestURI
		config := hystrix.CommandConfig{
			Timeout: 300,
		}
		hystrix.ConfigureCommand(name, config)

		ctx, cancel := context.WithCancel(r.Context())
		req := r.WithContext(ctx)
		if err := hystrix.Do(name,
			func() error {
				defer cancel()
				h.ServeHTTP(w, req)
				return nil
			},
			func(err error) error {
				defer cancel()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return err
			},
		); err != nil {
			log.Println("hystrix breaker err: jla, err")
			return
		}
	})
}
