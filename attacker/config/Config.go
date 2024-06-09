package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type MysqlConfig struct {
	Host   string `json:"host" toml:"host"`
	Port   int    `json:"port" toml:"port"`
	User   string `json:"user" toml:"user"`
	Passwd string `json:"password" toml:"password"`
	DbName string `json:"database" toml:"database"`
}

type Config struct {
	HttpPort    int         `json:"http_port" toml:"http_port"`
	HttpHost    string      `json:"http_host" toml:"http_host"`
	ExecuteRpc  string      `json:"execute_rpc" toml:"execute_rpc"`
	BeaconRpc   string      `json:"beacon_rpc" toml:"beacon_rpc"`
	MetricsPort int         `json:"metrics_port" toml:"metrics_port"`
	Strategy    string      `json:"strategy" toml:"strategy"`
	DbConfig    MysqlConfig `json:"mysql" toml:"mysql"`
	SwagHost    string      `json:"swag_host" toml:"swag_host"`
	RewardFile  string      `json:"reward_file" toml:"reward_file"`
}

var _cfg *Config = nil

func (conf *Config) MetricPort() int {
	return conf.MetricsPort
}

func ParseConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("get config failed", "err", err)
		panic(err)
	}
	err = toml.Unmarshal(data, &_cfg)
	// err = json.Unmarshal(data, &_cfg)
	if err != nil {
		log.Error("unmarshal config failed", "err", err)
		panic(err)
	}
	return _cfg, nil
}

func GetConfig() *Config {
	return _cfg
}

var (
	DefaultCors    = []string{"localhost"} // Default cors domain for the apis
	DefaultVhosts  = []string{"localhost"} // Default virtual hosts for the apis
	DefaultOrigins = []string{"localhost"} // Default origins for the apis
	DefaultPrefix  = ""                    // Default prefix for the apis
	DefaultModules = []string{}            // enable all module.
	//DefaultModules = []string{"time", "block", "attest"}
)

const (
	APIBatchItemLimit         = 2000
	APIBatchResponseSizeLimit = 250 * 1000 * 1000
)
