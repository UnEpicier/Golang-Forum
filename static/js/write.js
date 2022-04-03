const suggests = () => {
	const suggests = document.getElementById('pr-categories')
		.innerText
		.split('/')

	const input = document.getElementById('category-input').value
	const sContainter = document.getElementById('category-suggests')

	if (input.length > 0) {
		const s = suggests.filter(suggest => suggest.toLowerCase().startsWith(input.toLowerCase()));
		sContainter.innerHTML = ''
		s.forEach(suggest => {
			const suggestEl = document.createElement('span')
			suggestEl.innerText = suggest
			sContainter.appendChild(suggestEl)
		})
	} else {
		sContainter.innerHTML = ''
	}
}