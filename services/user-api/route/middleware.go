package route

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/cache/redis"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
)

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

func JwtAuth(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		if path == "/v1/user/signup" || path == "/v1/user/signin" {
			ctx.Next()
			return
		}

		header := ctx.Request.Header.Get("Authorization")
		if header == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := parseToken(header)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		res, err := authCli.Valid(token)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		if res.Expired {
			ctx.JSON(http.StatusOK, response.Response{
				Code:  response.StatusExpired,
				Error: "signin expired",
			})
			return
		}
		log.Println(res.UserRole, path)
		if allow, err := enforcer.Enforce(res.UserRole, path, method); err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, response.Response{
				Code:  response.StatusForbidden,
				Error: err.Error(),
			})
			return
		} else if !allow {
			ctx.AbortWithStatusJSON(http.StatusForbidden, response.Response{
				Code:  response.StatusForbidden,
				Error: errors.New("no permission").Error(),
			})
			return
		}

		ctx.Request.Header.Set("userid", res.UserID)
		ctx.Request.Header.Set("userRole", res.UserRole)
		ctx.Next()
	}
}

func PrintlnPath(ctx *gin.Context) {
	log.Println("Request Path = ", ctx.Request.URL.Path)
	ctx.Next()
}

func GetUserInfoFromRedis(ctx *gin.Context) {
	conn := redis.Get()
	defer conn.Close()
	userName := ctx.Param("userName")
	if userName == "" {
		ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "missing params",
		})
		return
	}
	log.Println("UserName: ", userName)
	log.Println("Check Exists")
	key := common.UserInfoKey(userName)
	exit, err := redigo.Bool(conn.Do("EXISTS", key))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
			Code:  response.StatusRedisCheck,
			Error: err.Error(),
		})
		return
	}
	if exit {
		log.Println("Exists, read from redis")
		data, err := redigo.Bytes(conn.Do("GET", key))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
				Code:  response.StatusRedisCheck,
				Error: err.Error(),
			})
			return
		} else {
			log.Println("Get Userinfo from redis")
			var user = new(pb.User)
			json.Unmarshal(data, user)
			ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
				Code: response.StatusOK,
				Data: user,
			})
			return
		}
	} else {
		log.Println("No Exists, read from database")
		ctx.Next()
	}
}
