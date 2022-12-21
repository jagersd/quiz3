function flipQuestionType(){
  let openQIdentifier = document.getElementById('openquestion').checked
  if(openQIdentifier == true){
    document.getElementById('options-section').style.display = 'none'
  } else {
    document.getElementById('options-section').style.display = 'block'
  }
}

function showAddSubjectForm(){
  document.getElementById("add-subject-form").style.display = "block"
}
