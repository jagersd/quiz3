package controllers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var templates *template.Template

type subject struct{
    Name string `json:"name"`
    Description string `json:"description"`
}

func New() http.Handler{
    templates = template.Must(templates.ParseGlob("public/*.html"))

    r := mux.NewRouter()
    //add static files
    r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/",http.FileServer(http.Dir("public/assets/"))))


    //routes
    r.HandleFunc("/", indexPage).Methods("GET")
    r.HandleFunc("/join", joinPage).Methods("GET")
    r.HandleFunc("/game/{quizSlug}", mainRoutine)
    r.HandleFunc("/create", createPage).Methods("GET")
    r.HandleFunc("/create", createQuiz).Methods("POST")
    r.HandleFunc("/getsubjects", allSubjects).Methods("GET")
    r.HandleFunc("/addsubject",addSubject).Methods("POST")
    r.HandleFunc("/addquestion",addQuestionPage).Methods("GET")
    r.HandleFunc("/addquestion",addQuestion).Methods("POST")

    //websocket endpoint
    r.HandleFunc("/ws/{quizSlug}", serveWs)

    return r
}

