const formData = new FormData();
const initSection = document.getElementById("init-section")

async function createGame(event){
    const createQuizForm = document.getElementById("init-game-form")
    event.preventDefault()
    
    let quizSlug = ""
    formData.append("player-name", document.getElementById("host-name").value)
    formData.append("subject-name", document.getElementById("subject-name").value)
    formData.append("question-amount", document.getElementById("question-amount").value)

    const response = await fetch("/create", {
        method: "post",
        body: formData
    })
    let responseString = ""
    responseString = await response.text()
    
    if (responseString != ""){
        responseElements = responseString.split("|")
        sessionStorage.setItem("playerSlug", responseElements[1])
        sessionStorage.setItem("playerName", document.getElementById("host-name").value)
        quizSlug = responseElements[0]

        document.getElementById("quiz-slug").innerText = quizSlug
        createQuizForm.style.display = "none"
        initSection.style.display = "block"
        document.getElementById("to-quiz-link").href="/game/"+quizSlug
    }
}

async function joinGame(event){
    event.preventDefault()
    quizCode = document.getElementById("quiz-slug").value

    formData.append("player-name", document.getElementById("player-name").value)
    formData.append("quiz-code", quizCode)

    const response = await fetch("/joingame", {
        method: "post",
        body: formData
    })
    let responseString = ""
    responseString = await response.text()

    if (responseString != "" && responseString != "Quiz already started" && responseString != "Error"){
        sessionStorage.setItem("playerSlug", responseString)
        sessionStorage.setItem("playerName", document.getElementById("player-name").value)

        initSection.style.display ="block"
        document.getElementById("to-quiz-link").href="/game/"+quizCode
        document.getElementById("join-quiz-form").style.display = "none"
    }

    if (responseString == "Quiz already started"){
        document.getElementById("already-started").style.display = "block"
    }

}
