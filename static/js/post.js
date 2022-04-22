function autoGrow() {
	var textarea = event.target;
	textarea.style.height = "5px";
	textarea.style.height = (textarea.scrollHeight) + "px";
}

if (document.getElementById('comment-textarea')) {
	document.getElementById('comment-textarea').style.height = "100px"
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

function showPopup() {
	document.getElementById("popup").classList.remove("hidden")
}

function closePopup() {
	document.getElementById("popup").classList.add("hidden")
}