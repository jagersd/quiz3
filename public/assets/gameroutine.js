connectToSocket()

const player = sessionStorage.getItem("playerSlug")
const inputSection = document.getElementById("input-section")
const waitingroom = document.getElementById("waiting-room")
const submitButton = document.getElementById("submit-answer-btn")
const warningMessage = document.getElementById("host-warning-message")
const playerResultDisplay = document.getElementById("show-player-result")
const questionCounterDisplay = document.getElementById("question-counter-display")
const questionImage = document.getElementById("question-image")
const waitingroomMessage = document.getElementById("waiting-room-tag-line")

let questionCounter = 0
let questionType = 0

function connectToSocket(){
    const submitForm = document.getElementById("answer-submit-form")
    const params = window.location.href.split("/")
    const quizId = params[params.length -1]
    const protocol = window.location.protocol === "https:" ? "wss://" : "ws://"

    const conn = new WebSocket(protocol+ document.location.host + "/ws/" + quizId)
    
    conn.onopen = () => {
        console.log("WebSocket connected!")
        conn.send(player + "|joined|" + "")
    }

    conn.onerror = (err) =>{
        console.log("Error in Socket connection", err)
    }

    conn.onclose = (event)=>{
        console.log("connection closed:", event)
        connectToSocket()
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

        if (messageContent.Attachment != ""){
            questionImage.src = messageContent.Attachment
            questionImage.style.display ="block"
        } else {
            questionImage.style.display = "none"
        }
        
        //display question number animation
        if (questionCounter >= 1){
            waitingroomMessage.innerText = "These players have yet to answer:"
            questionCounterDisplay.innerText = "Question " + questionCounter
            questionCounterDisplay.classList.add("next-question-display")
        }

        // setup the participant view section
        if (messageContent.QuestionType == 1) {
            for (const key in messageContent.Options){
                const option = document.createElement("input")
                option.setAttribute("class", "answer-option")
                option.setAttribute("type", "radio")
                option.setAttribute("name", "answer")
                option.setAttribute("id","option"+key)
                option.setAttribute("value",messageContent.Options[key])
                const label = document.createElement("label")
                label.setAttribute("for","option"+key)
                label.innerText = messageContent.Options[key]
                
                const lineBreak = document.createElement("br")

                inputSection.append(option)
                inputSection.append(label)
                inputSection.append(lineBreak)

            }
        } else {
            const openAnwer = document.createElement("input")
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
            submitButton.style.display = "none"
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
    const subtotalSection = document.getElementById("subtotals")
    let totalContent = ""
    for (const key in totalObjects){
        totalContent += key + " : " + totalObjects[key] + "<br>"
    }
    subtotalSection.innerHTML = totalContent
}

function pushCurrentResult(resultObjects){
    const currentResultSection = document.getElementById("current-question-result")
    let resultContent = ""
    const waitfor = []
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
        submitButton.style.opacity = "1"
    }
}

function setSubmitButtonState(host, started){
    if (host == sessionStorage.getItem("playerSlug")){
        if (started == true){
            warningMessage.style.display = "none"
            submitButton.value = "Next Question"
            submitButton.disabled = true
            submitButton.style.opacity = "0.6"
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

questionImage.addEventListener("click", ()=>{
    const fullImage = document.getElementById("question-image-full");
    fullImage.src = questionImage.src
    fullImage.style.display = "block"

    fullImage.addEventListener("click", ()=>{
        fullImage.style.display = "none"
    })
})
