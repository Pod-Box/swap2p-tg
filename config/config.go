package config

import "github.com/spf13/viper"

type Config struct {
	Token  string `yaml:"token"`
	Swap2p Swap2p `yaml:"swap2p"`
}

type Swap2p struct {
	Port             string `yaml:"port"`
	Host             string `yaml:"host"`
	Path             string `yaml:"path"`
	GetDataPath      string `yaml:"get_data_path"`
	SetWalletPath    string `yaml:"set_wallet_path"`
	SetUserStatePath string `yaml:"set_user_state_path"`
	AllTradesPath    string `yaml:"all_trades_path"`
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

func (s *Swap2p) GetDataByChatIDPath() string {
	return s.GetDataPath
}
func (s *Swap2p) GetSetWalletPath() string {
	return s.SetWalletPath
}
func (s *Swap2p) GetSetUserStatePath() string {
	return s.SetUserStatePath
}
func (s *Swap2p) GetAllTradesPath() string {
	return s.AllTradesPath
}
