const createQuizForm = document.getElementById("init-game-form")
const initSection = document.getElementById("init-section")

async function createGame(event){
    event.preventDefault()
    
    let quizSlug = ""
    let formData = new FormData();
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
        quizSlug = responseElements[0]

        document.getElementById("quiz-slug").innerText = quizSlug
        createQuizForm.style.display = "none"
        initSection.style.display = "block"
        document.getElementById("to-quiz-link").href="/game/"+quizSlug
    }
}
