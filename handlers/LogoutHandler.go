package forum

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Cookie("user")
	if user != nil {
		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("UPDATE user SET uuid = '' WHERE uuid = ?", user.Value)
		if err != nil {
			log.Fatal(err)
		}
		db.Close()

		cookie := http.Cookie{
			Name:    "user",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		}
		http.SetCookie(w, &cookie)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
