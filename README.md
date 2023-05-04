# Forum

> Projet réalisé par ROULLAND Roxanne, LAROUMANIE Gabriel, VASSEUR Alexis

## Sommaire
- [Project Structure](#project-structure)
- [Setup](#setup)
- [Start](#start)

## Projet Structure

```text
/
	server/						The server is here
		server.go
	static/						HTML files are here
		assets/					Contains all images or other files
			profile/			All profiles pictures with the default one
				default.png
		css/					CSS files
			admin/				CSS files for the admin pages
			user/				CSS files for the user and auth pages
			...
		error/					404 error page
			/404.html
		js/						All the JavaScript files
			user/				JS for the user and auth pages
			...
		layout/					Layout file, act like a template for every pages
			base.html
		pages/					All the HTML pages are here
			admin/				There are the admin one
			user/				Here, for the user and auth pages
			...
	handler.go					Golang file to handle and render every pages
	structures.go				Golang file to store the structures (use to send or receive data to the pages)
	utils.go					All the useful functions
```

## Setup
To setup your forum, you need first to clone this repository.
Then install sqlite3 and run in a terminal the file `forum.sql`.
It will create all your database tables.

## Start
To run the server, open a terminal and go to the project root directory.
Then run the command `go run ./server/server.go`.