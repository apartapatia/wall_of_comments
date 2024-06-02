package redis

import (
	"errors"

	"github.com/apartapatia/wall_of_comments/internal/config"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var ErrRedisConnect = errors.New("redis connection error")

func GetRepo(cfg config.RedisConfig) (*Repo, error) {
	rc, err := newClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Repo{db: rc, validate: validator.New()}, nil
}

func newClient(cfg config.RedisConfig) (*redis.Client, error) {
	logrus.Info("connecting to redis")

	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if _, err := rc.Ping().Result(); err != nil {
		logrus.Error(ErrRedisConnect.Error())
		return nil, ErrRedisConnect
	}

	rc.FlushAll()
	return rc, nil
}
