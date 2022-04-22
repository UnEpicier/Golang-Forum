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

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/category.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page f.Page
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true
	}

	filter := ""
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		filter = r.FormValue("filter")
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

		if f.Contains(uuid_list, uuid) {
			row, err = db.Query("SELECT * FROM category WHERE id = ?", uuid)
			if err != nil {
				log.Fatal(err)
			}

			var cat f.Category
			for row.Next() {
				err = row.Scan(&cat.ID, &cat.Name, &cat.CreationDate, &cat.Pinned, &cat.LastUpdate)
				if err != nil {
					log.Fatal(err)
				}
			}

			posts := []f.Post{}

			row, err = db.Query("SELECT * FROM user as u INNER JOIN post as p ON u.id = p.user_id WHERE p.category_id = ? ORDER By p.last_update DESC", uuid)
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				var uid string
				var post f.Post
				err = row.Scan(&post.User.ID, &post.User.Uuid, &post.User.ProfilePic, &post.User.Username, &post.User.Email, &post.User.Password, &post.User.Role, &post.User.CreationDate, &post.User.Biography, &post.User.LastSeen, &post.ID, &post.Title, &post.Content, &post.CreationDate, &uid, &post.CategoryId, &post.Pinned, &post.LastUpdate)
				if err != nil {
					log.Fatal(err)
				}

				posts = append(posts, post)
			}

			for i := 0; i < len(posts); i++ {
				row, err = db.Query("SELECT COUNT(*) FROM vote WHERE post_id = ? AND type = 'like'", posts[i].ID)
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					err = row.Scan(&posts[i].Likes)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()

				row, err = db.Query("SELECT COUNT(*) FROM vote WHERE post_id = ? AND type = 'dislikes'", posts[i].ID)
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					err = row.Scan(&posts[i].Dislikes)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()

				row, err = db.Query("SELECT COUNT(*) FROM comment WHERE post_id = ?", posts[i].ID)
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					err = row.Scan(&posts[i].CommentNB)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()
			}

			// Sort posts by last update
			if filter == "Latest" {
				for i := 0; i < len(posts); i++ {
					for j := i + 1; j < len(posts); j++ {
						if posts[i].LastUpdate.(time.Time).Before(posts[j].LastUpdate.(time.Time)) {
							temp := posts[i]
							posts[i] = posts[j]
							posts[j] = temp
						}
					}
				}
			} else if filter == "Oldest" {
				for i := 0; i < len(posts); i++ {
					for j := i + 1; j < len(posts); j++ {
						if posts[i].LastUpdate.(time.Time).After(posts[j].LastUpdate.(time.Time)) {
							temp := posts[i]
							posts[i] = posts[j]
							posts[j] = temp
						}
					}
				}
			} else if filter == "Most Liked" {
				for i := 0; i < len(posts); i++ {
					for j := i + 1; j < len(posts); j++ {
						if posts[i].Likes < posts[j].Likes {
							temp := posts[i]
							posts[i] = posts[j]
							posts[j] = temp
						}
					}
				}
			} else if filter == "Most Disliked" {
				for i := 0; i < len(posts); i++ {
					for j := i + 1; j < len(posts); j++ {
						if posts[i].Likes > posts[j].Likes {
							temp := posts[i]
							posts[i] = posts[j]
							posts[j] = temp
						}
					}
				}
			}

			for i := 0; i < len(posts); i++ {
				posts[i].CreationDate = posts[i].CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				posts[i].LastUpdate = posts[i].LastUpdate.(time.Time).Format("01/02/2006 15:04:05")

				posts[i].User.CreationDate = posts[i].User.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				posts[i].User.LastSeen = posts[i].User.LastSeen.(time.Time).Format("01/02/2006 15:04:05")
			}

			type c struct {
				Category f.Category
				Posts    []f.Post
				Filter   string
			}

			page.Content = c{Category: cat, Posts: posts, Filter: filter}

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
