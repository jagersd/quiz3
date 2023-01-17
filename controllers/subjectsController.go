package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"quiz3/dbconn"
	"quiz3/models"
)


type subject struct{
    Id uint `json:"id"`
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

func getSubjects(questionAmount int) []subject{
    var subjects []subject
    var err error

    if questionAmount == 0 {
        err = dbconn.DB.Model(models.Subject{}).Find(&subjects).Error
    } else {
        var availableSubjectIds []int
        dbconn.DB.Model(models.Question{}).Select("subject_id").Group("subject_id").Having("COUNT(subject_id) > ?", questionAmount).Find(&availableSubjectIds)
        err = dbconn.DB.Find(&subjects, availableSubjectIds).Error
    }

    if err != nil{
        log.Println(err)
    }
    return subjects
}
