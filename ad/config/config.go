package config

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	ProductionEnv bool   `json:"production_env"`
	LogRoot       string `json:"log_root"`
	BindAddr      string `json:"bind_addr"`
	DBName        string `json:"db_name"`
	DBSourceAms   string `json:"db_source_ams"`
	DBSourceCad   string `json:"db_source_cad"`
	LoggerLevel   uint8  `json:"logger_level"`
	EnableOrmLog  bool   `json:"enable_orm_log"`
	EnableHttpLog bool   `json:"enable_http_log"`

	EnablePprof      bool `json:"enable_pprof"`
	EnablePrometheus bool `json:"enable_prometheus"`
}

var c config

func init() {
	c.ProductionEnv = false
	c.LogRoot = "../log/"
	c.BindAddr = "http://127.0.0.1:26880" // todo
	c.DBName = "mysql"
	c.DBSourceAms = "root:000000@(127.0.0.1:3306)/ams?charset=utf8&parseTime=True&loc=Local"
	c.DBSourceCad = "root:000000@(127.0.0.1:3306)/cad?charset=utf8&parseTime=True&loc=Local"
	c.LoggerLevel = 1
	c.EnableOrmLog = true
	c.EnableHttpLog = true

	c.EnablePprof = false
	c.EnablePrometheus = false
}

func (c *config) Init() {
	if c.ProductionEnv {
		if c.LoggerLevel < 1 {
			c.LoggerLevel = 2
		}
	}
}

func LoadConfig(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var ctmp config
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

func GetLogRoot() string {
	return c.LogRoot
}

func GetBindAddr() string {
	return c.BindAddr
}

func GetDBName() string {
	return c.DBName
}

func GetDBSourceAms() string {
	return c.DBSourceAms
}

func GetDBSourceCad() string {
	return c.DBSourceCad
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

func IsPprofEnabled() bool {
	return c.EnablePprof
}

func IsPrometheusEnabled() bool {
	return c.EnablePrometheus
}
