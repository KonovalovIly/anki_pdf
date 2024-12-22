package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/KonovalovIly/anki_pdf/database/model"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func SetupEnvironment() (*Application, *sql.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	cfg := SetupConfig()
	storage, db := SetupStorage(cfg)

	log.Println("Connecting to database")

	app := &Application{
		Config:  cfg,
		Storage: storage,
	}

	return app, db
}

func SetupConfig() Config {
	return Config{
		Addr: GetString("ADDR", ":8080"),
		Db: DbConfig{
			Addr:         GetString("DB_ADDR", "DB_ADDR"),
			MaxOpenConns: GetInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns: GetInt("DB_MAX_IDLE_CONNS", 10),
			MaxIdleTime:  GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}

func SetupStorage(cfg Config) (model.Storage, *sql.DB) {
	dbConn, err := New(cfg.Db.Addr, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleConns, cfg.Db.MaxIdleTime)
	if err != nil {
		log.Panic("Error creating database connection %v", err)
	}
	return model.NewStorage(dbConn), dbConn
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val

}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}
