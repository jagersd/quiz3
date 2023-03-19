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

	"github.com/gorilla/mux"
)

func createQuiz(w http.ResponseWriter, r *http.Request){
    var subjectId uint

    cleanUp()
    
    questionAmount,_ := strconv.Atoi(r.FormValue("question-amount"))

    dbconn.DB.Model(&models.Subject{}).Where("name = ?", r.FormValue("subject-name")).Select("id").First(&subjectId)

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
    quizStates[newQuiz.QuizSlug]=&quizState{
        Host: setHost.PlayerSlug,
        Started: false,
        QuestionCounter: 0,
        LastQuestion: uint(questionAmount),
        CurrentQuestion: "",
        CurrentResult: make(map[string]uint8),
        Total: make(map[string]uint),
    }

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
    dbconn.DB.Model(&models.Question{}).Where("subject_id = ?", subjectId).Pluck("id",&questionIds)

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

func joinGame(w http.ResponseWriter, r *http.Request){
    var quiz models.Quiz
    
    // check whether quiz exists
    dbconn.DB.Model(&models.Quiz{}).Where("quiz_slug = ?", r.FormValue("quiz-code")).First(&quiz)
    if quiz.ID == 0{
        fmt.Fprintf(w, "Error")
        return
    }
    
    // check whether playername already enrolled if not, a new one is created
    var newResult models.Result

    dbconn.DB.Where(&models.Result{QuizId:quiz.ID, PlayerName:r.FormValue("player-name")}).First(&newResult)

    if newResult.ID != 0{
        fmt.Fprintf(w,newResult.PlayerSlug)
        return
    }
    // check whether the quiz already started

    if quiz.ID != 0 && quiz.Started == true{
        fmt.Fprintf(w, "Quiz already started")
        return
    }

    // Everything checked, let the player join
    newResult.QuizId = quiz.ID
    newResult.PlayerSlug = createSlug()
    newResult.PlayerName = r.FormValue("player-name")
    dbconn.DB.Create(&newResult)

    fmt.Fprintf(w, newResult.PlayerSlug)
}

func showResults(w http.ResponseWriter, r *http.Request){
    quizSlug := mux.Vars(r)["quizSlug"]
    var finalresult []models.Result
    quizId := getQuizId(quizSlug)

    dbconn.DB.Model(&models.Result{}).
    Where("quiz_id = ?", quizId).
    Select("player_name", "total", "is_host").
    Order("total desc").
    Find(&finalresult)


    templates.ExecuteTemplate(w, "finished.html", finalresult)
}




