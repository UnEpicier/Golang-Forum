const showCreatePopup = () => {
	document.getElementById('delete-popup').classList.remove('hidden')
	document.body.style.overflow = "hidden"
	document.getElementById('category_name').focus();
}

const hideCreatePopup = () => {
	document.getElementById('delete-popup').classList.add('hidden')
	document.body.style.overflow = "auto"
}

const deleteCategory = () => {
	const id = event.target.getAttribute('data-id')
	window.location.href = `/delete?type=category&id=${id}`
}