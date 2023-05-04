const password = document.getElementById("passwd")
const confirm_password = document.getElementById("confpasswd");
const button = document.getElementById("button")
const error = document.getElementById('error')

confirm_password.addEventListener("input", () => {
	if (password.value != confirm_password.value) {
		button.disabled = true
		error.innerHTML = "The passwords don't match"
	} else {
		button.disabled = false
		error.innerHTML = ""
	}
})


const email = document.getElementById("email")
const confirm_email = document.getElementById("confemail");
console.log(confirm_email)

confirm_email.addEventListener("input", () => {
	if (email.value != confirm_email.value) {
		button.disabled = true
		error.innerHTML = "The emails don't match"
	} else {
		button.disabled = false
		error.innerHTML = ""
	}
})