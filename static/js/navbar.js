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