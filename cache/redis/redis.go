package redis

import (
	"fmt"
	"time"

	// "github.com/garyburd/redigo/redis"
	"github.com/gomodule/redigo/redis"
	"github.com/miRemid/kira/common"
)

var (
	pool *redis.Pool
)

// newRedisPool : 创建redis连接池
func newRedisPool(host, pass string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 1. 打开连接
			c, err := redis.Dial("tcp", host)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			// 2. 访问认证
			if pass != "" {
				if _, err = c.Do("AUTH", pass); err != nil {
					fmt.Println(err)
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool(common.Getenv("REDIS_ADDRESS", "127.0.0.1:6379"), common.Getenv("REDIS_PASSWORD", ""))
}

func Get() redis.Conn {
	return pool.Get()
}
