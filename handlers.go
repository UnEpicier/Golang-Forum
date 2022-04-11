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

	tplt := template.Must(template.ParseFiles("./static/error/error.html"))

	err := tplt.Execute(w, tplt)
	if err != nil {
		log.Fatal(err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/index.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page Page
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true
	}

	err := tplt.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}

/*
	FORUM
*/
func ForumHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/forum.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page Page
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true
	}

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

	page.Content = forum

	err = tplt.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/post.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page Page
	page.Logged = false
	page.Error = ""

	cookie, _ := r.Cookie("user")

	if cookie != nil {
		page.Logged = true
	}

	if r.URL.Query().Has("id") {
		uuid := r.URL.Query().Get("id")

		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		uuid_list := []string{}
		row, err := db.Query("SELECT uuid FROM posts")
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			var uuid string
			err = row.Scan(&uuid)
			if err != nil {
				log.Fatal(err)
			}
			uuid_list = append(uuid_list, uuid)
		}
		row.Close()

		if contains(uuid_list, uuid) {
			var post Post
			row, err = db.Query("SELECT * FROM posts WHERE uuid = ?", uuid)
			if err != nil {
				log.Fatal(err)
			}
			var uid string
			for row.Next() {
				err = row.Scan(&post.Uuid, &post.Title, &post.Content, &post.Created, &uid, &post.Likes, &post.Dislikes, &post.Category)
				if err != nil {
					log.Fatal(err)
				}
			}
			row.Close()

			row, err = db.Query("SELECT * FROM users WHERE uuid = ?", uid)
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				var user User
				err = row.Scan(&user.Uuid, &user.Username, &user.Email, &user.Password, &user.Role, &user.Joined, &user.Description)
				if err != nil {
					log.Fatal(err)
				}
				user.Joined = strings.Replace(user.Joined, "-", "/", -1)
				fmt.Println(user.Joined)
				post.User = user
			}
			row.Close()

			row, err = db.Query("SELECT * FROM comments WHERE post = ?", post.Uuid)
			if err != nil {
				log.Fatal(err)
			}
			var comments []Comment
			for row.Next() {
				var comment Comment
				comment.User = User{}
				err = row.Scan(&comment.Uuid, &comment.Content, &comment.Created, &uid, &comment.Post, &comment.Likes, &comment.Dislikes)
				if err != nil {
					log.Fatal(err)
				}
				comments = append(comments, comment)
			}
			row.Close()

			for comment := range comments {
				row, err = db.Query("SELECT * FROM users WHERE uuid = ?", uid)
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					var user User
					err = row.Scan(&user.Uuid, &user.Username, &user.Email, &user.Password, &user.Role, &user.Joined, &user.Description)
					if err != nil {
						log.Fatal(err)
					}
					user.Joined = strings.Replace(user.Joined, "-", "/", -1)
					comments[comment].User = user
				}
				row.Close()
			}

			type c struct {
				Post     Post
				Comments []Comment
			}

			var content c
			content.Post = post
			content.Comments = comments
			page.Content = content

			db.Close()

			err = tplt.Execute(w, page)
			if err != nil {
				log.Fatal(err)
			}

		} else {
			http.Redirect(w, r, "/forum", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/category.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page Page
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true
	}

	if r.URL.Query().Has("id") {
		uuid := r.URL.Query().Get("id")

		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		uuid_list := []string{}
		row, err := db.Query("SELECT uuid FROM categories")
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			var uuid string
			err = row.Scan(&uuid)
			if err != nil {
				log.Fatal(err)
			}
			uuid_list = append(uuid_list, uuid)
		}

		if contains(uuid_list, uuid) {
			row, err = db.Query("SELECT name FROM categories WHERE uuid = ?", uuid)
			if err != nil {
				log.Fatal(err)
			}

			var name string
			for row.Next() {
				err = row.Scan(&name)
				if err != nil {
					log.Fatal(err)
				}
			}

			posts := []Post{}

			row, err = db.Query("SELECT * FROM users as u INNER JOIN posts as p ON u.uuid = p.user WHERE p.uuid = ?", uuid)
			if err != nil {
				log.Fatal(err)
			}
			// var uid string
			for row.Next() {
				var post Post
				// err = row.Scan(&post.User.Uuid, &post.User.Username, &post.User.Email, &post.User.Password, &post.User.Role, &post.User.Joined, &post.User.Description, &post.Title, &post.Content, &post.Created, &uid, &post.Likes, &post.Dislikes, &post.Category)
				err = row.Scan(&post.User.Username)
				if err != nil {
					log.Fatal(err)
				}
				posts = append(posts, post)
				fmt.Println("wsh", post)
			}

			type c struct {
				Name  string
				Posts []Post
			}

			page.Content = c{Name: name, Posts: posts}
		} else {
			http.Redirect(w, r, "/forum", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}

	err := tplt.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/write.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page Page
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true

		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		var categories []string
		row, err := db.Query("SELECT name FROM categories")
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			var name string
			err = row.Scan(&name)
			if err != nil {
				log.Fatal(err)
			}
			categories = append(categories, name)
		}
		page.Content = strings.Join(categories, "/")

		if r.Method == "POST" {
			r.ParseForm()
			category := r.FormValue("category")
			title := r.FormValue("title")
			content := strings.Replace(r.FormValue("content"), "\r\n", "<br>", -1)
			var uid string

			if title == "" || content == "" || category == "" {
				page.Error = "All fields are required"
			} else {
				row, err := db.Query("SELECT uuid, name FROM categories")
				if err != nil {
					log.Fatal(err)
				}
				var categoryName string
				found := false
				for row.Next() {
					err = row.Scan(&uid, &categoryName)
					if err != nil {
						log.Fatal(err)
					}
					if strings.EqualFold(categoryName, category) {
						category = categoryName
						found = true
						break
					}
				}
				row.Close()

				if !found {
					uid = uuid.New().String()
					_, err = db.Exec("INSERT INTO categories (uuid, name, link) VALUES (?, ?, ?)", uid, capitalize(category), "/category?id="+uid)
					if err != nil {
						log.Fatal(err)
					}
				}

				_, err = db.Exec("INSERT INTO posts (uuid, title, content, created, user, likes, dislikes, category) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", uuid.New().String(), title, content, time.Now().Format("02-01-2006"), cookie.Value, 0, 0, uid)
				if err != nil {
					log.Fatal(err)
				}

				db.Close()

				http.Redirect(w, r, "/forum", http.StatusSeeOther)
			}
		}

		err = tplt.Execute(w, page)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}

/*
	ADMIN
*/
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/admin/admin.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page Page
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true

		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		row, err := db.Query("SELECT role FROM users WHERE uuid = ?", cookie.Value)
		if err != nil {
			log.Fatal(err)
		}
		var role string
		for row.Next() {
			err = row.Scan(&role)
			if err != nil {
				log.Fatal(err)
			}
		}
		row.Close()
		if role == "Admin" {

			err := tplt.Execute(w, page)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		db.Close()
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

/*
	USER
*/

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/user/login.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page Page
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

		if db_email == email && CheckPasswordhash(password, db_password) {

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
				cookie.Expires = time.Now().AddDate(5, 0, 0)
			}
			http.SetCookie(w, &cookie)

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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/user/register.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page Page
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
		for row.Next() {
			row.Scan(&count)
		}
		row.Close()
		if count > 0 {
			create = false
			page.Error = "Email is already used"
		}
		row, err = db.Query("SELECT COUNT(*) FROM users WHERE username = ?", username)
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

	err := tplt.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/user/profile.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page Page
	page.Logged = false
	page.Error = ""

	cookie, _ := r.Cookie("user")

	if cookie != nil {
		page.Logged = true
		var user User

		user.Posts = []Post{}
		user.Comments = []Comment{}

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
		row.Close()
		user.Joined = strings.Replace(user.Joined, "-", "/", -1)

		// Posts
		var count int

		row, err = db.Query("SELECT COUNT(*) FROM posts WHERE user = ?", user.Uuid)
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			err := row.Scan(&count)
			if err != nil {
				log.Fatal(err)
			}
		}

		if count > 0 {
			row, err = db.Query("SELECT uuid, title, content, created, likes, dislikes, category FROM posts WHERE user = ?", user.Uuid)
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				var post Post
				post.User = user
				err = row.Scan(&post.Uuid, &post.Title, &post.Content, &post.Created, &post.Likes, &post.Dislikes, &post.Category)
				if err != nil {
					log.Fatal(err)
				}
				user.Posts = append(user.Posts, post)
			}
		}

		// Comments
		row, err = db.Query("SELECT COUNT(*) FROM comments WHERE user = ?", user.Uuid)
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			err := row.Scan(&count)
			if err != nil {
				log.Fatal(err)
			}
		}

		if count > 0 {
			row, err = db.Query("SELECT uuid, content, created, post, likes, dislikes FROM comments WHERE user = ?", user.Uuid)
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				var comment Comment
				comment.User = user
				err = row.Scan(&comment.Uuid, &comment.Content, &comment.Created, &comment.Post, &comment.Likes, &comment.Dislikes)
				if err != nil {
					log.Fatal(err)
				}
				user.Comments = append(user.Comments, comment)
			}
		}

		/*
			Settings
		*/
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				log.Fatal(err)
			}

			/* Username */
			if r.FormValue("form") == "username" {
				current := r.FormValue("username")
				new := r.FormValue("newusername")
				password := r.FormValue("passwd")

				var count int
				row, err := db.Query("SELECT COUNT(*) FROM users WHERE uuid = ? AND username = ?", cookie.Value, current)
				if err != nil {
					log.Fatal(err)
				}

				for row.Next() {
					err = row.Scan(&count)
					if err != nil {
						log.Fatal(err)
					}
				}

				if count > 0 {
					var db_pass string

					row, err = db.Query("SELECT password FROM users WHERE uuid = ? AND username = ?", cookie.Value, current)
					if err != nil {
						log.Fatal(err)
					}
					for row.Next() {
						err = row.Scan(&db_pass)
						if err != nil {
							log.Fatal(err)
						}
					}

					if CheckPasswordhash(password, db_pass) {
						_, err = db.Exec("UPDATE users SET username = ? WHERE uuid = ?", new, cookie.Value)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						page.Error = "[Username] Wrong password"
					}
				} else {
					page.Error = "[Username] This is not your current username"
				}
			} else if r.FormValue("form") == "email" {
				current := r.FormValue("oldemail")
				new := r.FormValue("newemail")
				confemail := r.FormValue("confemail")
				password := r.FormValue("passwd")

				if new != confemail {
					page.Error = "[Email] Emails don't match"
				} else {
					var count int
					row, err := db.Query("SELECT COUNT(*) FROM users WHERE uuid = ? AND email = ?", cookie.Value, current)
					if err != nil {
						log.Fatal(err)
					}

					for row.Next() {
						err = row.Scan(&count)
						if err != nil {
							log.Fatal(err)
						}
					}

					if count > 0 {
						var db_pass string

						row, err = db.Query("SELECT password FROM users WHERE uuid = ? AND email = ?", cookie.Value, current)
						if err != nil {
							log.Fatal(err)
						}
						for row.Next() {
							err = row.Scan(&db_pass)
							if err != nil {
								log.Fatal(err)
							}
						}

						if CheckPasswordhash(password, db_pass) {
							_, err = db.Exec("UPDATE users SET email = ? WHERE uuid = ?", new, cookie.Value)
							if err != nil {
								log.Fatal(err)
							}
						} else {
							page.Error = "[Email] Wrong password"
						}
					} else {
						page.Error = "[Email] This is not your current mail address"
					}
				}
			} else if r.FormValue("form") == "password" {
				current := r.FormValue("oldpswd")
				new := r.FormValue("newpswd")
				conf := r.FormValue("confpswd")

				if new != conf {
					page.Error = "[Password] Passwords don't match"
				} else {
					var db_pass string
					row, err := db.Query("SELECT password FROM users WHERE uuid = ?", cookie.Value)
					if err != nil {
						log.Fatal(err)
					}

					for row.Next() {
						err = row.Scan(&db_pass)
						if err != nil {
							log.Fatal(err)
						}
					}

					if CheckPasswordhash(current, db_pass) {
						hash, err := HashPassword(new)
						if err != nil {
							log.Fatal(err)
						}

						_, err = db.Exec("UPDATE users SET password = ? WHERE uuid = ?", hash, cookie.Value)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						page.Error = "[Password] This is not your current password"
					}
				}
			} else if r.FormValue("form") == "delete" {
				password := r.FormValue("passwd")

				var db_pass string
				row, err := db.Query("SELECT password FROM users WHERE uuid = ?", cookie.Value)
				if err != nil {
					log.Fatal(err)
				}

				for row.Next() {
					err = row.Scan(&db_pass)
					if err != nil {
						log.Fatal(err)
					}
				}

				if CheckPasswordhash(password, db_pass) {
					_, err = db.Exec("DELETE FROM users WHERE uuid = ?", cookie.Value)
					if err != nil {
						log.Fatal(err)
					}
					cookie := http.Cookie{
						Name:    "user",
						Value:   "",
						Path:    "/",
						Expires: time.Unix(0, 0),
					}
					http.SetCookie(w, &cookie)
					http.Redirect(w, r, "/", http.StatusFound)
				} else {
					page.Error = "[Delete] This is not your current password"
				}
			}

		}

		db.Close()

		page.Content = user

		err = tplt.Execute(w, page)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:    "user",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
