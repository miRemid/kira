package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/micro/cli/v2"
	microClient "github.com/micro/go-micro/v2/client"
	"github.com/micro/micro/v2/plugin"

	"github.com/miRemid/kira/client"
	"github.com/miRemid/kira/common/response"
)

var (
	user *regexp.Regexp
	e    *casbin.Enforcer
)

func init() {
	user, _ = regexp.Compile("/user/(me|deleteUser|getUserList|updateUser)/*")
	e, _ = casbin.NewEnforcer("./casbin/model.conf", "./casbin/permission.csv")
}

func needJWT(uri string) bool {
	return user.MatchString(uri)
}

func parseToken(header string) (string, error) {
	split := strings.Split(header, " ")
	if len(split) != 2 {
		return "", errors.New("invalid token struct")
	}
	if split[0] != "Bearer" {
		return "", errors.New("invalid prefix")
	}
	return split[1], nil
}

func NewPlugin(mCli microClient.Client) plugin.Plugin {
	var authClient *client.AuthClient
	return plugin.NewPlugin(
		plugin.WithName("auth"),
		plugin.WithInit(func(ctx *cli.Context) error {
			authClient = client.NewAuthClient(mCli)
			return nil
		}),
		plugin.WithHandler(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 1. 判断当前请求路径是否需要进行jwt鉴定
				path := r.URL.RequestURI()
				if needJWT(path) {
					header := r.Header.Get("Authorization")
					if header == "" {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					tokenString, err := parseToken(header)
					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					res, err := authClient.Valid(tokenString)
					if err != nil || res.Expired {
						log.Println("invalid token: ", err.Error())
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					// casbin
					path := r.URL.Path
					method := r.Method
					if allow, err := e.Enforce(res.UserRole, path, method); err != nil {
						data, _ := json.Marshal(response.Response{
							Code:  response.StatusForbidden,
							Error: err.Error(),
						})
						w.WriteHeader(http.StatusForbidden)
						w.Write(data)
						return
					} else if !allow {
						data, _ := json.Marshal(response.Response{
							Code:  response.StatusForbidden,
							Error: "no permission",
						})
						w.WriteHeader(http.StatusForbidden)
						w.Write(data)
						return
					}
					r.Header.Set("userid", res.UserID)
					r.Header.Set("userRole", res.UserRole)
				}
				h.ServeHTTP(w, r)
			})
		}),
	)
}
