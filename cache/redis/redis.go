package redis

import (
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/miRemid/kira/common"
)

var (
	pool      *redis.Pool
	redisPool sync.Pool
)

// newRedisPool : 创建redis连接池
func newRedisPool(host, pass string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   5000,
		Wait:        true,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

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
	redisPool = sync.Pool{
		New: func() interface{} {
			return newRedisPool(common.Getenv("REDIS_ADDRESS", "127.0.0.1:6379"), common.Getenv("REDIS_PASSWORD", ""))
		},
	}
}

func Get() redis.Conn {
	p := redisPool.Get().(*redis.Pool)
	res := p.Get()
	redisPool.Put(p)
	return res
}
