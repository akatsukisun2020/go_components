package redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type RedisZset struct {
	opts RedisOptions // redis 配置信息
}

func NewRedisZset(opts RedisOptions) *RedisZset {
	return &RedisZset{
		opts: opts,
	}
}

type ZsetItem struct {
	Score int64
	Value interface{}
}

// ZAdd zset新增元素
func (kv *RedisZset) ZAdd(ctx context.Context, key string, items []*ZsetItem) error {
	if len(items) == 0 {
		return fmt.Errorf("redisZset.ZADD Do error, key:%s, items is empty", key)
	}

	var args []interface{}
	args = append(args, key)
	for _, item := range items {
		args = append(args, item.Score, item.Value)
	}

	_, err := NewRedisPool(kv.opts).Do(ctx, "ZADD", args...)
	if err != nil {
		return fmt.Errorf("redisHash.ZADD Do error, key:%s, errmsg:%v", key, err.Error())
	}

	return nil
}

// ZRange 按照"索引范围"返回，查询zset中的成员
func (kv *RedisZset) ZRange(ctx context.Context, key string, start, end int64) ([]*ZsetItem, error) {
	reply, err := redis.Values(NewRedisPool(kv.opts).Do(ctx, "ZRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		return nil, fmt.Errorf("redisHash.ZRange Do error, key:%s, errmsg:%v", key, err.Error())
	}

	if len(reply)%2 != 0 {
		return nil, fmt.Errorf("redisHash.ZRange Do error, key:%s, errmsg: reply format error", key)
	}

	var ret []*ZsetItem
	for i := 0; i < len(reply); i += 2 {
		if filed, ok := reply[i].(int64); !ok {
			return nil, fmt.Errorf("redisHash.ZRange Do error, key:%s, errmsg: reply format error", key)
		} else {
			ret = append(ret, &ZsetItem{filed, reply[i+1]})
		}
	}

	return ret, nil
}

// ZRangebyscore 按照"分数范围"返回，查询zset中的成员
func (kv *RedisZset) ZRangebyscore(ctx context.Context, key string, minScore, maxScore int64) ([]*ZsetItem, error) {
	reply, err := redis.Values(NewRedisPool(kv.opts).Do(ctx, "ZRANGEBYSCORE", key, minScore, maxScore, "WITHSCORES"))
	if err != nil {
		return nil, fmt.Errorf("redisHash.ZRangebyscore Do error, key:%s, errmsg:%v", key, err.Error())
	}

	if len(reply)%2 != 0 {
		return nil, fmt.Errorf("redisHash.ZRangebyscore Do error, key:%s, errmsg: reply format error", key)
	}

	var ret []*ZsetItem
	for i := 0; i < len(reply); i += 2 {
		if filed, ok := reply[i].(int64); !ok {
			return nil, fmt.Errorf("redisHash.ZRangebyscore Do error, key:%s, errmsg: reply format error", key)
		} else {
			ret = append(ret, &ZsetItem{filed, reply[i+1]})
		}
	}

	return ret, nil
}

// ZRem 删除zset指定成员
func (kv *RedisZset) ZRem(ctx context.Context, key string, values []interface{}) error {
	var args []interface{}
	args = append(args, key)
	for _, v := range values {
		args = append(args, v)
	}
	_, err := NewRedisPool(kv.opts).Do(ctx, "ZREM", args...)
	if err != nil {
		return fmt.Errorf("redisHash.ZRem Do error, key:%s, errmsg:%v", key, err.Error())
	}

	return nil
}
