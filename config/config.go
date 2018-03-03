package config

import (
	"log"
)

type DBInfo struct {
	User   string
	Pass   string
	Host   string
	DBName string
}

type Config struct {
	Debug   bool
	DBDebug bool
	Version string
	Port    string

	DBInfo
}

func NewConfig() *Config {
	c := new(Config)

	c.DBInfo.Host = "127.0.0.1"
	c.DBInfo.User = "root"
	c.DBInfo.Pass = "root"
	c.DBInfo.DBName = "realdb"

	c.Debug = true
	c.DBDebug = false
	c.Version = "1.0.1"
	c.Port = ":5680"
	log.Printf("Debug:[%v], DBDebug:[%v], Version:[%s], DBInfo:[User:[%s], Pass:[%s], Host:[%s], DBName:[%s]].",
		c.Debug, c.DBDebug, c.Version, c.User, c.Pass, c.Host, c.DBName, )

	return c
}
