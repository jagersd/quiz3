connectToSocket()

const player = sessionStorage.getItem("playerSlug")
const inputSection = document.getElementById("input-section")
const waitingroom = document.getElementById("waiting-room")
const submitButton = document.getElementById("submit-answer-btn")
const warningMessage = document.getElementById("host-warning-message")
const playerResultDisplay = document.getElementById("show-player-result")
const questionCounterDisplay = document.getElementById("question-counter-display")

let questionCounter = 0
let questionType = 0

function connectToSocket(){
    const submitForm = document.getElementById("answer-submit-form")
    let conn;
    const params = window.location.href.split("/")
    const quizId = params[params.length -1]

    conn = new WebSocket("ws://"+ document.location.host + "/ws/" + quizId)
    
    conn.onopen = () => {
        console.log("WebSocket connected!")
        conn.send(player + "|joined|" + "")
    }

    conn.onerror = (err) =>{
        console.log("Error in Socket connection", err)
    }

    conn.onclose = (event)=>{
        console.log("connection closed:", event)
    }

    submitForm.onsubmit = (e) => {
        e.preventDefault
        
        questionCounterDisplay.classList.remove("next-question-display")
        questionCounterDisplay.innerText = ""
        if(questionCounter == null){
            conn.send(player + "|answer|" + "lets start")
        } else {
            let answer = null 
            if (questionType == 2){
                answer = document.getElementById("open-answer").value
            }
            if (questionType == 1){
                answer = document.querySelector('input[name="answer"]:checked').value
            }
            if (answer != "" || answer != null){
                conn.send(player + "|answer|" + answer)
                submitButton.style.display = "none"
            }
        }
        
        return false
    }

    conn.onmessage = (event)=>{
        parseMessage(event.data)
    }

}

function parseMessage(message){
    let msgType = ""
    let msgContent = ""
    if (message.includes("|")){
        msgArray = message.split("|")
        msgType = msgArray[0]
        msgContent = msgArray[1]
    } else {
        msgType = "generic"
        msgContent = message
    }
    
    if (msgType == "generic"){
        pushToFront(JSON.parse(msgContent))
    }
}

function pushToFront(messageContent){
    //fill waitingroom if quiz has not started
    if (messageContent){
        let players = ""
        for (const key in messageContent.CurrentResult){
            if (messageContent.CurrentResult[key] == 0){
                players += key+"<br>"
            }
        }
        waitingroom.innerHTML = players
    }
    setSubmitButtonState(messageContent.Host, messageContent.Started)
    //populate answers form when the Host moved to a new one
    if (questionCounter != messageContent.QuestionCounter){
        submitButton.style.display = "block"
        questionCounter = messageContent.QuestionCounter
        inputSection.innerHTML=""
        questionType = messageContent.QuestionType
        
        //display question number animation
        if (questionCounter >= 1){
            questionCounterDisplay.innerText = "Question " + questionCounter
            questionCounterDisplay.classList.add("next-question-display")
        }

        // setup the participant view section
        if (messageContent.QuestionType == 1) {
            for (const key in messageContent.Options){
                let option = document.createElement("input")
                option.setAttribute("class", "answer-option")
                option.setAttribute("type", "radio")
                option.setAttribute("name", "answer")
                option.setAttribute("id","option"+key)
                option.setAttribute("value",messageContent.Options[key])
                let label = document.createElement("label")
                label.setAttribute("for","option"+key)
                label.innerText = messageContent.Options[key]
                
                let lineBreak = document.createElement("br")

                inputSection.append(option)
                inputSection.append(label)
                inputSection.append(lineBreak)

            }
        } else {
            let openAnwer = document.createElement("input")
            openAnwer.setAttribute("type","text")
            openAnwer.setAttribute("name","answer")
            openAnwer.setAttribute("id","open-answer")
            inputSection.append(openAnwer)
        }
    }
    // setupHost items
    if (sessionStorage.getItem("playerSlug") == messageContent.Host){
        const hostSection = document.getElementById("host-section")
        hostSection.style.display = "block"
        const questionBody = document.getElementById("question-body")
        questionBody.innerText = messageContent.CurrentQuestion
        document.getElementById("correct-answer").innerText = messageContent.Answer

        if (questionCounter > 0){
            pushSubtotals(messageContent.Total)
            pushCurrentResult(messageContent.CurrentResult)
        }

        if (questionCounter == messageContent.LastQuestion){
            const params = window.location.href.split("/")
            document.getElementById("quiz-finished").style.display = "block"
            document.getElementById("to-results-link").href = "/finished/"+params[params.length -1]
        }
    }

    //show player result
    if (questionCounter == messageContent.QuestionCounter && messageContent.Host != sessionStorage.getItem("playerSlug") && messageContent.Started == true){
       for (const key in messageContent.CurrentResult) {
           if (key == sessionStorage.getItem("playerName")){
               switch(messageContent.CurrentResult[key]){
                   case 1:
                       playerResultDisplay.style.display = "block"
                       document.getElementById("result-current-question").innerText = "That was correct!"
                       break;
                   case 3:
                       playerResultDisplay.style.display = "block"
                       playerResultDisplay.style.color = "red"
                       document.getElementById("result-current-question").innerText = "Hmm...no good"
                       break;
                    default:
                       playerResultDisplay.style.display = "none"
                       playerResultDisplay.style.color = "#ffe26a"
               }
           }
       }
    }
}

function pushSubtotals(totalObjects){
    let subtotalSection = document.getElementById("subtotals")
    let totalContent = ""
    for (const key in totalObjects){
        totalContent += key + " : " + totalObjects[key] + "<br>"
    }
    subtotalSection.innerHTML = totalContent
}

function pushCurrentResult(resultObjects){
    let currentResultSection = document.getElementById("current-question-result")
    let resultContent = ""
    let waitfor = []
    for (const key in resultObjects){
        switch(resultObjects[key]){
            case 3:
                resultContent += key + " : incorrect <br>"
                break;
            case 1:
                resultContent += key + " : correct <br>"
                break;
            default:
                waitfor.push(key)
                resultContent += key + " : <br>"
        }
    }
    currentResultSection.innerHTML = resultContent
    if (waitfor.length <= 1){
        submitButton.disabled = false
    }
}

function setSubmitButtonState(host, started){
    if (host == sessionStorage.getItem("playerSlug")){
        if (started == true){
            warningMessage.style.display = "none"
            submitButton.value = "Next Question"
            submitButton.disabled = true
        } else {
            warningMessage.style.display = "block"
            submitButton.value = "Start the Quiz"
        }
    }

    if (sessionStorage.getItem("playerSlug") != host){
        warningMessage.innerText = "Please wait for the host to start the game :)"
        if (started == false ){
            submitButton.style.display = "none"
        } else {
            warningMessage.style.display = "none"
        }
    }
}
