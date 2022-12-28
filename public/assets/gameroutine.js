connectToSocket()
const player = sessionStorage.getItem("playerSlug")

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

    submitForm.onsubmit = () => {
        conn.send(player + "|answer|" + "hallo!")
        return false
    }

    conn.onmessage = (event)=>{
        parseMessage(event.data)
        console.log(event.data)
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
}

