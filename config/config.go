package config

import "github.com/spf13/viper"

type Config struct {
	Token  string `yaml:"token"`
	Swap2p Swap2p `yaml:"swap2p"`
}

type Swap2p struct {
	Host         string `yaml:"host"`
	RedirectHost string `yaml:"redirectHost"`
}

func ReadConfig(path string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (s *Swap2p) GetHost() string {
	return s.Host
}

func (s *Swap2p) GetRedirectHost() string {
	return s.RedirectHost
}
