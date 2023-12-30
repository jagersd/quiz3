package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"quiz3/dbconn"
	"quiz3/models"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type quizState struct {
	Host            string
	Started         bool
	CurrentQuestion string
	Attachment      string
	QuestionType    uint
	Options         []string
	Answer          string
	QuestionCounter uint
	LastQuestion    uint
	CurrentResult   map[string]uint8
	Total           map[string]uint
}

var (
	quizStates = make(map[string]*quizState)
)

func mainRoutine(w http.ResponseWriter, r *http.Request) {
	quizSlug := mux.Vars(r)["quizSlug"]

	quiz, ok := quizStates[quizSlug]
	if !ok {
		templates.ExecuteTemplate(w, "errcatcher.htlm", "State lost for this quiz")
	}
	updateQuizState(quiz)

	templates.ExecuteTemplate(w, "gameroutine.html", quizSlug)
}

func kickPlayer(w http.ResponseWriter, r *http.Request) {
	quizSlug := r.PostFormValue("quiz-slug")
	playerName := r.PostFormValue("player-name")
	hostSlug := r.PostFormValue("host-slug")
	if quizSlug == "" || playerName == "" {
		log.Println("quizSlug: ", quizSlug, " hostSlug ", hostSlug)
		return
	}

	quiz, ok := quizStates[quizSlug]
	if !ok {
		log.Println("quizState not found")
		return
	}

	if quiz.Host != hostSlug {
		log.Println("Request does not come from host")
		return
	}

	quizId := getQuizIdByHost(hostSlug)
	if quizId == 0 {
		log.Println("quizId not found")
		return
	}

	err := dbconn.DB.Exec("DELETE FROM results WHERE quiz_id = ? AND player_name = ? LIMIT 1;", quizId, playerName).Error
	if err != nil {
		log.Println(err)
		return
	}

	delete(quiz.CurrentResult, playerName)
	delete(quiz.Total, playerName)

	updateQuizState(quiz)

	_, hostmsg := createResponse(hostSlug, "joined", "", quiz)
	x := message{hostmsg, quizSlug}
	h.sendToHost <- x

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(playerName + " removed from game."))
}

func updateQuizState(quiz *quizState) {
	quizId := getQuizIdByHost(quiz.Host)

	counter := quiz.QuestionCounter
	if counter == 0 {
		counter = 1
	}

	type currentResult struct {
		Playername string
		Current    uint8
		Total      uint
	}

	var results []currentResult

	//update current question result
	toPullResult := fmt.Sprintf("result%d AS current", counter)
	if toPullResult == "result0" {
		toPullResult = "result1"
	}

	dbconn.DB.Table("results").
		Where("quiz_id = ?", quizId).
		Select("player_name AS playername", toPullResult, "total").
		Find(&results)

	for _, r := range results {
		quiz.CurrentResult[r.Playername] = r.Current
		quiz.Total[r.Playername] = r.Total
	}
}

func getQuizIdByHost(hostSlug string) uint {
	var quizId uint

	dbconn.DB.Model(&models.Result{}).
		Where("player_slug = ? AND is_host = ?", hostSlug, true).
		Select("quiz_id").
		First(&quizId)
	return quizId
}

func getQuizId(quizSlug string) uint {
	var quizId uint

	dbconn.DB.Model(&models.Quiz{}).
		Where("quiz_slug = ?", quizSlug).
		Select("id").
		First(&quizId)
	return quizId
}

func createResponse(player string, messageType string, message string, room *quizState) ([]byte, []byte) {
	var playerResponse []byte
	var hostResponse []byte

	if messageType == "answer" && player == room.Host {
		room.moveToNextQuestion()
	}

	if messageType == "answer" || messageType == "joined" {
		if player != room.Host {
			evaluateAnswer(player, message, room)
		}
		hostResponse, _ = json.Marshal(room)

		type participantResponse struct {
			Started         bool
			QuestionCounter uint
			Attachment      string
			LastQuestion    uint
			QuestionType    uint
			Options         []string
			CurrentResult   map[string]uint8
		}

		var responseToParticipant participantResponse
		responseToParticipant.Started = room.Started
		responseToParticipant.QuestionCounter = room.QuestionCounter
		responseToParticipant.Attachment = room.Attachment
		responseToParticipant.LastQuestion = room.LastQuestion
		responseToParticipant.QuestionType = room.QuestionType
		responseToParticipant.Options = room.Options
		responseToParticipant.CurrentResult = room.CurrentResult

		playerResponse, _ = json.Marshal(responseToParticipant)

		return playerResponse, hostResponse
	}

	errorMessage, _ := json.Marshal("error|no message type found")
	return errorMessage, errorMessage
}

func (quiz *quizState) moveToNextQuestion() {
	counter := quiz.QuestionCounter

	var quizId uint

	dbconn.DB.Model(&models.Result{}).Where("player_slug = ? AND is_host = ?", quiz.Host, true).Select("quiz_id").First(&quizId)

	//Register quiz as Started
	if counter == 0 {
		dbconn.DB.Model(&models.Quiz{ID: quizId}).Update("started", true)
		quiz.Started = true
	}

	var questionString string
	dbconn.DB.Model(&models.Quiz{}).Where("id = ?", quizId).Select("questions").First(&questionString)

	questionArray := strings.Split(questionString, ",")

	if counter >= uint(len(questionArray)) {
		return
	}

	var question models.Question
	dbconn.DB.Model(&models.Question{}).Where("id = ?", questionArray[counter]).First(&question)

	quiz.CurrentQuestion = question.Body
	quiz.Answer = question.Answer
	quiz.QuestionType = question.Type

	if question.Attachment != "" {
		quiz.Attachment = "../quiz-images/" + strconv.FormatUint(uint64(question.ID), 10) + "." + question.Attachment
	} else {
		quiz.Attachment = ""
	}

	var options []models.Option
	quiz.Options = nil

	if question.Type == 1 {
		dbconn.DB.Model(&models.Option{}).Where("question_id = ?", question.ID).Find(&options)

		for _, q := range options {
			quiz.Options = append(quiz.Options, q.Option)
		}
		quiz.Options = append(quiz.Options, quiz.Answer)
	}

	counter += 1
	quiz.QuestionCounter = counter

	updateQuizState(quiz)
}

func evaluateAnswer(player string, answer string, room *quizState) {
	if answer == "" {
		return
	}

	//no evealuations needed when the quiz has not started
	if room.QuestionCounter == 0 {
		return
	}

	//get result to see whether it is already provided
	var result uint8
	column := fmt.Sprintf("result%d", room.QuestionCounter)
	dbconn.DB.Model(&models.Result{}).Where("player_slug = ?", player).Select(column).First(&result)

	if result != 0 {
		return
	}

	if equalizeAnswer(room.Answer) == equalizeAnswer(answer) {
		dbconn.DB.Model(&models.Result{}).Where("player_slug = ?", player).Update(column, 1)
		dbconn.DB.Model(&models.Result{}).Where("player_slug = ?", player).Update("total", gorm.Expr("total + ?", 1))
	} else {
		dbconn.DB.Model(&models.Result{}).Where("player_slug = ?", player).Update(column, 3)
	}

	updateQuizState(room)
}

func equalizeAnswer(in string) string {
	in = strings.ReplaceAll(in, " ", "")
	return strings.ToLower(in)
}

func cleanUp() {
	var quizzesInMemory []string
	dbconn.DB.Model(&models.Quiz{}).Where("updated_at < ?", time.Now().Add(-time.Hour*3)).Pluck("quiz_slug", &quizzesInMemory)

	for _, q := range quizzesInMemory {
		if _, ok := quizStates[q]; ok {
			delete(quizStates, q)
		}
	}
}
