const toggleNav = () => {
	const button = document.getElementsByClassName('nav-toggler')[0]
	const navbar = document.getElementsByClassName('navbar')[0]

	if (button.classList.contains('active')) {
		button.classList.remove('active')
		navbar.classList.remove('active')
	} else {
		button.classList.add('active')
		navbar.classList.add('active')
	}
}

const setupNav = () => {
	const path = window.location.pathname
	const elements = document.getElementsByClassName("nav-element")
	for (let i = 0; i < elements.length; i++) {
		const element = elements[i]
		if (element.getAttribute("href") === path) {
			element.classList.add("active")
		}
	}
}
setupNav()