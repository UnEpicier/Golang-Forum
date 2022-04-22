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
	if (radio.checked == true) {
		radio.checked = false
	} else {
		radio.checked = true
	}

	if (radio.name == 'like') {
		const other = radio.nextElementSibling.nextElementSibling.nextElementSibling
		if (other.checked == true) {
			other.checked = false
		}
	} else if (radio.name == 'dislike') {
		const other = radio.previousElementSibling.previousElementSibling.previousElementSibling
		if (other.checked == true) {
			other.checked = false
		}
	}

	radio.parentElement.submit()
}
