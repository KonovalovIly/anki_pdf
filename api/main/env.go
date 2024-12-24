package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	route "github.com/KonovalovIly/anki_pdf/api/route"
	"github.com/KonovalovIly/anki_pdf/database/storage"
	database_utils "github.com/KonovalovIly/anki_pdf/database/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func SetupEnvironment() (*route.Application, *sql.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	cfg := SetupConfig()
	storage, db := SetupStorage(cfg)

	log.Println("Connecting to database")

	app := &route.Application{
		Config:  cfg,
		Storage: storage,
	}

	return app, db
}

func SetupConfig() route.Config {
	return route.Config{
		Addr: GetString("ADDR", ":8080"),
		Db: route.DbConfig{
			Addr:         GetString("DB_ADDR", "DB_ADDR"),
			MaxOpenConns: GetInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns: GetInt("DB_MAX_IDLE_CONNS", 10),
			MaxIdleTime:  GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}

func SetupStorage(cfg route.Config) (storage.Storage, *sql.DB) {
	dbConn, err := database_utils.New(cfg.Db.Addr, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleConns, cfg.Db.MaxIdleTime)
	if err != nil {
		log.Panic("Error creating database connection %v", err)
	}
	return storage.NewStorage(dbConn), dbConn
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
