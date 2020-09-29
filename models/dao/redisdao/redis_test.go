package redisdao

import (
	"github.com/guoruibiao/rediscustomsync/models/dao"
	"strings"
	"testing"
)

func init() {
	redisReadNode := dao.RedisNode{Host: "10.64.9.59", Port: 8379, Db: 0, Auth: ""}
	redisWriteNode := dao.RedisNode{Host: "10.138.43.130", Port: 8379, Db: 0, Auth: ""}
	config := dao.AppConfig{
		Source: redisReadNode,
		Destination: redisWriteNode,
	}
	Init(config)
}

func TestRedisPool_GetString(t *testing.T) {
	if value, err := ReadPool.GetString("name"); err != nil {
		t.Error(err)
	}else{
		t.Log(value)
	}
	if value, err := ReadPool.GetString("address"); err != nil {
		t.Error(err)
	}else{
		t.Log(value)
	}
}

func TestRedisPool_SetString(t *testing.T) {
	key   := "string"
	value := "rediscustomsync"
	if err := WritePool.SetString(key, value); err != nil {
		t.Error(err)
	}else{
		t.Log("Set Success.")
	}
}

func TestRedisPool_MGet(t *testing.T) {
	if values, err := ReadPool.MGet("name", "address"); err != nil {
		t.Error(err)
	}else{
		t.Log(values)
	}
}

func TestRedisPool_MSet(t *testing.T) {
	values := []interface{}{
		"name", "Tiger",
		"age", 26,
		"school", "DLUT",
	}

	if err := WritePool.MSet(values...); err != nil {
		t.Error(err)
	}else{
		t.Log("MSet Success.")
	}
}

func TestRedisPool_Type(t *testing.T) {
	keys := []string{
		"name",
		"listkey",
		"setkey",
		"hashkey",
		"shortvideo-busconf-all",
	}
	for _, key := range keys {
		if tp, err := ReadPool.Type(key); err != nil {
			t.Error(err)
		}else{
			t.Log(key, tp)
		}
	}
}

func TestRedisPool_Expire(t *testing.T) {
	key := "string"
	ttl := 20
	if err := WritePool.Expire(key, ttl); err != nil {
		t.Error(err)
	}else{
		t.Log("Expire Success.")
	}
}

func TestRedisPool_TTL(t *testing.T) {
	key := "name"
	if ttl, err := ReadPool.TTL(key); err != nil {
		t.Error(err)
	}else {
		t.Log("Key=", key, ", ttl=", ttl)
	}
}

func TestRedisPool_Hgetall(t *testing.T) {
	result, err := ReadPool.Hgetall("hashkey")
	if err != nil {
		t.Error(err)
	}else{
		t.Log(result)
	}
}

func TestRedisPool_Hmset(t *testing.T) {
	key := "hashkey"
	values := []interface{}{
		"name", "tiger",
		"home", "mars",
		"ratio", 1.23456789,
	}
	if err := WritePool.Hmset(key, values...); err != nil {
		t.Error(err)
	}else{
		t.Log("HMset Success.")
	}
}

func TestRedisPool_Lrange(t *testing.T) {
	if values, err := ReadPool.Lrange("listkey"); err != nil {
		t.Error(err)
	}else{
		t.Log(values)
	}
}

func TestRedisPool_Lpush(t *testing.T) {
	key := "listkey"
	values := []interface{}{
		1,2,3,
		"a", "b", "c",
		1.234567,
	}

	if err := WritePool.Lpush(key, values...); err != nil {
		t.Error(err)
	}else{
		t.Log("Lpush Success.")
	}
}

func TestRedisPool_Zrange(t *testing.T) {
	if values, err := ReadPool.Zrange("shortvideo-busconf-all"); err != nil {
		t.Error(err)
	}else{
		t.Log(values)
	}
}

func TestRedisPool_Zadd(t *testing.T) {
	key := "zsetkey"
	members := map[string]float64{
		"zhangsan": 123,
		"lisi"    : 99.99,
		"wangwu"  : 78,
	}

	if err := WritePool.Zadd(key, members); err != nil {
		t.Error(err)
	}else{
		t.Log("Zadd Success.")
	}
}

func TestRedisPool_Smembers(t *testing.T) {
	if values, err := ReadPool.Smembers("setkey"); err != nil {
		t.Error(err)
	}else{
		t.Log(values)
	}
}

func TestRedisPool_Sadd(t *testing.T) {
	key := "setkey"
	values := []interface{}{
		1,2,3,
		"a", "b", "c",
		1.234567,
	}

	if err := WritePool.Sadd(key, values...); err != nil {
		t.Error(err)
	}else {
		t.Log("Sadd Success.")
	}
}

func TestRedisPool_Pattern(t *testing.T) {
	pattern := "t*"
	if keys, err := ReadPool.Pattern(pattern); err != nil {
		t.Error(err)
	}else{
		t.Log("Pattern=", pattern, ", keys=", strings.Join(keys, ", "))
	}
}