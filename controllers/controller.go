package controllers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var templates *template.Template

func New() http.Handler {
	go h.run()

	templates = template.Must(templates.ParseGlob("public/*.html"))

	r := mux.NewRouter()
	//add static files
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets/"))))
	r.PathPrefix("/quiz-images/").Handler(http.StripPrefix("/quiz-images/", http.FileServer(http.Dir("public/quiz-images/"))))

	//routes
	r.HandleFunc("/", indexPage).Methods("GET")
	r.HandleFunc("/join", joinPage).Methods("GET")
	r.HandleFunc("/joingame", joinGame).Methods("POST")
	r.HandleFunc("/game/{quizSlug}", mainRoutine).Methods("GET")
	r.HandleFunc("/finished/{quizSlug}", showResults).Methods("GET")
	r.HandleFunc("/create", createPage).Methods("GET")
	r.HandleFunc("/create", createQuiz).Methods("POST")
	r.HandleFunc("/getsubjects", allSubjects).Methods("GET")
	r.HandleFunc("/addsubject", addSubject).Methods("POST")
	r.HandleFunc("/addquestion", addQuestionPage).Methods("GET")
	r.HandleFunc("/addquestion", addQuestion).Methods("POST")
	r.HandleFunc("/kickplayer", kickPlayer).Methods("POST")

	//websocket endpoint
	r.HandleFunc("/ws/{quizSlug}", serveWs)

	return r
}
