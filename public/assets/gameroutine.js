connectToSocket()
const player = sessionStorage.getItem("playerSlug")
const inputSection = document.getElementById("input-section")
const waitingroom = document.getElementById("waiting-room")
const submitButton = document.getElementById("submit-answer-btn")
let questionCounter = null
let questionType = 0


function connectToSocket(){
    const submitForm = document.getElementById("answer-submit-form")
    let conn;
    const params = window.location.href.split("/")
    const quizId = params[params.length -1]

    conn = new WebSocket("ws://"+ document.location.host + "/ws/" + quizId)
    
    conn.onopen = () => {
        console.log("WebSocket connected!")
        conn.send(player + "|joined|" + sessionStorage.getItem("playerName") + " just joined the game joined")
    }

    conn.onerror = (err) =>{
        console.log("Error in Socket connection", err)
    }

    conn.onclose = (event)=>{
        console.log("connection closed:", event)
    }

    submitForm.onsubmit = (e) => {
        e.preventDefault
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
        console.log(msgArray)
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
    if (messageContent.Started == false){
        let players = ""
        for (const key in messageContent.CurrentResult){
            players += key+"<br>"
        }
        waitingroom.innerHTML = players
        return
    } else {
        waitingroom.style.display = "none"
    }
    //populate answers form when the Host moved to a new one
    if (questionCounter == messageContent.QuestionCounter){
        return
    } else {
        submitButton.style.display = "block"
        questionCounter = messageContent.QuestionCounter
        inputSection.innerHTML=""
        questionType = messageContent.QuestionType
    }
    
    if (messageContent.QuestionType == 1) {
        for (const key in messageContent.Options){
            let option = document.createElement("input")
            option.setAttribute("type", "radio")
            option.setAttribute("name", "answer")
            option.setAttribute("id","option"+key)
            option.setAttribute("value",messageContent.Options[key])
            let label = document.createElement("label")
            label.setAttribute("for","option"+key)
            label.innerText = messageContent.Options[key]

            inputSection.append(option)
            inputSection.append(label)

        }
    } else {
        let openAnwer = document.createElement("input")
        openAnwer.setAttribute("type","text")
        openAnwer.setAttribute("name","answer")
        openAnwer.setAttribute("id","open-answer")
        inputSection.append(openAnwer)
    }

}

