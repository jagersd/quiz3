package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"quiz3/dbconn"
	"quiz3/models"
)


type subject struct{
    Name string `json:"name"`
    Description string `json:"description"`
}


func addSubject(w http.ResponseWriter, r *http.Request){
    
    subjectIn := models.Subject{
		Name:        r.FormValue("subject-name"),
		Description: r.FormValue("description"),
	}

    if subjectIn.Name == "" || subjectIn.Description == ""{
        return
    }

    dbconn.DB.Create(&subjectIn)

    templates.ExecuteTemplate(w, "addquestion.html", nil)
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
