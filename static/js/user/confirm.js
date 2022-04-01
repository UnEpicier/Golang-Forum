const password = document.getElementById("passwd")
const confirm_password = document.getElementById("confpasswd");
const button = document.getElementById("button")
const error = document.getElementById('error')

confirm_password.addEventListener("input", () => {
  if (password.value != confirm_password.value) {
    button.disabled = true
    error.innerHTML = "Vos deux mots de passe sont différents"
  } else {
    button.disabled = false
    error.innerHTML = ""
  }
})