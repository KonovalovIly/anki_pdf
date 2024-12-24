package route

import (
	"github.com/KonovalovIly/anki_pdf/database/storage"
)

type Application struct {
	Config  Config
	Storage storage.Storage
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
