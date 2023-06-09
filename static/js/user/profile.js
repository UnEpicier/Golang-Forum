const toggleTab = (id) => {
	const tabs = document.getElementsByClassName("tab");

	for (const tab of tabs) {
		if (tab.id === id) {
			tab.classList.remove("hidden");
		} else {
			tab.classList.add("hidden");
		}
	}

	// add anchor to the current url

	const tabLinks = document.getElementsByClassName("tab-btn");
	for (const btn of tabLinks) {
		if ((btn.innerText).replace(" ", "").toLowerCase() === id) {
			btn.classList.add("active");
		} else {
			btn.classList.remove("active");
		}
	}

}

const setupTabs = () => {
	const url = window.location.href
	const anchor = url.split("#")[1]

	const tabs = document.getElementsByClassName("tab")
	for (let i = 0; i < tabs.length; i++) {
		if (tabs[i].id === anchor) {
			toggleTab(anchor)
		}
	}
}
if (document.getElementsByClassName('tabs')) {
	setupTabs()
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

const showPopup = (id) => {
	document.getElementById(id).classList.remove('hidden')
	document.body.style.overflow = "hidden"
}

const hidePopup = (id) => {
	document.getElementById(id).classList.add('hidden')
	document.body.style.overflow = "auto"
}

const autoGrow = () => {
	const element = event.target
	element.style.height = "5px"
	element.style.height = (element.scrollHeight) + "px"
}

const sendBio = () => {
	event.preventDefault()

	const bio = event.target.value

	const form = document.createElement('form')
	form.method = "POST"
	form.classList.add('hidden')

	const input = document.createElement('input')
	input.type = "hidden"
	input.name = "bio"
	input.value = bio

	form.appendChild(input)

	const f = document.createElement('input')
	f.type = "hidden"
	f.name = "form"
	f.value = "biography"

	form.appendChild(f)

	document.body.appendChild(form)

	form.submit()

	document.body.removeChild(form)
}

/*
On page ready
*/
if (window.location.pathname == "/user/profile") {
	document.getElementById('bio').style.height = "5px"
	document.getElementById('bio').style.height = (document.getElementById('bio').scrollHeight) + "px"

	showError()
}

const sendPic = () => {
	event.preventDefault()

	const pic = event.target.parentElement.submit()
}