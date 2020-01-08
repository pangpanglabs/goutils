package cache

import (
	"errors"
	"net/url"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

var (
	redisWaitingTime   = time.Second
	redisExpireTime    = time.Hour * 24
	redisMaxIdle       = 5
	errorInvalidScheme = errors.New("invalid Redis database URI scheme")
)

type Redis struct {
	*redis.Pool
	ExpireTime time.Duration
	Converter  Converter
}

func WithExpireTime(d time.Duration) func(*Redis) {
	return func(r *Redis) {
		r.ExpireTime = d
	}
}
func WithGobConverter() func(*Redis) {
	return func(r *Redis) {
		r.Converter = GobConverter{}
	}
}

func NewRedis(uri string, options ...func(*Redis)) *Redis {
	redis := &Redis{
		Pool: &redis.Pool{
			MaxIdle:     redisMaxIdle,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redisConnFromUri(uri)
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
		ExpireTime: redisExpireTime,
		Converter:  JsonConverter{},
	}
	for _, option := range options {
		if option != nil {
			option(redis)
		}
	}
	return redis
}

func (r *Redis) LoadOrStore(key string, value interface{}, getter func() (interface{}, error)) (loadFromCache bool, err error) {
	if err := r.getFromRedis(key, value); err == nil {
		return true, nil
	}

	v, err := getter()
	if err != nil {
		return false, err
	}

	if err := writeTo(v, value); err != nil {
		return false, err
	}
	if v != nil {
		go r.setToRedis(key, v)
	}
	return false, nil
}

func (r *Redis) Load(key string, value interface{}) (ok bool) {
	if err := r.getFromRedis(key, value); err == nil {
		return true
	}
	return false
}
func (r *Redis) Store(key string, value interface{}) {
	go r.setToRedis(key, value)
}

func (r *Redis) Delete(key string) error {
	redisConn := r.Get()
	defer redisConn.Close()

	err := redisConn.Send("DEL", key)
	if err != nil {
		logrus.WithField("key", key).WithError(err).Error("Delete From Redis Error")
	}
	return err
}

func (r *Redis) setToRedis(k string, v interface{}) {
	data, err := r.Converter.Encode(v)
	if err != nil {
		logrus.WithError(err).Error("Set To Redis Error")
		return
	}

	redisConn := r.Get()
	defer redisConn.Close()

	if err := redisConn.Send("SETEX", k, r.ExpireTime.Seconds(), data); err != nil {
		logrus.WithError(err).Error("Set To Redis Error")
	}
	logrus.WithField("key", k).Info("Set To Redis")
}
func (r *Redis) getFromRedis(k string, v interface{}) error {
	redisConn := r.Get()
	reply, err := redis.Bytes(redisConn.Do("GET", k))
	redisConn.Close()
	if err != nil {
		if err == redis.ErrNil {
			logrus.WithField("key", k).Info("Not found")
		}
		return err
	}

	return r.Converter.Decode(reply, v)
}

func redisConnFromUri(uriString string) (redis.Conn, error) {
	uri, err := url.Parse(uriString)
	if err != nil {
		return nil, err
	}

	var network string
	var host string
	var password string
	var db string

	switch uri.Scheme {
	case "redis":
		network = "tcp"
		host = uri.Host
		if uri.User != nil {
			password, _ = uri.User.Password()
		}
		if len(uri.Path) > 1 {
			db = uri.Path[1:]
		}
	case "unix":
		network = "unix"
		host = uri.Path
	default:
		return nil, errorInvalidScheme
	}

	conn, err := redis.Dial(network, host)
	if err != nil {
		return nil, err
	}

	if password != "" {
		_, err := conn.Do("AUTH", password)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	if db != "" {
		_, err := conn.Do("SELECT", db)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}
