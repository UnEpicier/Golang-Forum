/*
 *
 ! CONSTANTS
 * 
 */
const bgColors = [
	'rgba(255, 99, 132, 0.2)',
	'rgba(255, 159, 64, 0.2)',
	'rgba(255, 205, 86, 0.2)',
	'rgba(75, 192, 192, 0.2)',
	'rgba(54, 162, 235, 0.2)',
	'rgba(153, 102, 255, 0.2)',
	'rgba(201, 203, 207, 0.2)'
];
const borderColors = [
	'rgb(255, 99, 132)',
	'rgb(255, 159, 64)',
	'rgb(255, 205, 86)',
	'rgb(75, 192, 192)',
	'rgb(54, 162, 235)',
	'rgb(153, 102, 255)',
	'rgb(201, 203, 207)'
];
const borderWidth = 1;
const hoverOffset = 2;
const legends = {
	title: {
		color: 'rgb(255, 255, 255)'
	},
	labels: {
		color: 'rgb(255, 255, 255)',
		fontColor: 'rgb(255, 255, 255)',
		fillStyle: 'rgb(255, 255, 255)'
	}
};
const scales = {
	y: {
		beginAtZero: true,
		grid: {
			borderColor: 'rgb(255, 255, 255)',
			tickColor: 'rgb(255, 255, 255)',
			color: 'rgba(255, 255, 255, .5)'
		},
		ticks: {
			color: 'rgb(255, 255, 255)'
		}
	},
	x: {
		grid: {
			borderColor: 'rgb(255, 255, 255)',
			tickColor: 'rgb(255, 255, 255)',
			color: 'rgba(255, 255, 255, .5)'
		},
		ticks: {
			color: 'rgb(255, 255, 255)'
		}
	}
};
const months = [
	'January',
	'February',
	'March',
	'April',
	'May',
	'June',
	'July',
	'August',
	'September',
	'October',
	'November',
	'December'
];


/*
 *
 ! CATEGORIES
 * 
 */
const catHTML = document.getElementsByClassName("data-category");
let cat_labels = [];
let cat_values = [];

for (let i = 0; i < catHTML.length; i++) {
	cat_labels.push(catHTML[i].children[0].innerText);
	cat_values.push(catHTML[i].children[1].innerText);
}

const cat_data = {
	labels: cat_labels,
	datasets: [{
		label: 'Categories',
		data: cat_values,
		backgroundColor: bgColors,
		borderColor: borderColors,
		borderWidth: borderWidth,
		hoverOffset: hoverOffset
	}],
};
const cat_config = {
	type: 'doughnut',
	data: cat_data,
	options: {
		plugins: {
			legend: legends
		}
	},
};

const catChart = new Chart(
	document.getElementById('categoriesChart'),
	cat_config
)

/*
 *
 ! CATEGORIES ACTIVITIES
 *
 */

const catActHTML = document.getElementsByClassName("data-catactivities");
let catAct_labels = [];
let catAct_values = [];

for (let i = 0; i < catActHTML.length; i++) {
	catAct_labels.push(catActHTML[i].children[0].innerText);
	catAct_values.push(catActHTML[i].children[1].innerText);
}

const catAct_data = {
	labels: months,
	datasets: [{
		label: 'Posts',
		data: [0, 20, 40, 5, 60, 55, 30, 20, 23, 10, 8, 3],
		backgroundColor: borderColors[0],
		borderColor: borderColors[0],
		borderWidth: borderWidth,
		hoverOffset: hoverOffset
	},
	{
		label: 'Posts',
		data: [10, 40, 2, 77, 4, 30, 10, 50, 30, 44, 66],
		backgroundColor: borderColors[1],
		borderColor: borderColors[1],
		borderWidth: borderWidth,
		hoverOffset: hoverOffset
	}],
};
const catAct_config = {
	type: 'line',
	data: catAct_data,
	options: {
		scales: scales,
		plugins: {
			legend: legends,
		}
	},
};

const catActChart = new Chart(
	document.getElementById('catActChart'),
	catAct_config
)

/*
 *
 ! USER INSCRIPTIONS
 * 
 */

const inscrHTML = document.getElementsByClassName("data-inscription");

let inscr_labels = [];
let inscr_values = [];

for (let i = 0; i < inscrHTML.length; i++) {
	inscr_labels.push(inscrHTML[i].children[0].innerText);
	inscr_values.push(inscrHTML[i].children[1].innerText);
}

const inscr_data = {
	labels: inscr_labels,
	datasets: [{
		label: 'User Inscriptions',
		data: inscr_values,
		backgroundColor: bgColors,
		borderColor: borderColors,
		borderWidth: borderWidth,
		hoverOffset: hoverOffset
	}],
};
const inscr_config = {
	type: 'bar',
	data: inscr_data,
	options: {
		scales: scales,
		plugins: {
			legend: {
				display: false
			}
		}
	},
};

const inscrChart = new Chart(
	document.getElementById('inscrChart'),
	inscr_config
)





/* WHEN LOADED, REMOVE DATA DIV */
document.getElementById('data').remove();