package forum

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("user")
	if cookie != nil {
		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		if r.URL.Query().Has("type") {
			if r.URL.Query().Get("type") == "post" {
				_, err = db.Exec("DELETE FROM post WHERE id = ?", r.URL.Query().Get("id"))
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("DELETE FROM comment WHERE post_id = ?", r.URL.Query().Get("id"))
				if err != nil {
					log.Fatal(err)
				}
				http.Redirect(w, r, "/user/profile", http.StatusFound)
			} else if r.URL.Query().Get("type") == "comment" {
				_, err = db.Exec("DELETE FROM comment WHERE id = ?", r.URL.Query().Get("id"))
				if err != nil {
					log.Fatal(err)
				}
				http.Redirect(w, r, "/user/profile", http.StatusFound)
			} else if r.URL.Query().Get("type") == "category" {
				_, err = db.Exec("DELETE FROM category WHERE id = ?", r.URL.Query().Get("id"))
				if err != nil {
					log.Fatal(err)
				}
				http.Redirect(w, r, "/forum", http.StatusFound)
			}
		}

		db.Close()
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
