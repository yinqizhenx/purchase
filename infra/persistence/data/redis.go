package data

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	"purchase/pkg/utils"
)

type redisConfig struct {
	Network      string  `json:"network"`
	Addr         string  `json:"addr"`
	Password     string  `json:"password"`
	Db           int32   `json:"db"`
	DialTimeout  float64 `json:"dial_timeout"`
	ReadTimeout  float64 `json:"read_timeout"`
	WriteTimeout float64 `json:"write_timeout"`
}

func NewRedis(c config.Config) (*redis.Client, error) {
	rdbCfg := &redisConfig{}
	err := utils.DecodeConfig(c, "data.redis", rdbCfg)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:         rdbCfg.Addr,
		Password:     rdbCfg.Password,
		DB:           int(rdbCfg.Db),
		DialTimeout:  utils.Float2Duration(rdbCfg.DialTimeout),
		WriteTimeout: utils.Float2Duration(rdbCfg.WriteTimeout),
		ReadTimeout:  utils.Float2Duration(rdbCfg.ReadTimeout),
	})
	// rdb.AddHook(redisotel.TracingHook{})
	return rdb, redisotel.InstrumentTracing(rdb)
}
