package route

import (
	"encoding/json"
	"log"
	"net/http"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/cache/redis"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
)

func PrintlnPath(ctx *gin.Context) {
	log.Println("Request Path = ", ctx.Request.URL.Path)
	ctx.Next()
}

func GetUserInfoFromRedis(ctx *gin.Context) {
	// Get Redis Conn from Conn Pool
	conn := redis.Get()
	defer conn.Close()
	// get userName from url
	userName := ctx.Param("userName")
	if userName == "" {
		ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "missing params",
		})
		return
	}
	// get redis info key
	key := common.UserInfoKey(userName)
	// check exist
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
