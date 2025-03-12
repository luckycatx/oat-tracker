package conf

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	TrackerID    string `mapstructure:"tracker_id"`
	Interval     uint   `mapstructure:"interval"`
	MinInterval  uint   `mapstructure:"min_interval"`
	NumShard     int    `mapstructure:"num_shard"`
	PeerLifetime int64  `mapstructure:"peer_lifetime"`
}

func Load() *Config {
	var cfg = &Config{}
	viper.AutomaticEnv()

	viper.SetConfigName("cfg")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	// Path for dev usage
	viper.AddConfigPath("../conf")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatal("Error unmarshalling config: ", err)
	}

	return cfg
}
