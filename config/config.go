package config

import (
	"github.com/kelseyhightower/envconfig"
)

type AOJConfig struct {
	API              string `default:"https://judgeapi.u-aizu.ac.jp"`
	DataAPI          string `default:"https://judgedat.u-aizu.ac.jp"`
	ResultBufferName string `default:"AOJ Status"`
	Mode             string `default:"release"`
	ID               string `default:""`
	RawPassword      string `default:""`
}

var conf AOJConfig

func init() {
	_ = envconfig.Process("aoj", &conf)
	// NOTE: not yet set Variable
	conf.API = "https://judgeapi.u-aizu.ac.jp"
}

func GetConfig() AOJConfig {
	return conf
}
