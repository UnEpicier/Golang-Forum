const toggleTab = (id) => {
	const tabs = document.getElementsByClassName("tab");

	for (const tab of tabs) {
		if (tab.id === id) {
			tab.classList.remove("hidden");
		} else {
			tab.classList.add("hidden");
		}
	}

	const tabLinks = document.getElementsByClassName("tab-btn");
	for (const btn of tabLinks) {
		if ((btn.innerText).toLowerCase() === id) {
			btn.classList.add("active");
		} else {
			btn.classList.remove("active");
		}
	}

}

const showError = () => {
	const error = document.getElementById("error").innerText
	if (error != "") {
		const errorCategorie = error.substring(error.indexOf("[") + 1, error.indexOf("]"))
		const message = error.substring(error.indexOf("]") + 1, error.length)

		if (errorCategorie == "Username") {
			document.getElementById('username').innerText = message
			document.getElementById('username').classList.remove('hidden')
		} else if (errorCategorie == "Email") {
			document.getElementById('email').innerText = message
			document.getElementById('email').classList.remove('hidden')
		} else if (errorCategorie == "Password") {
			document.getElementById('password').innerText = message
			document.getElementById('password').classList.remove('hidden')
		}
	}
}
showError()