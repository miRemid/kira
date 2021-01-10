package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/micro/v2/plugin"

	authClient "github.com/miRemid/kira/services/auth/client"

	"github.com/casbin/casbin/v2"
)

var (
	authCli *authClient.AuthClient
	skip    = []string{"/user/signin", "/user/signup"}
)

func checkSkip(path string) bool {
	for i := 0; i < len(skip); i++ {
		if path == skip[i] {
			return true
		}
	}
	return false
}

func checkFile(path string) bool {
	splits := strings.Split(path, "/")
	return splits[1] == "file"
}

func parseToken(header string) (userID string, userRole string, err error) {
	split := strings.Split(header, " ")
	if len(split) != 2 {
		return "", "", errors.New("invalid token struct")
	}
	if split[0] != "Bearer" {
		return "", "", errors.New("invalid prefix")
	}
	res, err := authCli.Valid(split[1])
	return res.UserID, res.UserRole, err
}

func NewPlugin(enforcer *casbin.Enforcer) plugin.Plugin {
	return plugin.NewPlugin(
		plugin.WithName("auth"),
		plugin.WithInit(func(ctx *cli.Context) (err error) {
			authCli = authClient.NewAuthClient(client.DefaultClient)
			return nil
		}),
		plugin.WithHandler(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := r.URL.Path
				method := r.Method
				// 1. check route, if route is sign in or signup
				// just skip
				if checkSkip(path) || checkFile(path) {
					h.ServeHTTP(w, r)
					return
				}

				// 2. get token from header
				header := r.Header.Get("Authorization")
				if header == "" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				userID, userRole, err := parseToken(header)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(err.Error()))
					return
				}
				// TODO casbin valid
				log.Printf("UserRole, Path, Method = (%s, %s, %s)", userRole, path, method)
				if allow, err := enforcer.Enforce(userRole, path, method); err != nil {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(err.Error()))
					return
				} else if !allow {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("no permission"))
					return
				}

				r.Header.Set("userid", userID)
				r.Header.Set("userRole", userRole)
				h.ServeHTTP(w, r)
			})
		}),
	)
}
