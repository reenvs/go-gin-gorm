package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Product       string `json:"product"`
	ProductionEnv bool   `json:"production_env"`
	ImsRoot       string `json:"ims_root"`
	LogRoot       string `json:"log_root"`
	LoggerLevel   uint8  `json:"logger_level"`
	EnableOrmLog  bool   `json:"enable_orm_log"`
	EnableHttpLog bool   `json:"enable_http_log"`
	ApiBindAddr   string `json:"api_bind_addr"`
	WebTitle      string `json:"title"`
	//ImageBindAddr string `json:"image_bind_addr"`
	OmsAddr string `json:"oms_addr"`
}

var c Config

func init() {
	c.ProductionEnv = false
	c.ImsRoot = "../../"
	c.LogRoot = "../log/"
	c.LoggerLevel = 1
	c.EnableOrmLog = true
	c.EnableHttpLog = true
	c.Product = "ssp"
	c.WebTitle = "ssp系统"
}

func (c *Config) Init() {
	if c.ProductionEnv {
		if c.LoggerLevel < 1 {
			c.LoggerLevel = 2
		}
	}

	if len(c.WebTitle) == 0 {
		c.WebTitle = "ssp系统"
	}
}

func LoadConfig(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var ctmp Config
	err = json.Unmarshal(b, &ctmp)
	if err != nil {
		return err
	}

	c = ctmp

	pc := &c
	pc.Init()

	fc, err := json.MarshalIndent(c, "", "    ")
	if err == nil {
		ioutil.WriteFile(path, fc, 0644)
	}

	return nil
}

func IsProductionEnv() bool {
	return c.ProductionEnv
}

func GetImsRoot() string {
	return c.ImsRoot
}

func GetLogRoot() string {
	return c.LogRoot
}

func GetLoggerLevel() uint8 {
	return c.LoggerLevel
}

func IsOrmLogEnabled() bool {
	return c.EnableOrmLog
}

func IsHttpLogEnabled() bool {
	return c.EnableHttpLog
}

func GetApiBindAddr() string {
	return c.ApiBindAddr
}

/*
func GetImageBindAddr() string {
	return c.ImageBindAddr
}
//*/

func GetProduct() string {
	return c.Product
}

func GetTitle() string {
	return c.WebTitle
}
