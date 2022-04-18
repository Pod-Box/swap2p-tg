package config

import "github.com/spf13/viper"

type Config struct {
	Swap2p Swap2p `yaml:"swap2p"`
}

type Swap2p struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
	Path string `yaml:"path"`
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

func (s *Swap2p) GetPath() string {
	return s.Path
}
