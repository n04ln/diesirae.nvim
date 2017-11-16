package config

import (
	"github.com/kelseyhightower/envconfig"
)

type AOJConfig struct {
	Endpoint    string `default:"https://judgeapi.u-aizu.ac.jp"`
	ID          string `default:""`
	RawPassword string `default:""`
}

var conf AOJConfig

func init() {
	_ = envconfig.Process("aoj", &conf)
	// NOTE: not yet set Variable
	conf.Endpoint = "https://judgeapi.u-aizu.ac.jp"
}

func GetConfig() AOJConfig {
	return conf
}
