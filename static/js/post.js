function autoGrow() {
	var textarea = event.target;

	if (textarea.name == "reason") {
		console.log(textarea.style.height)
		if (parseInt((textarea.style.height).slice(0, -2)) >= 400) {
			textarea.style.overflow = 'auto'
		} else {
			textarea.style.height = "5px";
			textarea.style.height = (textarea.scrollHeight) + "px";
		}
	} else {
		textarea.style.height = "5px";
		textarea.style.height = (textarea.scrollHeight) + "px";
	}
}

if (document.getElementById('comment-textarea')) {
	document.getElementById('comment-textarea').style.height = "100px"
}
if (document.getElementById('report-textarea')) {
	document.getElementById('report-textarea').style.height = "100px"
}

function unselectRadio() {
	event.preventDefault()

	const radio = event.target.parentElement.previousElementSibling
	if (radio.checked == true || radio.hasAttribute('checked') == true) {
		radio.checked = false
		radio.removeAttribute('checked')
		event.target.classList.remove('liked')
	} else {
		radio.checked = true
		radio.setAttribute('checked', '')
		event.target.classList.add('liked')
	}

	if (radio.name == 'like') {
		const other = radio.nextElementSibling.nextElementSibling.nextElementSibling
		if (other.checked == true || other.hasAttribute('checked') == true) {
			other.checked = false
			other.removeAttribute('checked')
			other.nextElementSibling.children[0].classList.remove('liked')
		}
	} else if (radio.name == 'dislike') {
		const other = radio.previousElementSibling.previousElementSibling.previousElementSibling
		if (other.checked == true) {
			other.checked = false
			other.removeAttribute('checked')
			other.nextElementSibling.children[0].classList.add('liked')
		}
	}

	radio.parentElement.submit()
}

function showLoggedPopup() {
	document.getElementById("logged").classList.remove("hidden")
	document.body.style.overflow = 'hidden'
}

function closeLoggedPopup() {
	document.getElementById("logged").classList.add("hidden")
	document.body.style.overflow = 'auto'
}

/*
 * @param {string} type	=> 'post', 'comment'
 * @param {string} id => 'id' = data-id
 */
function showReportPopup(type) {
	const id = event.target.getAttribute('data-id')
	const uid = event.target.getAttribute('data-uid')

	if (type == "post" || type == "comment") {
		document.getElementById("report").classList.remove("hidden")
		document.body.style.overflow = 'hidden'

		document.getElementById("report-title").innerText = document.getElementById("report-title").innerText + ' ' + type

		const form = document.getElementById("report-form")

		const type_ = document.createElement('input')
		type_.type = 'hidden'
		type_.name = 'type'
		type_.value = type
		form.appendChild(type_)

		const id_ = document.createElement('input')
		id_.type = 'hidden'
		id_.name = 'id'
		id_.value = id
		form.appendChild(id_)

		const uid_ = document.createElement('input')
		uid_.type = 'hidden'
		uid_.name = 'uid'
		uid_.value = uid
		form.appendChild(uid_)

		if (type == "comment") {
			const idp = event.target.getAttribute('data-pid')
			const pid_ = document.createElement('input')
			pid_.type = 'hidden'
			pid_.name = 'pid'
			pid_.value = idp
			form.appendChild(pid_)
		}
	}
}

function closePopup() {
	event.target.parentElement.parentElement.classList.add("hidden")
	document.body.style.overflow = 'auto'
}