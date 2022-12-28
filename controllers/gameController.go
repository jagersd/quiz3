package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"quiz3/dbconn"
	"quiz3/models"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type quizState struct{
    Host string
    Started bool
    CurrentQuestion string
    QuestionType uint
    Options []string
    Answer string
    QuestionCounter uint
    CurrentResult map[string]uint8
    Total map[string]uint
}

var quizStates = make(map[string]quizState)

func mainRoutine(w http.ResponseWriter, r *http.Request) {
    quizSlug := mux.Vars(r)["quizSlug"]
    
    _, ok := quizStates[quizSlug]
    if !ok{
        //try to recreate quizState
        quizStates[quizSlug] = quizState{
            Host: getHost(quizSlug),
            Started: false,
            CurrentResult: make(map[string]uint8),

        }
   }

   updateQuizState(quizStates[quizSlug], quizSlug)
   
   templates.ExecuteTemplate(w,"gameroutine.html", nil)
}

func updateQuizState(quiz quizState, quizSlug string){
    quizId := getQuizId(quizSlug)
    var result []models.Result
    
    //update current question result
    toPullResult := fmt.Sprintf("result%d",quiz.QuestionCounter)
    if toPullResult == "result0"{
        toPullResult = "result1"
    }
    dbconn.DB.Where("quiz_id = ?", quizId).Select("player_name", toPullResult, "total").Find(&result) 
    
    for _,r := range result{
        quiz.CurrentResult[r.PlayerName] = r.Result1
        quiz.Total[r.PlayerName] = r.Total
    }
}

func getHost(quizSlug string) string{
    var host string
    dbconn.DB.Model(&models.Result{QuizId: getQuizId(quizSlug), IsHost: true}).Select("player_slug").First(&host)

    return host
}

func getQuizId(quizSlug string) uint{
    var quizId uint

    dbconn.DB.Model(&models.Quiz{QuizSlug: quizSlug}).Select("id").First(&quizId)
    return quizId
}

func createResponse(player string, messageType string, message string, room quizState) ([]byte, []byte){
    var playerResponse []byte
    var hostResponse []byte

    if messageType == "answer" && player == room.Host{
        moveToNextQuestion(&room)
    }

    if messageType == "answer" || messageType == "joined"{
        hostResponse,_ = json.Marshal(room)

        type participantResponse struct{
            QuestionType uint
            Options []string
            CurrentResult map[string]uint8
        }

        var responseToParticipant participantResponse
        responseToParticipant.QuestionType = room.QuestionType
        responseToParticipant.Options = room.Options
        responseToParticipant.CurrentResult = room.CurrentResult

        playerResponse,_ = json.Marshal(responseToParticipant)

        return playerResponse, hostResponse
    }

    errorMessage,_ := json.Marshal("error|no message type found")
    return errorMessage, errorMessage
}

func moveToNextQuestion(quiz *quizState){
    counter := quiz.QuestionCounter
    var quizId uint

    dbconn.DB.Model(&models.Result{}).Where("is_host = ?", true).Select("quiz_id").First(&quizId)
    
    //Register quiz as Started
    if counter == 0 {
        dbconn.DB.Model(&models.Quiz{ID:quizId}).Update("started", true)
    }

    var questionString string
    dbconn.DB.Model(&models.Quiz{ID:quizId}).Select("questions").First(&questionString)

    questionArray := strings.Split(questionString,",")

    var question models.Question
    dbconn.DB.Model(&models.Question{}).Where("id = ?", questionArray[counter]).First(&question)

    quiz.CurrentQuestion = question.Body
    quiz.Answer = question.Answer
    quiz.QuestionType = question.Type
    
    if question.Type == 1{
        var options []models.Option
        dbconn.DB.Model(&models.Option{}).Where("question_id = ?", question.ID).Find(&options)

        for _,q := range options{
            quiz.Options = append(quiz.Options, q.Option)   
        }

    }

    quiz.QuestionCounter += 1

}

func cleanUp(){
    var quizzesInMemory []string
    dbconn.DB.Model(&models.Quiz{}).Where("updated_at < ?", time.Now().Add(-time.Hour * 2)).Pluck("quiz_slug", &quizzesInMemory) 

    for _,q := range quizzesInMemory{
        if _, ok := quizStates[q]; ok{
            delete(quizStates, q)
        }
    }
}

