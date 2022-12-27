package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"quiz3/dbconn"
	"quiz3/models"
	"time"

	"github.com/gorilla/mux"
)

type quizState struct{
    Host string
    Started bool
    CurrentQuestion string
    Options []string
    Answer string
    QuestionCounter int
    LastResult map[string]int
    Total map[string]int
}

var quizStates = make(map[string]quizState)

func mainRoutine(w http.ResponseWriter, r *http.Request) {
    quizSlug := mux.Vars(r)["quizSlug"]
    
    _, ok := quizStates[quizSlug]
    if !ok{   
        quizStates[quizSlug] = quizState{
            Host: getHost(quizSlug),
            Started: false,
        }
   }

   fmt.Println("active states: ", quizStates)
   templates.ExecuteTemplate(w,"gameroutine.html",quizStates[quizSlug])
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

func createResponse(player string, message string, room quizState) ([]byte, []byte){
    var playerResponse []byte
    var hostResponse []byte

    hostResponse,_ = json.Marshal(room)

    type participantResponse struct{
        Options []string
    }

    var responseToParticipant participantResponse
    responseToParticipant.Options = room.Options

    playerResponse,_ = json.Marshal(responseToParticipant)

    return playerResponse, hostResponse

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

