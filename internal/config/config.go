// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	DB DBConf `json:"db"`
}

type DBConf struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}
