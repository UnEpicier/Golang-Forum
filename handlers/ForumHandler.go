package forum

import (
	"database/sql"
	f "forum"
	"html/template"
	"log"
	"net/http"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

func ForumHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/forum.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page f.Page
	page.Logged = false

	var forum f.Forum

	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true

		row, err := db.Query("SELECT role FROM user WHERE uuid = ?", cookie.Value)
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			err = row.Scan(&forum.Role)
			if err != nil {
				log.Fatal(err)
			}
		}
		row.Close()
	}

	var categories []f.Category

	if r.Method == "POST" {
		err = r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		title := r.FormValue("category_name")

		if forum.Role == "Admin" {
			_, err := db.Exec("INSERT INTO category (`name`, `creation_date`, `pinned`, `last_update`) VALUES (?, ?, 0, ?)", title, time.Now(), time.Now())
			if err != nil {
				log.Fatal(err)
			}
		}
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
			var category f.Category
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
		forum.Categories = []f.Category{}
		forum.Error = "No categories found"
	}

	db.Close()

	page.Content = forum

	err = tplt.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}
