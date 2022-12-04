package fxmodule

import (
	"github.com/go-playground/validator"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"siigo.com/order/src/api/config"
)

var CacheModule = fx.Options(
	fx.Provide(
		NewRedisClient,
	),
)

// NewRedisClient Create Connection redis
func NewRedisClient(viper *viper.Viper) *redis.Client {

	redisOtion := &config.RedisConfiguration{}

	if !viper.IsSet("redis") {
		logrus.Warnf("redis section not found.")
		client := redis.NewClient(&redis.Options{
			Addr:        "",
			Password:    "",
			DB:          0,
			ReadTimeout: 0,
		})
		return client
	}

	err := viper.UnmarshalKey("redis", redisOtion)
	if err != nil {
		panic("redis section not allowed structure")
	}

	validate := validator.New()
	err = validate.Struct(redisOtion)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:        redisOtion.Addr,
		Password:    redisOtion.Password,
		DB:          redisOtion.Db,
		ReadTimeout: redisOtion.Timeout,
	})
	pong, err := client.Ping().Result()
	logrus.Infof("Status Redis: "+pong, err)
	return client
}
