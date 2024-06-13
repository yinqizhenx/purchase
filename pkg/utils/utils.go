package utils

import (
	"os"
	"path"
	"time"

	"github.com/go-kratos/kratos/v2/config"
)

func Float2Duration(v float64) time.Duration {
	return time.Duration(v * float64(time.Second))
}

// DecodeConfigX val must be a point, if err panic
func DecodeConfigX(c config.Config, key string, val interface{}) {
	dataCfg := c.Value(key)
	err := dataCfg.Scan(val)
	if err != nil {
		panic(err)
	}
}

// DecodeConfig val must be a point,
func DecodeConfig(c config.Config, key string, val interface{}) error {
	dataCfg := c.Value(key)
	err := dataCfg.Scan(val)
	return err
}

func GetRootPath() (string, error) {
	execPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	root := path.Dir(execPath)
	return root, nil
}
