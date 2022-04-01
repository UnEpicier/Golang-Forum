const toggleTab = (id) => {
	const tabs = document.getElementsByClassName("tab");

	for (const tab of tabs) {
		if (tab.id === id) {
			tab.classList.remove("hidden");
		} else {
			tab.classList.add("hidden");
		}
	}

	const tabLinks = document.getElementsByClassName("tab-btn");
	for (const btn of tabLinks) {
		if ((btn.innerText).toLowerCase() === id) {
			btn.classList.add("active");
		} else {
			btn.classList.remove("active");
		}
	}

}