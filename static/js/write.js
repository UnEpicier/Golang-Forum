const suggests = () => {
	const suggests = document.getElementById('pr-categories')
		.innerText
		.split('/')

	const input = document.getElementById('category-input').value
	const sContainter = document.getElementById('category-suggests')

	if (input.length > 0) {
		const s = suggests.filter(suggest => suggest.toLowerCase().startsWith(input.toLowerCase()))
		sContainter.innerHTML = ''
		s.forEach(suggest => {
			const suggestEl = document.createElement('span')
			suggestEl.innerText = suggest
			suggestEl.onclick = () => {
				document.getElementById('category-input').value = suggest
				sContainter.innerHTML = ''
			}
			sContainter.appendChild(suggestEl)

			if (input === suggest) {
				document.getElementsByClassName('btn')[0].disabled = false
				sContainter.innerHTML = ''
			} else {
				document.getElementsByClassName('btn')[0].disabled = true
			}
		})
	} else {
		sContainter.innerHTML = ''
		document.getElementsByClassName('btn')[0].disabled = true
	}
}