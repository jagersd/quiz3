package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"quiz3/dbconn"
	"quiz3/models"
	"strconv"
	"strings"
	"time"
)


func createQuiz(w http.ResponseWriter, r *http.Request){
    var subjectId uint
    
    questionAmount,_ := strconv.Atoi(r.FormValue("question-amount"))

    dbconn.DB.Model(&models.Subject{Name: r.FormValue("subject-name")}).Select("id").First(&subjectId)

    newQuiz := models.Quiz{
        QuizSlug: createSlug(),
        Questions: getQuestions(subjectId, questionAmount),
    }

    dbconn.DB.Create(&newQuiz)

    setHost := models.Result{
        QuizId: newQuiz.ID,
        PlayerName: r.FormValue("player-name"),
        PlayerSlug: createSlug(),
        IsHost: true,
    }
    
    dbconn.DB.Create(&setHost)

    fmt.Fprintf(w, newQuiz.QuizSlug+"|"+setHost.PlayerSlug)
}

func createSlug() string{
    rand.Seed(time.Now().UnixNano())

    var runes = []rune("abcdefghijklmnopqrstuvwxyz")
    s := make([]rune,6)
    for i := range s{
        s[i] = runes[rand.Intn(len(runes))]
    }

    return string(s)
}

func getQuestions(subjectId uint, questionAmount int) string{
    var questionIds []uint
    dbconn.DB.Model(&models.Question{ID: subjectId}).Pluck("ID", &questionIds)
    
    if len(questionIds) < questionAmount{
        questionAmount = len(questionIds)
    }

    rand.Shuffle(len(questionIds), func(i, j int) { questionIds[i], questionIds[j] = questionIds[j], questionIds[i] })


    var returnString string
    for i:=0; i<questionAmount ; i++{
        returnString += fmt.Sprintf("%d,",questionIds[i])
    }

    return strings.TrimSuffix(returnString,",")
}




