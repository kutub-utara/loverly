package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"loverly/lib/log"

	"github.com/redis/go-redis/v9"
)

type RedisCfg struct {
	Conn *redis.Client
	log  log.Interface
}

type Redis interface {
	WithCache(ctx context.Context, key string, dest interface{}, valFunc func() (interface{}, error)) error
	DelWithPattern(ctx context.Context, pattern string) error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, duration time.Duration) error
	Del(ctx context.Context, key string) error
}

func Init(ctx context.Context, log log.Interface, addr, password string) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	}

	redisClient := redis.NewClient(opts)
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		log.Error(ctx, fmt.Sprintf("init redis fail: %+v", err))
		return nil, err
	}

	return redisClient, nil
}

func InitRedis(ctx context.Context, log log.Interface, addr, password string) (Redis, error) {

	opts := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	}

	redisClient := redis.NewClient(opts)
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		log.Error(ctx, fmt.Sprintf("init redis fail: %+v", err))
		return nil, err
	}

	rediss := RedisCfg{
		Conn: redisClient,
		log:  log,
	}

	return &rediss, nil
}

func (rds *RedisCfg) WithCache(ctx context.Context, key string, dest interface{}, valFunc func() (interface{}, error)) error {
	val, err := rds.Conn.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		rds.log.Error(ctx, fmt.Sprintf("error when get data redis: %v", err))
	}

	if val != "" {
		err := json.Unmarshal([]byte(val), dest)
		if err == nil {
			return nil
		}

		rds.log.Error(ctx, fmt.Sprintf("error when unmarshal data redis:  %v", err))
	}

	data, err := valFunc()
	if err != nil {
		rds.log.Error(ctx, fmt.Sprintf("error function params:  %v", err))
		return err
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		rds.log.Error(ctx, fmt.Sprintf("error when marshal dataJSON for redis:  %v", err))
		return err
	}

	err = rds.Conn.Set(ctx, key, dataJSON, 0).Err()
	if err != nil {
		rds.log.Error(ctx, fmt.Sprintf("error when set data redis:  %v", err))
	}

	err = json.Unmarshal(dataJSON, dest)
	if err != nil {
		rds.log.Error(ctx, fmt.Sprintf("error when unmarshal data for return: %v", err))
		return err
	}

	return nil
}

func (rds *RedisCfg) Get(ctx context.Context, key string) (string, error) {
	val, err := rds.Conn.Get(ctx, key).Result()
	if err != nil {
		rds.log.Error(ctx, fmt.Sprintf("error when get data redis:  %v", err))
		return "", err
	}

	return val, nil
}

func (rds *RedisCfg) Set(ctx context.Context, key string, value string, duration time.Duration) error {

	err := rds.Conn.Set(ctx, key, value, duration).Err()
	if err != nil {
		rds.log.Error(ctx, fmt.Sprintf("error when set data redis:  %v", err))
		return err
	}

	return nil
}

func (rds *RedisCfg) Del(ctx context.Context, key string) error {
	err := rds.Conn.Del(ctx, key).Err()
	if err != nil {
		rds.log.Error(ctx, fmt.Sprintf("error when delete data redis:  %v", err))
		return err
	}

	return nil
}

func (rds *RedisCfg) DelWithPattern(ctx context.Context, pattern string) error {

	var cursor uint64
	var keys []string

	for {

		var err error
		keys, cursor, err = rds.Conn.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			rds.log.Error(ctx, fmt.Sprintf("something wrong when scan keys: %v", err))
			return err
		}

		if len(keys) == 0 {
			break
		}

		err = rds.Conn.Del(ctx, keys...).Err()
		if err != nil {
			rds.log.Error(ctx, fmt.Sprintf("something wrong when deleted data: %v", err))
			return err
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
