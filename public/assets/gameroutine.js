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
    }

    conn.onerror = (err) =>{
        console.log("Error in Socket connection", err)
    }

    conn.onclose = (event)=>{
        console.log("connection closed:", event)
    }

    submitForm.onsubmit = () => {
        conn.send(player + "|" + "hallo!")
        return false
    }

    conn.onmessage = (event)=>{
        console.log("The following was received through the websocket:", event)
    }

}


