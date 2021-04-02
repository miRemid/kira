package cors

import (
	"net/http"

	"github.com/micro/micro/v2/plugin"
	"github.com/rs/cors"
)

func NewPlugin() plugin.Plugin {
	return plugin.NewPlugin(
		plugin.WithName("cors"),
		plugin.WithHandler(func(h http.Handler) http.Handler {
			hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.ServeHTTP(w, r)
			})
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cors.New(cors.Options{
					AllowedOrigins:   []string{"*"},
					AllowedMethods:   []string{"*"},
					AllowedHeaders:   []string{"*"},
					AllowCredentials: true,
				}).ServeHTTP(w, r, hf)
			})
		}),
	)
}
