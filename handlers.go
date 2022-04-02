package forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		IndexHandler(w, r)
		return
	}

	tplt := template.Must(template.ParseFiles("./static/error.html"))

	err := tplt.Execute(w, tplt)
	if err != nil {
		log.Fatal(err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/index.html", "./static/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	err := tplt.ExecuteTemplate(w, "base", tplt)
	if err != nil {
		log.Fatal(err)
	}
}

/*
	FORUM
*/
func ForumHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/forum.html", "./static/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var forum Forum
	var categories []Category
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	row, err := db.Query("SELECT COUNT(*) FROM categories")
	if err != nil {
		log.Fatal(err)
	}
	var count int
	for row.Next() {
		err = row.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}
	row.Close()

	if count > 0 {
		row, err = db.Query("SELECT * FROM categories")
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			var category Category
			err = row.Scan(&category.Uuid, &category.Name, &category.Link)
			if err != nil {
				log.Fatal(err)
			}
			categories = append(categories, category)
		}
		row.Close()
		forum.Categories = categories
		forum.Error = ""
	} else {
		forum.Categories = []Category{}
		forum.Error = "No categories found"
	}

	db.Close()

	err = tplt.ExecuteTemplate(w, "base", forum)
	if err != nil {
		log.Fatal(err)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/post.html", "./static/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	err := tplt.ExecuteTemplate(w, "base", tplt)
	if err != nil {
		log.Fatal(err)
	}
}

/*
	ADMIN
*/

/*
	USER
*/

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/user/login.html", "./static/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}

		email := r.FormValue("email")
		password := r.FormValue("passwd")

		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		row, err := db.Query("SELECT `uuid`, `email`, `password` FROM users WHERE email = '" + email + "' LIMIT 1")
		if err != nil {
			log.Fatal(err)
		}
		var db_uuid string
		var db_email string
		var db_password string
		for row.Next() {
			err = row.Scan(&db_uuid, &db_email, &db_password)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println(db_email == email, CheckPasswordhash(password, db_password), db_email == email && CheckPasswordhash(password, db_password))

		if db_email == email && CheckPasswordhash(password, db_password) {
			cookie := http.Cookie{
				Name:    "user",
				Value:   db_uuid,
				Expires: time.Now().Add(time.Hour * 24),
			}
			http.SetCookie(w, &cookie)
			fmt.Println("Cookie set")
			http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
		} else {
			fmt.Println("Not working")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}
	}

	err := tplt.ExecuteTemplate(w, "base", tplt)
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/user/register.html", "./static/base.html"}
	tplt := template.Must(template.ParseFiles(files...))
	backError := ""

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}

		id := uuid.New().String()
		username := r.FormValue("username")
		email := r.FormValue("email")
		confemail := r.FormValue("confemail")
		password := r.FormValue("passwd")
		confpassword := r.FormValue("confpasswd")
		joined := time.Now().Format("02-01-2006")

		//In addition to the verification in JS, confirm that emails and passwords are the same
		if email != confemail || password != confpassword {
			http.Redirect(w, r, "/user/register", http.StatusSeeOther)
			return
		}

		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		create := true

		// Watch if the user isn't already registered
		var count int
		row, err := db.Query("SELECT COUNT(*) FROM users WHERE email = ?", email)
		if err != nil {
			log.Fatal(err)
		}
		row.Scan(&count)
		row.Close()
		if count > 0 {
			create = false
			backError = "Email is already used"
		}
		row, err = db.Query("SELECT COUNT(*) FROM users WHERE username = ?", username)
		if err != nil {
			log.Fatal(err)
		}
		row.Scan(&count)
		row.Close()
		if count > 0 {
			create = false
			backError = "Username is already used"
		}

		// The user is not registered, we create a new user
		if create {
			pass, err := HashPassword(password)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("INSERT INTO users (`uuid`, `username`, `email`, `password`, `role`, `joined`, `description`) VALUES (?, ?, ?, ?, \"Member\", ?, \"\")", id, username, email, pass, joined)
			if err != nil {
				log.Fatal(err)
			}

			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}

		db.Close()
	}

	err := tplt.ExecuteTemplate(w, "base", backError)
	if err != nil {
		log.Fatal(err)
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/user/profile.html", "./static/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	cookie, _ := r.Cookie("user")

	if cookie != nil {
		var user User

		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		row, err := db.Query("SELECT * FROM users WHERE uuid = ? LIMIT 1", cookie.Value)
		if err != nil {
			log.Fatal(err)
		}

		for row.Next() {
			err = row.Scan(&user.Uuid, &user.Username, &user.Email, &user.Password, &user.Role, &user.Joined, &user.Description)
			if err != nil {
				log.Fatal(err)
			}
		}
		user.Joined = strings.Replace(user.Joined, "-", "/", -1)

		row.Close()
		db.Close()

		err = tplt.ExecuteTemplate(w, "base", user)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}
