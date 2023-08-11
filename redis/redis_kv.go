package redis

import (
	"context"
	"fmt"
)

type RedisKV struct {
	opts RedisOptions // redis 配置信息
}

func NewRedisKV(opts RedisOptions) *RedisKV {
	return &RedisKV{
		opts: opts,
	}
}

func (kv *RedisKV) Get(ctx context.Context, key string) (interface{}, error) {
	reply, err := NewRedisPool(kv.opts).Do(ctx, "GET", key)
	if err != nil {
		return nil, fmt.Errorf("RedisKV.Get Do error, key:%s, errmsg:%v", key, err.Error())
	}

	return reply, err
}

func (kv *RedisKV) Set(ctx context.Context, key string, value interface{}) error {

	_, err := NewRedisPool(kv.opts).Do(ctx, "SET", key, value)
	if err != nil {
		return fmt.Errorf("RedisKV.Set Do error, key:%s, errmsg:%v", key, err.Error())
	}
	return nil
}

func (kv *RedisKV) Del(ctx context.Context, key string) error {
	_, err := NewRedisPool(kv.opts).Do(ctx, "DEL", key)
	if err != nil {
		return fmt.Errorf("RedisKV.Del error, key:%s, errmsg:%v", key, err.Error())
	}
	return nil
}
