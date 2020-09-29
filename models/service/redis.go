package service

import (
	"errors"
	"fmt"
	"github.com/guoruibiao/rediscustomsync/library"
	"github.com/guoruibiao/rediscustomsync/models/dao"
	"github.com/guoruibiao/rediscustomsync/models/dao/redisdao"
	"time"
)

func TransferPatterns(config dao.AppConfig)(success, failed []string, err error) {
	if len(config.Patterns) == 0 {
		return
	}

	for _, pattern := range config.Patterns {
		keys, err := redisdao.ReadPool.Pattern(pattern)
		if err != nil {
			fmt.Println("Pattern " + pattern + "transfer failed")
			continue
		}

		fmt.Printf("Pattern: %s, found keys: len=%d\n", pattern, len(keys))
		for _, key := range keys {
			fmt.Println("transfering " + key + "...")
			if transferErr := transfer(config, key); transferErr != nil {
				failed = append(failed, key)
			}else{
				success = append(success, key)
			}
			// 设置休眠时间
			time.Sleep(time.Duration(config.Interval) * time.Millisecond)
		}
	}
	return
}

func TransferKeysfile(config dao.AppConfig) (success, failed []string, err error){
	lines, err := library.Readlines(config.Keysfile)
	if err != nil {
		return
	}

	for _, line := range lines {
		fmt.Println("transfering " + line + "...")
		if transferErr := transfer(config, line); transferErr != nil {
			failed = append(failed, line)
		}else{
			success = append(success, line)
		}
		// 设置休眠时间
		time.Sleep(time.Duration(config.Interval) * time.Millisecond)
	}
	return
}

func transfer(config dao.AppConfig, key string) (err error) {
	tp , err := redisdao.ReadPool.Type(key)
	if err != nil {
		return
	}

	switch tp {
	case "string":
		if value, err := redisdao.ReadPool.GetString(key); err == nil {
			err = redisdao.WritePool.SetString(key, value)
		}
	case "list":
		if values, err := redisdao.ReadPool.Lrange(key); err == nil {
			err = redisdao.WritePool.Lpush(key,values...)
		}
	case "set":
		if values, err := redisdao.ReadPool.Smembers(key); err == nil {
			err = redisdao.WritePool.Sadd(key, values...)
		}
	case "hash":
		if members, err := redisdao.ReadPool.Hgetall(key); err == nil {
			err = redisdao.WritePool.Hmset(key, members)
		}
	case "zset":
		if members, err := redisdao.ReadPool.Zrange(key); err == nil {
			err = redisdao.WritePool.Zadd(key, members)
		}
	default:
		err = errors.New("unsupported redisdao key type")
	}

	// TTL 处理
	if config.EnableTTL {
		ttl, ttlErr := redisdao.ReadPool.TTL(key)
		if ttlErr != nil {
			err = ttlErr
			return
		}

		// 对大于 0 的 ttl 才有必要进行 expire 处理
		if ttl > 0 {
			err = redisdao.WritePool.Expire(key, ttl)
		}
	}
	return
}