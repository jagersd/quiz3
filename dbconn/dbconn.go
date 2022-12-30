package dbconn

import (
	"fmt"
	"log"
	"quiz3/models"

	"gopkg.in/ini.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func getDsn() string {
	cfile, err := ini.Load("conf.ini")
	if err != nil {
		log.Fatal(err)
	}
	return cfile.Section("dbconn").Key("connectionString").String()
}

func Connect(migrate bool) {
	dsn := getDsn()
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        //change to Silent before production
        Logger: logger.Default.LogMode(logger.Silent),
    })
	if err != nil {
		log.Fatal("failed to connect database", err)
	} else {
		fmt.Println("connected to database!")
	}

    if migrate == true{
        fmt.Println("Executing database migration")
        connection.AutoMigrate(
            &models.Question{},
            &models.Quiz{},
            &models.Option{},
            &models.Subject{},
            &models.Result{},
            &models.Qtype{},
        )
    }

    DB = connection
}

func AddDefaults(){
    Connect(false)
    	qtypes := []models.Qtype{
		{Description: "Multiple Choice"},
		{Description: "Open Question"},
	}
	DB.Create(&qtypes)
}
