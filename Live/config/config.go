package config

import "github.com/spf13/viper"

var (
	Conf   *config
	Secret *config
)

func init() {
	ci := &confInfo{
		ConfName: "confs",
		ConfType: "yaml",
		ConfPath: "configs/conf",
	}
	v := getConf(ci)
	Conf = &config{
		viper: v,
	}

	sci := &confInfo{
		ConfName: "secrets",
		ConfType: "yaml",
		ConfPath: "configs/secret",
	}
	sv := getConf(sci)
	Secret = &config{
		viper: sv,
	}
}

type config struct {
	viper *viper.Viper
}

type confInfo struct {
	ConfName string
	ConfType string
	ConfPath string
}

func getConf(ci *confInfo) *viper.Viper {
	v := viper.New()
	v.SetConfigName(ci.ConfName)
	v.SetConfigType(ci.ConfType)
	v.AddConfigPath(ci.ConfPath)
	_ = v.ReadInConfig()
	return v
}

func (c *config) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *config) GetInt(key string) int {
	return c.viper.GetInt(key)
}
