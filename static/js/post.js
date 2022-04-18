function autoGrow() {
	var textarea = event.target;
	textarea.style.height = "5px";
	textarea.style.height = (textarea.scrollHeight) + "px";
}

if (document.getElementById('comment-textarea')) {
	document.getElementById('comment-textarea').style.height = "100px"
}