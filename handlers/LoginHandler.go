package forum

import (
	"database/sql"
	f "forum"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/user/login.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page f.Page
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true
		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}

		email := r.FormValue("email")
		password := r.FormValue("passwd")
		keepAlive := r.FormValue("keep-alive")

		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		row, err := db.Query("SELECT email, password FROM user WHERE email = ? LIMIT 1", email)
		if err != nil {
			log.Fatal(err)
		}
		var db_email string
		var db_password string
		for row.Next() {
			err = row.Scan(&db_email, &db_password)
			if err != nil {
				log.Fatal(err)
			}
		}
		row.Close()

		if db_email == email && f.CheckPasswordhash(password, db_password) {
			db_uuid := uuid.New().String()

			_, err = db.Exec("UPDATE user SET uuid = ? WHERE email = ?", db_uuid, email)
			if err != nil {
				log.Fatal(err)
			}

			cookie := http.Cookie{
				Name:       "user",
				Value:      db_uuid,
				Path:       "/",
				Domain:     "",
				Expires:    time.Time{},
				RawExpires: "",
				MaxAge:     0,
				Secure:     false,
				HttpOnly:   false,
				SameSite:   0,
				Raw:        "",
				Unparsed:   []string{},
			}
			if keepAlive == "on" {
				cookie.Expires = time.Now().AddDate(20, 0, 0)
			}
			http.SetCookie(w, &cookie)

			_, err = db.Exec("UPDATE user SET last_seen = ? WHERE email = ?", time.Now(), db_email)
			if err != nil {
				log.Fatal(err)
			}

			http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}
	}

	err := tplt.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}
