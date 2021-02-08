package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type (
	RedisClient struct {
		pool *redis.Pool
	}

	RedisConfig struct {
		Host            string
		Port            int
		Password        string
		MaxIdle         int
		MaxActive       int
		MaxConnLifetime time.Duration
		IdleTimeout     time.Duration
	}

	RedisOption interface {
		apply(o *RedisConfig)
	}

	hostOption        string
	portOption        int
	passwordOption    string
	maxIdleOption     int
	maxActiveOption   int
	maxLifetimeOption time.Duration
	idleTimeoutOption time.Duration
)

const (
	defaultHost            = "127.0.0.1"
	defaultPort            = 6379
	defaultMaxIdle         = 3
	defaultMaxActive       = 12
	defaultMaxConnLifetime = 3600 * time.Second
	defaultIdleTimeout     = 240 * time.Second
)

func NewRedisClient(options ...RedisOption) RedisClient {
	conf := &RedisConfig{
		Host:            defaultHost,
		Port:            defaultPort,
		MaxIdle:         defaultMaxIdle,
		MaxActive:       defaultMaxActive,
		MaxConnLifetime: defaultMaxConnLifetime,
		IdleTimeout:     defaultIdleTimeout,
	}
	for _, o := range options {
		o.apply(conf)
	}
	return RedisClient{pool: newPool(conf)}
}

func newPool(conf *RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:         conf.MaxIdle,
		MaxActive:       conf.MaxActive,
		MaxConnLifetime: conf.MaxConnLifetime,
		IdleTimeout:     conf.IdleTimeout,
		Dial:            func() (redis.Conn, error) { return redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port)) },
		Wait:            true,
	}
}

func (h hostOption) apply(o *RedisConfig) { o.Host = string(h) }

func (p portOption) apply(o *RedisConfig) { o.Port = int(p) }

func (p passwordOption) apply(o *RedisConfig) { o.Password = string(p) }

func (p maxIdleOption) apply(o *RedisConfig) { o.MaxIdle = int(p) }

func (p maxActiveOption) apply(o *RedisConfig) { o.MaxActive = int(p) }

func (p maxLifetimeOption) apply(o *RedisConfig) { o.MaxConnLifetime = time.Duration(p) }

func (p idleTimeoutOption) apply(o *RedisConfig) { o.IdleTimeout = time.Duration(p) }

func SetHost(h string) RedisOption { return hostOption(h) }

func SetPort(port int) RedisOption { return portOption(port) }

func SetPassword(pwd string) RedisOption { return passwordOption(pwd) }

func SetMaxIdle(idle int) RedisOption { return maxIdleOption(idle) }

func SetMaxActive(active int) RedisOption { return maxActiveOption(active) }

func SetMaxLifetime(l time.Duration) RedisOption { return maxLifetimeOption(l) }

func SetIdleTimeout(t time.Duration) RedisOption { return idleTimeoutOption(t) }

func (r RedisClient) Get(key string) (string, error) {
	return r.DoString("get")
}

func (r RedisClient) Set(key, value string) (int, error) {
	return r.DoInt("set", key, value)
}

func (r RedisClient) SetNx(key, value string) (int, error) {
	return r.DoInt("setnx", key, value)
}

func (r RedisClient) SetExpired(key, value string, expired int) (int, error) {
	return r.DoInt("set", key, value, "ex", expired)
}

func (r RedisClient) SetNxExpired(key, value string, expired int) (int, error) {
	return r.DoInt("setnx", key, value, "ex", expired)
}

func (r RedisClient) Del(key string) (int, error) {
	return r.DoInt("del", key)
}

func (r RedisClient) DoString(command string, args ...interface{}) (string, error) {
	conn := r.getConn()
	defer r.closeConn(conn)
	return redis.String(conn.Do(command, args...))
}

func (r RedisClient) DoInt(command string, args ...interface{}) (int, error) {
	conn := r.getConn()
	defer r.closeConn(conn)
	return redis.Int(conn.Do(command, args...))
}

func (r RedisClient) getConn() redis.Conn {
	return r.pool.Get()
}

func (r RedisClient) closeConn(conn redis.Conn) {
	_ = conn.Close()
}
