package config

type AOJConfig struct {
	Endpoint string `default:"https://judgeapi.u-aizu.ac.jp"`
}

func GetConfig() (AOJConfig, error) {
	return conf, nil
}
