package config

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

type Config struct {
	DataSource []struct {
		Drive    string `json:"drive"`
		Name     string `json:"name"`
		Url      string `json:"url"`
		User     string `json:"user"`
		Password string `json:"password"`
		Db       string `json:"db"`
	} `json:"dataSource"`

	Redis struct {
		Uri      string `json:"uri"`
		Password string `json:"password"`
		Db       int    `json:"db"`
	} `json:"redis"`

	Canal struct {
		Uri         string `json:"uri"`
		Port        int    `json:"port"`
		Name        string `json:"name"`
		Password    string `json:"password"`
		DefaultDest string `json:"default_dest"`
		SoTO        int32  `json:"so_to"`
		IdleTO      int32  `json:"idle_to"`
	} `json:"canal"`

	AMQPUrl string `json:"amqp_url"`

	ElasticSearch struct {
		Host  string `json:"host"`
		Debug bool   `json:"debug"`
	} `json:"elasticSearch"`
}

var (
	JsonConfig = LoadJsonFileConfig()
)

func LoadJsonFileConfig() *Config {
	var jsonFile *os.File
	var e error

	env := os.Getenv("ACTIVE")
	dir, _ := os.Getwd()
	if env == "prd" {
		jsonFile, e = os.Open(path.Join(dir, "config-prd.json"))
	} else {
		jsonFile, e = os.Open(path.Join(dir, "config-dev.json"))
	}
	if e != nil {
		log.Println("Can't find the file: config.json", e)
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	cfg := new(Config)
	err := decoder.Decode(cfg)
	if err != nil {
		log.Println("config.json Decoder Error:", err)
	}
	return cfg
}
