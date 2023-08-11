package redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type redisHash struct {
	opts RedisOptions // redis 配置信息
}

func NewRedisHash(opts RedisOptions) *redisHash {
	return &redisHash{
		opts: opts,
	}
}

type HashItem struct {
	Field string
	Value interface{}
}

func (kv *redisHash) HGet(ctx context.Context, key, field string) (*HashItem, error) {
	reply, err := NewRedisPool(kv.opts).Do(ctx, "HGET", key, field)
	if err != nil {
		return nil, fmt.Errorf("redisHash.HGET Do error, key:%s, field:%s, errmsg:%v", key, field, err.Error())
	}

	return &HashItem{field, reply}, nil
}

func (kv *redisHash) HGetAll(ctx context.Context, key string) ([]*HashItem, error) {
	reply, err := redis.Values(NewRedisPool(kv.opts).Do(ctx, "HGETALL", key))
	if err != nil {
		return nil, fmt.Errorf("redisHash.HGETALL Do error, key:%s, errmsg:%v", key, err.Error())
	}

	if len(reply)%2 != 0 {
		return nil, fmt.Errorf("redisHash.HGETALL Do error, key:%s, errmsg: reply format error", key)
	}

	var ret []*HashItem
	for i := 0; i < len(reply); i += 2 {
		if filed, ok := reply[i].([]byte); !ok {
			return nil, fmt.Errorf("redisHash.HGETALL Do error, key:%s, errmsg: reply format error", key)
		} else {
			ret = append(ret, &HashItem{string(filed), reply[i+1]})
		}
	}

	return ret, nil
}

func (kv *redisHash) HSet(ctx context.Context, key, field string, value interface{}) error {
	_, err := NewRedisPool(kv.opts).Do(ctx, "HSET", key, field, value)
	if err != nil {
		return fmt.Errorf("redisHash.HSET Do error, key:%s, field:%s, errmsg:%v", key, field, err.Error())
	}
	return nil
}

func (kv *redisHash) HDel(ctx context.Context, key string, fields []string) error {
	var args []interface{}
	args = append(args, key)
	for _, field := range fields {
		args = append(args, field)
	}
	_, err := NewRedisPool(kv.opts).Do(ctx, "HDEL", args...)
	if err != nil {
		return fmt.Errorf("redisHash.HDEL error, key:%s, fields:%v, errmsg:%v", key, fields, err.Error())
	}
	return nil
}
