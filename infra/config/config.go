package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/google/wire"

	"purchase/pkg/utils"
)

type DataConfig struct {
	Database struct {
		Driver string `json:"driver"`
		Source string `json:"source"`
	} `json:"database"`
	Redis struct {
		Network      string  `json:"network"`
		Addr         string  `json:"addr"`
		Password     string  `json:"password"`
		Db           int32   `json:"db"`
		DialTimeout  float64 `json:"dial_timeout"`
		ReadTimeout  float64 `json:"read_timeout"`
		WriteTimeout float64 `json:"write_timeout"`
	} `json:"redis"`
}

var ProviderSet = wire.NewSet(NewConfig)

func NewDataConfig() *DataConfig {
	return &DataConfig{}
}

func NewConfig() (config.Config, error) {
	cfgFile := getConfigFilePath()
	cfg := config.New(
		config.WithSource(
			file.NewSource(cfgFile),
			// conf.New()
		),
	)
	if err := cfg.Load(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func getConfigFilePath() string {
	root, err := utils.GetRootPath()
	if err != nil {
		panic(err)
	}
	confFile := root + "/conf/config.yaml"
	return confFile
}
