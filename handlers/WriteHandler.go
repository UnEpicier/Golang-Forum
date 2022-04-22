package forum

import (
	"database/sql"
	f "forum"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func WriteHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/write.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page f.Page
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true

		if r.URL.Query().Has("action") {
			action := r.URL.Query().Get("action")

			if action == "create" || action == "edit" {
				write := f.Write{"", action, f.Post{}}

				db, err := sql.Open("sqlite3", "./forum.db")
				if err != nil {
					log.Fatal(err)
				}

				row, err := db.Query("SELECT name FROM category")
				if err != nil {
					log.Fatal(err)
				}
				var categories []string
				for row.Next() {
					var cat string
					err = row.Scan(&cat)
					if err != nil {
						log.Fatal(err)
					}
					categories = append(categories, cat)
				}

				write.Categories = strings.Join(categories, "/")

				if action == "edit" {
					if r.URL.Query().Has("id") {
						postID := r.URL.Query().Get("id")

						row, err := db.Query("SELECT id FROM user WHERE uuid = ?", cookie.Value)
						if err != nil {
							log.Fatal(err)
						}
						var userID int
						for row.Next() {
							err = row.Scan(&userID)
							if err != nil {
								log.Fatal(err)
							}
						}
						row.Close()

						row, err = db.Query("SELECT COUNT(*) FROM post WHERE id = ? AND user_id = ?", postID, userID)
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
							if r.Method == "POST" {
								err = r.ParseForm()
								if err != nil {
									log.Fatal(err)
								}
								category := r.FormValue("category")
								title := r.FormValue("title")
								content := r.FormValue("content")
								content = strings.Replace(content, "\r\n", "<br/>", -1)

								row, err = db.Query("SELECT id FROM category WHERE name = ?", category)
								if err != nil {
									log.Fatal(err)
								}
								var catID int
								for row.Next() {
									err = row.Scan(&catID)
									if err != nil {
										log.Fatal(err)
									}
								}
								row.Close()

								_, err := db.Exec("UPDATE post SET title = ?, content = ?, category_id = ?, last_update = ? WHERE id = ?", title, content, catID, time.Now(), postID)
								if err != nil {
									log.Fatal(err)
								}

								_, err = db.Exec("UPDATE category SET last_update = ? WHERE category_id = ?", time.Now(), catID)
								if err != nil {
									log.Fatal(err)
								}

								http.Redirect(w, r, "/post?id="+postID, http.StatusFound)
							}

							row, err = db.Query("SELECT * FROM post AS p INNER JOIN category AS c ON p.category_id = c.id WHERE p.id = ? AND p.user_id = ?", postID, userID)
							if err != nil {
								log.Fatal(err)
							}
							var post f.Post
							for row.Next() {
								var skip interface{}
								err = row.Scan(&post.ID, &post.Title, &post.Content, &post.CreationDate, &skip, &post.CategoryId, &post.Pinned, &post.LastUpdate, &post.Category.ID, &post.Category.Name, &post.Category.CreationDate, &post.Category.Pinned, &post.Category.LastUpdate)
								if err != nil {
									log.Fatal(err)
								}
								post.CreationDate = post.CreationDate.(time.Time).Format("2006/01/02 15:04:05")
								post.LastUpdate = post.LastUpdate.(time.Time).Format("2006/01/02 15:04:05")

								post.Category.CreationDate = post.Category.CreationDate.(time.Time).Format("2006/01/02 15:04:05")
								post.Category.LastUpdate = post.Category.LastUpdate.(time.Time).Format("2006/01/02 15:04:05")
							}
							row.Close()
							write.Post = post
						} else {
							http.Redirect(w, r, "/forum", http.StatusFound)
						}
					} else {
						http.Redirect(w, r, "/forum", http.StatusFound)
					}
				} else if action == "create" {
					if r.URL.Query().Has("category-id") {
						row, err := db.Query("SELECT COUNT(*) FROM category WHERE id = ?", r.URL.Query().Get("category-id"))
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
							catID := r.URL.Query().Get("category-id")

							row, err = db.Query("SELECT name FROM category WHERE id = ?", catID)
							if err != nil {
								log.Fatal(err)
							}
							var catName string
							for row.Next() {
								err = row.Scan(&catName)
								if err != nil {
									log.Fatal(err)
								}
							}
							row.Close()
							write.Post.Category.Name = catName
						}
					}

					if r.Method == "POST" {
						err := r.ParseForm()
						if err != nil {
							log.Fatal(err)
						}
						category := r.FormValue("category")
						title := r.FormValue("title")
						content := r.FormValue("content")
						content = strings.Replace(content, "\r\n", "<br/>", -1)

						row, err := db.Query("SELECT id FROM user WHERE uuid = ?", cookie.Value)
						if err != nil {
							log.Fatal(err)
						}
						var userID int
						for row.Next() {
							err = row.Scan(&userID)
							if err != nil {
								log.Fatal(err)
							}
						}
						row.Close()

						row, err = db.Query("SELECT id FROM category WHERE name = ?", category)
						if err != nil {
							log.Fatal(err)
						}
						var catID int
						for row.Next() {
							err = row.Scan(&catID)
							if err != nil {
								log.Fatal(err)
							}
						}
						row.Close()

						acTime := time.Now()
						_, err = db.Exec("INSERT INTO post (title, content, creation_date, user_id, category_id, pinned, last_update) VALUES (?, ?, ?, ?, ?, 0, ?)", title, content, acTime, userID, catID, acTime)
						if err != nil {
							log.Fatal(err)
						}

						_, err = db.Exec("UPDATE category SET last_update = ? WHERE id = ?", acTime, catID)
						if err != nil {
							log.Fatal(err)
						}

						row, err = db.Query("SELECT id FROM post WHERE user_id = ? AND title = ? AND content = ? AND creation_date = ? AND category_id = ?", userID, title, content, acTime, catID)
						if err != nil {
							log.Fatal(err)
						}
						var pid int
						for row.Next() {
							err = row.Scan(&pid)
							if err != nil {
								log.Fatal(err)
							}
						}
						http.Redirect(w, r, "/post?id="+strconv.Itoa(pid), http.StatusFound)
					}
				}

				db.Close()
				page.Content = write

				err = tplt.Execute(w, page)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				http.Redirect(w, r, "/forum", http.StatusFound)
			}
		} else {
			http.Redirect(w, r, "/forum", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	}
}
