package controllers

import (
	"fmt"
	"net/http"
	"quiz3/dbconn"
	"quiz3/models"
	"strconv"
)

func addQuestion(w http.ResponseWriter, r *http.Request){
    var question models.Question

    questionType,_ := strconv.ParseUint(r.FormValue("question-type"),10,8)

    dbconn.DB.Model(&models.Subject{Name: r.FormValue("subject-name")}).Select("id").First(&question.SubjectId)
    
    question.Body = r.FormValue("question-body")
    question.Answer = r.FormValue("question-answer")
    question.Type = uint(questionType)
    question.Attachment = r.FormValue("attachment")

    dbconn.DB.Create(&question)
    
    if questionType == 1 {
        for i:=1; i<=6; i++{
            if r.FormValue("question-option"+fmt.Sprintf("%d",i)) != ""{
                var option models.Option
                option.QuestionId = question.ID
                option.Option = r.FormValue("question-option"+fmt.Sprintf("%d",i))
                dbconn.DB.Create(&option)
            }
        }
    }

    templates.ExecuteTemplate(w,"addquestion.html", getSubjects())
}
