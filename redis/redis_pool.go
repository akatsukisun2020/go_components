package redis

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	tcpNetWork = "tcp"
)

var gRedisPool map[string]*redis.Pool
var gLock sync.RWMutex
var gOnce sync.Once

func init() {
	gOnce.Do(func() {
		gRedisPool = make(map[string]*redis.Pool)
	})
}

type redisPoolConf struct {
	maxIdle        int           // 最大空闲连接数
	maxActive      int           // 最大活跃连接数
	idleTimeout    time.Duration // 空闲连接超时时间
	defaultTimeout time.Duration // 默认超时时间

	password string // 秘钥
	addr     string // 地址
}

var gDefaultRedisPoolConf = &redisPoolConf{
	maxIdle:        200,
	maxActive:      runtime.GOMAXPROCS(0) * 2500,
	idleTimeout:    3 * time.Minute,
	defaultTimeout: time.Second,
}

type RedisOption func(rc *redisPoolConf)
type RedisOptions []RedisOption

func WithMaxIdleOption(maxActive int) RedisOption {
	return func(rc *redisPoolConf) {
		if rc == nil {
			return
		}
		rc.maxActive = maxActive
	}
}

func WithMaxActiveOption(maxActive int) RedisOption {
	return func(rc *redisPoolConf) {
		if rc == nil {
			return
		}
		rc.maxActive = maxActive
	}
}

func WithIdleTimeoutOption(idleTimeout time.Duration) RedisOption {
	return func(rc *redisPoolConf) {
		if rc == nil {
			return
		}
		rc.idleTimeout = idleTimeout
	}
}

func WithDefaultTimeoutOption(defaultTimeout time.Duration) RedisOption {
	return func(rc *redisPoolConf) {
		if rc == nil {
			return
		}
		rc.defaultTimeout = defaultTimeout
	}
}

func WithPasswordOption(password string) RedisOption {
	return func(rc *redisPoolConf) {
		if rc == nil {
			return
		}
		rc.password = password
	}
}

func WithAddrOption(addr string) RedisOption {
	return func(rc *redisPoolConf) {
		if rc == nil {
			return
		}
		rc.addr = addr
	}
}

// redisPool redis池
type redisPool struct {
	opts RedisOptions
}

func NewRedisPool(opts RedisOptions) *redisPool {
	return &redisPool{
		opts: opts,
	}
}

func (rp *redisPool) Do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
	conn, err := rp.getRedisConn(ctx)
	defer conn.Close() // 归还连接
	if err != nil {
		return nil, err
	}

	return conn.Do(cmd, args...)
}

// getRedisConn 获取redis连接
func (rp *redisPool) getRedisConn(ctx context.Context) (redis.Conn, error) {
	rc := gDefaultRedisPoolConf
	for _, opt := range rp.opts {
		opt(rc)
	}

	gLock.RLock()
	if pool, ok := gRedisPool[rc.addr]; ok {
		conn, err := pool.GetContext(ctx)
		gLock.RUnlock()
		return conn, err
	}
	gLock.RUnlock()

	// 没有创建连接池的，则创建
	gLock.Lock()
	defer gLock.Unlock()
	if pool, ok := gRedisPool[rc.addr]; ok {
		return pool.GetContext(ctx)
	}

	pool := &redis.Pool{
		MaxIdle:     rc.maxIdle,
		MaxActive:   rc.maxActive,
		IdleTimeout: rc.idleTimeout,
		Wait:        true,
		DialContext: func(ctx context.Context) (redis.Conn, error) { // TODO:测试各种配置选项
			dialOpts := []redis.DialOption{
				redis.DialReadTimeout(rc.defaultTimeout),
				redis.DialWriteTimeout(rc.defaultTimeout),
				redis.DialConnectTimeout(rc.defaultTimeout),
				redis.DialPassword(rc.password),
			}
			return redis.DialContext(ctx, tcpNetWork, rc.addr, dialOpts...)
		},
	}
	gRedisPool[rc.addr] = pool // 放入连接池管理器中

	return pool.GetContext(ctx)
}
