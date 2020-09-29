package redisdao

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/guoruibiao/rediscustomsync/models/dao"
	"strconv"
)
type RedisPool struct{
	*redis.Pool
}

var (
	ReadPool RedisPool
	WritePool RedisPool
)



func Init(config dao.AppConfig) {
	// 初始化 ReadPool
	readPool := &redis.Pool{
		// Other pool configuration not shown in this example.
		Dial: func () (redis.Conn, error) {
			server := fmt.Sprintf("%s:%d", config.Source.Host, config.Source.Port)
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if config.Source.Auth != "" {
				if _, err := c.Do("AUTH", config.Source.Auth); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", config.Source.Db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	ReadPool = RedisPool{
		readPool,
	}

	// 初始化 WritePool
	writePool := &redis.Pool{
		// Other pool configuration not shown in this example.
		Dial: func () (redis.Conn, error) {
			server := fmt.Sprintf("%s:%d", config.Destination.Host, config.Destination.Port)
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if config.Source.Auth != "" {
				if _, err := c.Do("AUTH", config.Destination.Auth); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", config.Destination.Db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	WritePool = RedisPool{
		writePool,
	}
}

func (pool *RedisPool) GetString(key string) (value string, err error) {
	conn := pool.Get()
	defer conn.Close()

	value, err = redis.String(conn.Do("GET", key))
	return
}

func (pool *RedisPool)SetString(key, value string)(err error) {
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, value)
	return
}

func (pool *RedisPool)MGet(keys ...string)(values[]string, err error) {
	conn := pool.Get()
	defer conn.Close()
	for _, key := range keys {
		value, err := pool.GetString(key)
		if err != nil {
			values = append(values, "")
			continue
		}
		values = append(values, value)
	}
	return
}

func (pool *RedisPool)MSet(values ...interface{})(err error) {
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("MSET", values...)
	return
}

func (pool *RedisPool)Type(key string) (value string, err error) {
	conn := pool.Get()
	defer conn.Close()

	value, err = redis.String(conn.Do("TYPE", key))
	return
}

func (pool *RedisPool) TTL(key string)(ttl int, err error) {
	conn := pool.Get()
	defer conn.Close()

	ttl, err = redis.Int(conn.Do("TTL", key))
	return
}

func (pool *RedisPool)Expire(key string, ttl int) (err error) {
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("EXPIRE", key, ttl)
	return
}

func (pool *RedisPool)Hgetall(key string)(result map[string]string, err error) {
	conn := pool.Get()
	defer conn.Close()

	result, err = redis.StringMap(conn.Do("HGETALL", key))
	return
}

func (pool *RedisPool)Hmset(key string, members ...interface{}) (err error) {
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(members)...)
	return
}

// TODO 返回值类型待定
func (pool *RedisPool)Lrange(key string)(values []interface{}, err error) {
	conn := pool.Get()
	defer conn.Close()

	values, err = redis.Values(conn.Do("LRANGE", key, 0, -1))
	return
}

func (pool *RedisPool)Lpush(key string, values ...interface{}) (err error) {
	conn := pool.Get()
	defer conn.Close()

	conn.Do("LPUSH", redis.Args{}.Add(key).AddFlat(values)...)
	return
}

func (pool *RedisPool)Zrange(key string)(result map[string]float64, err error) {
	conn := pool.Get()
	defer conn.Close()

	values, err := redis.StringMap(conn.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		return
	}
	result = make(map[string]float64)
	for member, value := range values {
		if floatValue, parseErr := strconv.ParseFloat(value, 64); parseErr != nil {
			result[member] = 0
		}else{
			result[member] = floatValue
		}
	}
	return
}

func (pool *RedisPool)Zadd(key string, members map[string]float64) (err error) {
	conn := pool.Get()
	defer conn.Close()

	values := make([]interface{}, len(members)*2)
	for member, value := range members {
		values = append(values, value, member)
	}
	_, err = conn.Do("ZADD", redis.Args{}.Add(key).AddFlat(values)...)
	return
}

// TODO 类型待定
func (pool *RedisPool)Smembers(key string) (values []interface{}, err error) {
	conn := pool.Get()
	defer conn.Close()

	values, err = redis.Values(conn.Do("SMEMBERS", key))
	return
}

func (pool *RedisPool)Sadd(key string, values ...interface{}) (err error) {
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("SADD", redis.Args{}.Add(key).AddFlat(values)...)
	return
}

func (pool *RedisPool)Pattern(pattern string)(keys []string, err error) {
	conn := pool.Get()
	defer conn.Close()

	var cursor int
	for {
		reply, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", pattern))
		if err != nil {
			break
		}
		if len(reply) != 2 {
			err = errors.New("RedisPool.Pattern reply not valid ")
			break
		}
		cursor, _ = redis.Int(reply[0], nil)
		k, _     := redis.Strings(reply[1], nil)
		if len(k) > 0 {
			keys = append(keys, k...)
		}

		if cursor == 0 {
			break
		}
	}

	return
}