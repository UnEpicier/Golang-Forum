package main

import (
	"fmt"
	forum "forum/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	// Serve files, css and js
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./static/assets/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js/"))))

	// Handle Pages templates + 404 error
	http.HandleFunc("/", forum.ErrorHandler)
	http.HandleFunc("/index", forum.IndexHandler)

	// Forum
	http.HandleFunc("/forum", forum.ForumHandler)
	http.HandleFunc("/post", forum.PostHandler)
	http.HandleFunc("/category", forum.CategoryHandler)
	http.HandleFunc("/write", forum.WriteHandler)

	// Admin
	http.HandleFunc("/admin", forum.AdminHandler)
	http.HandleFunc("/admin/delete", forum.AdminDeleteHandler)

	// User
	http.HandleFunc("/user/login", forum.LoginHandler)
	http.HandleFunc("/user/register", forum.RegisterHandler)
	http.HandleFunc("/user/profile", forum.ProfileHandler)
	http.HandleFunc("/delete", forum.DeleteHandler)
	http.HandleFunc("/user/logout", forum.LogoutHandler)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server started: http://localhost:" + port)

	err := http.ListenAndServe("localhost:"+port, nil)
	if err != nil {
		log.Fatal("Unable to start the server:\n", err)
	}
}
