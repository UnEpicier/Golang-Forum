package forum

import (
	"database/sql"
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

	row, err := db.Query("SELECT COUNT(*) FROM category")
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
		row, err = db.Query("SELECT * FROM category ORDER BY last_update DESC")
		if err != nil {
			log.Fatal(err)
		}

		for row.Next() {
			var category Category
			err = row.Scan(&category.ID, &category.Name, &category.CreationDate, &category.Pinned, &category.LastUpdate)
			if err != nil {
				log.Fatal(err)
			}
			categories = append(categories, category)
		}
		row.Close()

		for i := 0; i < len(categories); i++ {
			categories[i].CreationDate = categories[i].CreationDate.(time.Time).Format("01/02/2006 15:04:05")
			categories[i].LastUpdate = categories[i].LastUpdate.(time.Time).Format("01/02/2006 15:04:05")
		}

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
		row, err := db.Query("SELECT id FROM post")
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
			row, err = db.Query("SELECT * FROM post WHERE id = ?", uuid)
			if err != nil {
				log.Fatal(err)
			}
			var uid string
			for row.Next() {
				err = row.Scan(&post.ID, &post.Title, &post.Content, &post.CreationDate, &post.UserID, &post.UpVotes, &post.DownVotes, &post.CategoryId, &post.Pinned, &post.LastUpdate)
				if err != nil {
					log.Fatal(err)
				}
				post.CreationDate = post.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				post.LastUpdate = post.LastUpdate.(time.Time).Format("01/02/2006 15:04:05")
			}
			row.Close()

			row, err = db.Query("SELECT * FROM user WHERE id = ?", uid)
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				var user User
				err = row.Scan(&user.ID, &user.Uuid, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreationDate, &user.Biography, &user.LastSeen)
				if err != nil {
					log.Fatal(err)
				}
				user.CreationDate = user.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				user.LastSeen = user.LastSeen.(time.Time).Format("01/02/2006 15:04:05")

				post.UserID = user
			}
			row.Close()

			row, err = db.Query("SELECT * FROM comment WHERE post_id = ?", post.ID)
			if err != nil {
				log.Fatal(err)
			}
			var comments []Comment
			for row.Next() {
				var comment Comment
				comment.UserID = User{}
				err = row.Scan(&comment.ID, &comment.Content, &comment.CreationDate, &uid, &comment.PostID, &comment.UpVotes, &comment.DownVotes, &comment.Pinned)
				if err != nil {
					log.Fatal(err)
				}
				comments = append(comments, comment)
			}
			row.Close()

			for comment := range comments {
				row, err = db.Query("SELECT * FROM user WHERE id = ?", uid)
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					var user User
					err = row.Scan(&user.ID, &user.Uuid, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreationDate, &user.Biography, &user.LastSeen)
					if err != nil {
						log.Fatal(err)
					}

					user.CreationDate = user.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
					user.LastSeen = user.LastSeen.(time.Time).Format("01/02/2006 15:04:05")

					comments[comment].UserID = user
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
		row, err := db.Query("SELECT id FROM category")
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			var uid string
			err = row.Scan(&uid)
			if err != nil {
				log.Fatal(err)
			}
			uuid_list = append(uuid_list, uid)
		}

		if contains(uuid_list, uuid) {
			row, err = db.Query("SELECT name FROM category WHERE id = ?", uuid)
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

			row, err = db.Query("SELECT * FROM user as u INNER JOIN post as p ON u.id = p.user_id WHERE p.category_id = ? ORDER BY p.last_update DESC", uuid)
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				var uid string
				var post Post
				err = row.Scan(&post.UserID.ID, &post.UserID.Uuid, &post.UserID.Username, &post.UserID.Email, &post.UserID.Password, &post.UserID.Role, &post.UserID.CreationDate, &post.UserID.Biography, &post.UserID.LastSeen, &post.ID, &post.Title, &post.Content, &post.CreationDate, &uid, &post.UpVotes, &post.DownVotes, &post.CategoryId, &post.Pinned, &post.LastUpdate)
				if err != nil {
					log.Fatal(err)
				}

				post.CreationDate = post.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				post.LastUpdate = post.LastUpdate.(time.Time).Format("01/02/2006 15:04:05")

				post.UserID.CreationDate = post.UserID.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				post.UserID.LastSeen = post.UserID.LastSeen.(time.Time).Format("01/02/2006 15:04:05")

				posts = append(posts, post)
			}

			type c struct {
				Name  string
				Posts []Post
			}

			page.Content = c{Name: name, Posts: posts}

			err := tplt.Execute(w, page)
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
		row, err := db.Query("SELECT name FROM category")
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
				row, err := db.Query("SELECT id, name FROM category")
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
					_, err = db.Exec("INSERT INTO category (name, creation_date, pinned, last_update) VALUES (?, ?, 0, ?)", capitalize(category), time.Now(), time.Now())
					if err != nil {
						log.Fatal(err)
					}
				}

				var id string
				row, err = db.Query("SELECT id FROM user WHERE uuid = ?", cookie.Value)
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					err = row.Scan(&id)
					if err != nil {
						log.Fatal(err)
					}
				}

				_, err = db.Exec("INSERT INTO post (title, content, creation_date, user_id, up_vote, down_vote, category_id, pinned, last_update) VALUES (?, ?, ?, ?, '{}', '{}', ?, 0, ?)", title, content, time.Now(), id, uid, time.Now())
				if err != nil {
					log.Fatal(err, " ligne 412")
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

		row, err := db.Query("SELECT role FROM user WHERE uuid = ?", cookie.Value)
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
			admin := Admin{}

			/*
				STATS
			*/
			stats := Stats{}

			// Categories
			CAT_ids := []int{}

			row, err = db.Query("SELECT id, name FROM category")
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				var id int
				c := AD_Categories{}
				err = row.Scan(&id, &c.Name)
				if err != nil {
					log.Fatal(err)
				}
				CAT_ids = append(CAT_ids, id)
				stats.Categories = append(stats.Categories, c)
			}
			row.Close()

			for i := 0; i < len(CAT_ids); i++ {
				row, err = db.Query("SELECT COUNT(*) FROM post WHERE category_id = ?", CAT_ids[i])
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					err = row.Scan(&stats.Categories[i].Count)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()
			}

			// Categries activities per month
			/* for i := 0; i < len(CAT_ids); i++ {
				row, err = db.Query("SELECT COUNT(*) FROM post WHERE category_id = ? AND creation_date > ?", CAT_ids[i], time.Now().AddDate(0, -1, 0))
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					err = row.Scan(&stats.Categories[i].Activity)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()
			} */

			// Users inscriptions
			var ui_month int
			var ui_count int
			row, err = db.Query("SELECT COUNT(*) AS count, strftime('%m', creation_date) as month FROM user GROUP BY month")
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				err = row.Scan(&ui_count, &ui_month)
				if err != nil {
					log.Fatal(err)
				}
				inscr := AD_Inscription{}
				inscr.Month = time.Month(ui_month).String()
				inscr.Count = ui_count
				stats.Inscriptions = append(stats.Inscriptions, inscr)
			}
			row.Close()

			admin.Stats = stats

			page.Content = admin
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

		if db_email == email && CheckPasswordhash(password, db_password) {
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
			pass, err := HashPassword(password)
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

		row, err := db.Query("SELECT COUNT(*) FROM user WHERE uuid = ?", cookie.Value)
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

		if count <= 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		row, err = db.Query("SELECT * FROM user WHERE uuid = ? LIMIT 1", cookie.Value)
		if err != nil {
			log.Fatal(err)
		}

		for row.Next() {
			err = row.Scan(&user.ID, &user.Uuid, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreationDate, &user.Biography, &user.LastSeen)
			if err != nil {
				log.Fatal(err)
			}
		}
		row.Close()

		user.CreationDate = user.CreationDate.(time.Time).Format("02/01/2006 15:04:05")
		user.LastSeen = user.LastSeen.(time.Time).Format("02/01/2006 15:04:05")

		// Posts
		row, err = db.Query("SELECT * FROM post AS p INNER JOIN category AS c ON p.category_id = c.id WHERE p.user_id = ? ORDER BY p.creation_date DESC", user.ID)
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			var skip int
			var post Post
			post.UserID = User{}
			err = row.Scan(&post.ID, &post.Title, &post.Content, &post.CreationDate, &skip, &post.UpVotes, &post.DownVotes, &post.CategoryId, &post.Pinned, &post.LastUpdate, &post.Category.ID, &post.Category.Name, &post.Category.CreationDate, &post.Category.Pinned, &post.Category.LastUpdate)
			if err != nil {
				log.Fatal(err)
			}

			post.CreationDate = post.CreationDate.(time.Time).Format("02/01/2006 15:04:05")
			post.LastUpdate = post.LastUpdate.(time.Time).Format("02/01/2006 15:04:05")

			user.Posts = append(user.Posts, post)
		}
		row.Close()

		// Comments
		row, err = db.Query("SELECT * FROM comment AS c INNER JOIN post AS p ON c.post_id = p.id ORDER BY c.creation_date DESC")
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			var skip int
			var comment Comment
			err = row.Scan(&comment.ID, &comment.Content, &comment.CreationDate, &skip, &skip, &comment.UpVotes, &comment.DownVotes, &comment.Pinned, &comment.PostID.ID, &comment.PostID.Title, &comment.PostID.Content, &comment.PostID.CreationDate, &skip, &comment.PostID.UpVotes, &comment.PostID.DownVotes, &comment.PostID.CategoryId, &comment.PostID.Pinned, &comment.PostID.LastUpdate)
			if err != nil {
				log.Fatal(err)
			}

			comment.CreationDate = comment.CreationDate.(time.Time).Format("02/01/2006 15:04:05")

			user.Comments = append(user.Comments, comment)
		}
		row.Close()

		/*
			Settings
		*/
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				log.Fatal(err)
			}

			if r.FormValue("form") == "biography" {
				biography := r.FormValue("bio")

				_, err := db.Exec("UPDATE user SET biography = ? WHERE uuid = ?", biography, cookie.Value)
				if err != nil {
					log.Fatal(err)
				}
				http.Redirect(w, r, "/user/profile", http.StatusFound)

			} else if r.FormValue("form") == "username" {
				current := r.FormValue("username")
				new := r.FormValue("newusername")
				password := r.FormValue("passwd")

				var count int
				row, err := db.Query("SELECT COUNT(*) FROM user WHERE uuid = ? AND username = ?", cookie.Value, current)
				if err != nil {
					log.Fatal(err)
				}

				for row.Next() {
					err = row.Scan(&count)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()

				if count > 0 {
					var db_pass string

					row, err = db.Query("SELECT password FROM user WHERE uuid = ? AND username = ?", cookie.Value, current)
					if err != nil {
						log.Fatal(err)
					}
					for row.Next() {
						err = row.Scan(&db_pass)
						if err != nil {
							log.Fatal(err)
						}
					}
					row.Close()

					if CheckPasswordhash(password, db_pass) {
						_, err = db.Exec("UPDATE user SET username = ? WHERE uuid = ?", new, cookie.Value)
						if err != nil {
							log.Fatal(err)
						}
						http.Redirect(w, r, "/user/profile", http.StatusFound)
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
					row, err := db.Query("SELECT COUNT(*) FROM user WHERE uuid = ? AND email = ?", cookie.Value, current)
					if err != nil {
						log.Fatal(err)
					}

					for row.Next() {
						err = row.Scan(&count)
						if err != nil {
							log.Fatal(err)
						}
					}
					row.Close()

					if count > 0 {
						var db_pass string

						row, err = db.Query("SELECT password FROM user WHERE uuid = ? AND email = ?", cookie.Value, current)
						if err != nil {
							log.Fatal(err)
						}
						for row.Next() {
							err = row.Scan(&db_pass)
							if err != nil {
								log.Fatal(err)
							}
						}
						row.Close()

						if CheckPasswordhash(password, db_pass) {
							_, err = db.Exec("UPDATE user SET email = ? WHERE uuid = ?", new, cookie.Value)
							if err != nil {
								log.Fatal(err)
							}
							http.Redirect(w, r, "/user/profile", http.StatusFound)
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
					row, err := db.Query("SELECT password FROM user WHERE uuid = ?", cookie.Value)
					if err != nil {
						log.Fatal(err)
					}

					for row.Next() {
						err = row.Scan(&db_pass)
						if err != nil {
							log.Fatal(err)
						}
					}
					row.Close()

					if CheckPasswordhash(current, db_pass) {
						hash, err := HashPassword(new)
						if err != nil {
							log.Fatal(err)
						}

						_, err = db.Exec("UPDATE user SET password = ? WHERE uuid = ?", hash, cookie.Value)
						if err != nil {
							log.Fatal(err)
						}
						http.Redirect(w, r, "/user/profile", http.StatusFound)
					} else {
						page.Error = "[Password] This is not your current password"
					}
				}
			} else if r.FormValue("form") == "delete" {
				password := r.FormValue("passwd")

				var db_pass string
				row, err := db.Query("SELECT password FROM user WHERE uuid = ?", cookie.Value)
				if err != nil {
					log.Fatal(err)
				}

				for row.Next() {
					err = row.Scan(&db_pass)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()

				if CheckPasswordhash(password, db_pass) {
					_, err = db.Exec("DELETE FROM user WHERE uuid = ?", cookie.Value)
					if err != nil {
						log.Fatal(err)
					}

					user := http.Cookie{
						Name:    "user",
						Value:   "",
						Path:    "/",
						Expires: time.Unix(0, 0),
					}
					http.SetCookie(w, &user)
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
