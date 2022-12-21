package controllers

import(
    "net/http"
)

func indexPage(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w,"index.html",nil)
}

func addQuestionPage(w http.ResponseWriter, r *http.Request){
    templates.ExecuteTemplate(w,"addquestion.html", getSubjects())
}

func joinPage(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w,"joinGame.html",nil)
}

func createPage(w http.ResponseWriter, r *http.Request){
    templates.ExecuteTemplate(w,"create.html", getSubjects())
}

func mainRoutine(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w,"gameroutine.html",nil)
}




