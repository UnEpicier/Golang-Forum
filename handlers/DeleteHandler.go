package forum

import (
	"database/sql"
	f "forum"
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
				var post []f.Post
				var comment []f.Comment
				rows, err := db.Query("SELECT * FROM post WHERE category_id = ?", r.URL.Query().Get("id"))
				if err != nil {
					log.Fatal(err)
				}
				for rows.Next() {
					var p f.Post
					err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreationDate, &p.User.ID, &p.Category.ID, &p.Pinned, &p.LastUpdate)
					if err != nil {
						log.Fatal(err)
					}
					post = append(post, p)
				}
				rows, err = db.Query("SELECT * FROM comment WHERE post_id = ?", r.URL.Query().Get("id"))
				if err != nil {
					log.Fatal(err)
				}
				for rows.Next() {
					var c f.Comment
					err = rows.Scan(&c.ID, &c.Content, &c.CreationDate, &c.User.ID, &c.PostID.ID, &c.Pinned)
					if err != nil {
						log.Fatal(err)
					}
					comment = append(comment, c)
				}
				rows.Close()

				for _, c := range comment {
					_, err := db.Exec("DELETE FROM Vote WHERE comment_id = ?", c.ID)
					if err != nil {
						log.Fatal(err)
					}
				}
				_, err = db.Exec("DELETE FROM comment WHERE post_id = ?", r.URL.Query().Get("id"))
				if err != nil {
					log.Fatal(err)
				}
				for _, p := range post {
					_, err := db.Exec("DELETE FROM Vote WHERE post_id = ?", p.ID)
					if err != nil {
						log.Fatal(err)
					}
				}
				_, err = db.Exec("DELETE FROM post WHERE category_id = ?", r.URL.Query().Get("id"))
				if err != nil {
					log.Fatal(err)
				}
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
