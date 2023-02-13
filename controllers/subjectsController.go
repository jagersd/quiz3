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
    Questioncount int `json:"questionCount"`
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
    var err error
    var subjects []subject

    if questionAmount == 0 {
        err = dbconn.DB.Model(&models.Subject{}).Select("id,name").Find(&subjects).Error
    } else {
        err = dbconn.DB.Model(models.Question{}).
        Select("questions.subject_id AS id, subjects.name, COUNT(questions.subject_id) AS questioncount").Group("name,subject_id").
        Joins("left join subjects on subjects.id = questions.subject_id").
        Having("COUNT(subject_id) > ?", questionAmount).
        Scan(&subjects).Error
    }

    if err != nil{
        log.Println(err)
    }

    return subjects
}
