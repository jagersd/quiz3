package main

import (
	"flag"
	"log"
	"net/http"
	"quiz3/controllers"
	"quiz3/dbconn"
)

func main() {
	migrate := CheckFlags()
	dbconn.Connect(migrate)

	handler := controllers.New()

	log.Println("The quiz3 Server wil start at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func CheckFlags() bool {
	addData := flag.Bool("d", false, "Add default db data")
	dbMigrate := flag.Bool("m", false, "Migrate to new database")
	cleanDb := flag.Bool("clean", false, "Remove all current quizzes")
	flag.Parse()
	if *addData == true {
		dbconn.AddDefaults()
	}
	if *cleanDb == true {
		dbconn.CleanDb()
	}
	if *dbMigrate == true {
		return true
	}

	return false
}
