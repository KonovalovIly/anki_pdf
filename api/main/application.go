package main

import (
	"github.com/KonovalovIly/anki_pdf/database/model"
)

type Application struct {
	Config  Config
	Storage model.Storage
}

type Config struct {
	Addr string
	Db   DbConfig
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
