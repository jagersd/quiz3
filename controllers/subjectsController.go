package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"quiz3/dbconn"
	"quiz3/models"
)


func addSubject(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "application/json")
    
    var input subject

    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        log.Println("There was an error decoding the request body into the struct", err)
    }

    subjectIn := models.Subject{
		Name:        input.Name,
		Description: input.Description,
	}

    dbconn.DB.Create(&subjectIn)

    w.WriteHeader(http.StatusOK)
}

func allSubjects(w http.ResponseWriter, r *http.Request){
    var subjects []subject
    err := dbconn.DB.Model(models.Subject{}).Select("Name","Description").Find(&subjects).Error
    if err != nil{
        log.Println("Could not find subjects", err)
    } else {
        json.NewEncoder(w).Encode(subjects)
    }
}

func getSubjects() []subject{
    var subjects []subject
    err := dbconn.DB.Model(models.Subject{}).Find(&subjects).Error
    if err != nil{
        log.Println(err)
    }
    return subjects
}
