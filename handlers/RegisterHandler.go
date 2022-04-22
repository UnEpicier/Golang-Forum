package forum

import (
	"database/sql"
	f "forum"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"time"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/user/register.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page f.Page
	page.Error = ""
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true
		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		confemail := r.FormValue("confemail")
		password := r.FormValue("passwd")
		confpassword := r.FormValue("confpasswd")

		if email != confemail || password != confpassword {
			http.Redirect(w, r, "/user/register", http.StatusSeeOther)
			return
		}

		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		create := true

		var count int
		row, err := db.Query("SELECT COUNT(*) FROM user WHERE email = ?", email)
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			row.Scan(&count)
		}
		row.Close()
		if count > 0 {
			create = false
			page.Error = "Email is already used"
		}
		row, err = db.Query("SELECT COUNT(*) FROM user WHERE username = ?", username)
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			row.Scan(&count)
		}
		row.Close()
		if count > 0 {
			create = false
			page.Error = "Username is already used"
		}

		if create {
			pass, err := f.HashPassword(password)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("INSERT INTO user (`uuid`, `username`, `email`, `password`, `role`, `creation_date`, `biography`, `last_seen`) VALUES ('', ?, ?, ?, \"Member\", ?, \"\", ?)", username, email, pass, time.Now(), time.Now())
			if err != nil {
				log.Fatal(err)
			}

			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}

		db.Close()
	}

	err := tplt.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}
