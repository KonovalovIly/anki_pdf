package main

import (
	"log"
)

func main() {
	app, db := SetupEnvironment()
	defer db.Close()
	log.Fatal(app.run(app.mount()))
}
