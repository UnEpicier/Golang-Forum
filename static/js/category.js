function applyFilters() {
	const filter = document.getElementById('filter').value;

	const form = document.createElement('form');
	form.method = 'POST';

	const input = document.createElement('input');
	input.type = 'hidden';
	input.name = 'filter';
	input.value = filter;

	form.appendChild(input);

	document.body.appendChild(form);
	form.submit();
	document.body.removeChild(form);
}