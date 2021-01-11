package count

import (
	"fmt"
	"net/http"

	"github.com/miRemid/kira/cache/redis"
	"github.com/micro/micro/v2/plugin"
)

// this plugin will count the times of each api
// and insert into redis

func NewPlugin() plugin.Plugin {
	return plugin.NewPlugin(
		plugin.WithName("count"),
		plugin.WithHandler(countHandler),
	)
}

func countHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get api path, ip
		api := r.RequestURI
		key := fmt.Sprintf("%s-%s", api, "count")
		ip := r.RemoteAddr

		// 2.
		conn := redis.Get()
		conn.Do("setbit", key, ip, 1)
	})
}
