package forum

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
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

	row, err = db.Query("SELECT COUNT(*) FROM category")
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
			row, err = db.Query("SELECT * FROM category WHERE id = ?", uuid)
			if err != nil {
				log.Fatal(err)
			}

			var cat Category
			for row.Next() {
				err = row.Scan(&cat.ID, &cat.Name, &cat.CreationDate, &cat.Pinned, &cat.LastUpdate)
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
				err = row.Scan(&post.User.ID, &post.User.Uuid, &post.User.ProfilePic, &post.User.Username, &post.User.Email, &post.User.Password, &post.User.Role, &post.User.CreationDate, &post.User.Biography, &post.User.LastSeen, &post.ID, &post.Title, &post.Content, &post.CreationDate, &uid, &post.CategoryId, &post.Pinned, &post.LastUpdate)
				if err != nil {
					log.Fatal(err)
				}

				post.CreationDate = post.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				post.LastUpdate = post.LastUpdate.(time.Time).Format("01/02/2006 15:04:05")

				post.User.CreationDate = post.User.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				post.User.LastSeen = post.User.LastSeen.(time.Time).Format("01/02/2006 15:04:05")

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

			type c struct {
				Category Category
				Posts    []Post
			}

			page.Content = c{Category: cat, Posts: posts}

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
			var uid string
			err = row.Scan(&uid)
			if err != nil {
				log.Fatal(err)
			}
			uuid_list = append(uuid_list, uid)
		}
		row.Close()

		if contains(uuid_list, uuid) {
			if r.Method == "POST" {
				content := r.FormValue("content")
				content = strings.Replace(content, "\r\n", "<br/>", -1)

				row, err = db.Query("SELECT id FROM user WHERE uuid = ?", cookie.Value)
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

				_, err = db.Exec("INSERT INTO comment (content, creation_date, user_id, post_id) VALUES (?, ?, ?, ?)", content, time.Now(), userID, uuid)
				if err != nil {
					log.Fatal(err)
				}

				_, err = db.Exec("UPDATE post SET last_update = ? WHERE id = ?", time.Now(), uuid)
				if err != nil {
					log.Fatal(err)
				}
			}

			var post Post
			row, err = db.Query("SELECT * FROM user AS u INNER JOIN post AS p ON u.id = p.user_id WHERE p.id = ?", uuid)
			if err != nil {
				log.Fatal(err)
			}
			var skip string
			for row.Next() {
				err = row.Scan(&post.User.ID, &post.User.Uuid, &post.User.ProfilePic, &post.User.Username, &post.User.Email, &post.User.Password, &post.User.Role, &post.User.CreationDate, &post.User.Biography, &post.User.LastSeen, &post.ID, &post.Title, &post.Content, &post.CreationDate, &skip, &post.CategoryId, &post.Pinned, &post.LastUpdate)
				if err != nil {
					log.Fatal(err)
				}

				post.User.CreationDate = post.User.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				post.User.LastSeen = post.User.LastSeen.(time.Time).Format("01/02/2006 15:04:05")

				post.CreationDate = post.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				post.LastUpdate = post.LastUpdate.(time.Time).Format("01/02/2006 15:04:05")
			}
			row.Close()

			row, err = db.Query("SELECT COUNT(*) FROM vote WHERE post_id = ? AND type = 'like'", post.ID)
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				err = row.Scan(&post.Likes)
				if err != nil {
					log.Fatal(err)
				}
			}
			row.Close()

			row, err = db.Query("SELECT COUNT(*) FROM vote WHERE post_id = ? AND type = 'dislikes'", post.ID)
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				err = row.Scan(&post.Dislikes)
				if err != nil {
					log.Fatal(err)
				}
			}
			row.Close()

			row, err = db.Query("SELECT * FROM comment AS c INNER JOIN user AS u ON c.user_id = u.id WHERE post_id = ?", post.ID)
			if err != nil {
				log.Fatal(err)
			}
			var comments []Comment
			var uid string
			for row.Next() {
				var comment Comment
				comment.User = User{}
				err = row.Scan(&comment.ID, &comment.Content, &comment.CreationDate, &uid, &skip, &comment.Pinned, &comment.User.ID, &comment.User.Uuid, &comment.User.ProfilePic, &comment.User.Username, &comment.User.Email, &comment.User.Password, &comment.User.Role, &comment.User.CreationDate, &comment.User.Biography, &comment.User.LastSeen)
				if err != nil {
					log.Fatal(err)
				}

				comment.CreationDate = comment.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				comment.User.CreationDate = comment.User.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				comment.User.LastSeen = comment.User.LastSeen.(time.Time).Format("01/02/2006 15:04:05")

				comments = append(comments, comment)
			}
			row.Close()

			for i := range comments {
				row, err = db.Query("SELECT COUNT(*) FROM vote WHERE comment_id = ? AND type = 'like'", comments[i].ID)
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					err = row.Scan(&comments[i].Likes)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()

				row, err = db.Query("SELECT COUNT(*) FROM vote WHERE comment_id = ? AND type = 'dislikes'", comments[i].ID)
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					err = row.Scan(&comments[i].Dislikes)
					if err != nil {
						log.Fatal(err)
					}
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

func WriteHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/write.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page Page
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true

		if r.URL.Query().Has("action") {
			action := r.URL.Query().Get("action")

			if action == "create" || action == "edit" {
				write := Write{"", action, Post{}}

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
							var post Post
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

			// Globals
			forum := AD_Forum{}
			row, err = db.Query("SELECT COUNT(*) FROM category")
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				err = row.Scan(&forum.Categories)
				if err != nil {
					log.Fatal(err)
				}
			}
			row.Close()
			row, err = db.Query("SELECT COUNT(*) FROM post")
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				err = row.Scan(&forum.Posts)
				if err != nil {
					log.Fatal(err)
				}
			}
			row.Close()
			row, err = db.Query("SELECT COUNT(*) FROM comment")
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				err = row.Scan(&forum.Comments)
				if err != nil {
					log.Fatal(err)
				}
			}
			row.Close()
			stats.Forum = forum

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
			sum_catAct := []AD_CatActivities{}
			for i := 0; i < len(CAT_ids); i++ {
				catAct := AD_CatActivities{}
				row, err = db.Query("SELECT name FROM category WHERE id = ?", CAT_ids[i])
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					err = row.Scan(&catAct.Name)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()

				for j := 1; j <= 12; j++ {
					m := strconv.Itoa(j)
					if j < 10 {
						m = "0" + m
					}
					row, err = db.Query("SELECT COUNT(*) FROM post WHERE category_id = ? AND strftime('%m', creation_date) = ?", CAT_ids[i], m)
					if err != nil {
						log.Fatal(err)
					}
					for row.Next() {
						err = row.Scan(&catAct.Activity[j-1])
						if err != nil {
							log.Fatal(err)
						}
					}
					row.Close()
				}
				sum_catAct = append(sum_catAct, catAct)
			}
			stats.CatActivities = sum_catAct

			// Users
			users := AD_Users{}
			row, err = db.Query("SELECT COUNT(*) FROM user")
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				err = row.Scan(&users.Total)
				if err != nil {
					log.Fatal(err)
				}
			}
			row.Close()
			row, err = db.Query("SELECT COUNT(*) FROM user WHERE role = 'Admin'")
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				err = row.Scan(&users.Admins)
				if err != nil {
					log.Fatal(err)
				}
			}
			row.Close()
			row, err = db.Query("SELECT COUNT(*) FROM user WHERE role = 'Mod'")
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				err = row.Scan(&users.Mods)
				if err != nil {
					log.Fatal(err)
				}
			}
			row.Close()
			row, err = db.Query("SELECT COUNT(*) FROM user WHERE role = 'Member'")
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				err = row.Scan(&users.Members)
				if err != nil {
					log.Fatal(err)
				}
			}
			row.Close()
			stats.Users = users

			// Users inscriptions
			var ui_month int
			var ui_count int
			row, err = db.Query("SELECT COUNT(*) AS count, strftime('%m', creation_date) as month FROM user WHERE strftime('%Y', creation_date) = strftime('%Y', date()) GROUP BY month")
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

			/* USERS TAB */
			row, err = db.Query("SELECT * FROM user")
			if err != nil {
				log.Fatal(err)
			}
			for row.Next() {
				user := User{}
				err = row.Scan(&user.ID, &user.Uuid, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreationDate, &user.Biography, &user.LastSeen)
				if err != nil {
					log.Fatal(err)
				}
				user.CreationDate = user.CreationDate.(time.Time).Format("2006/01/02 15:04:05")
				user.LastSeen = user.LastSeen.(time.Time).Format("2006/01/02 15:04:05")

				admin.Users = append(admin.Users, user)
			}
			row.Close()

			/* REPORTS TAB */
			reports := []Report{}

			row, err = db.Query("SELECT * FROM report")
			if err != nil {
				log.Fatal(err)
			}

			uids := []interface{}{}
			pids := []interface{}{}
			cids := []interface{}{}
			for row.Next() {
				report := Report{}
				var userID interface{}
				var postID interface{}
				var commentID interface{}
				err = row.Scan(&report.ID, &report.Type, &report.Reason, &report.CreationDate, &userID, &postID, &commentID)
				if err != nil {
					log.Fatal(err)
				}
				report.CreationDate = report.CreationDate.(time.Time).Format("2006/01/02 15:04:05")

				uids = append(uids, userID)
				pids = append(pids, postID)
				cids = append(cids, commentID)
				reports = append(reports, report)
			}
			row.Close()

			for i := 0; i < len(reports); i++ {
				var skip interface{}
				if reports[i].Type == "user" {
					row, err = db.Query("SELECT * FROM user WHERE id = ?", uids[i])
					if err != nil {
						log.Fatal(err)
					}
					for row.Next() {
						err = row.Scan(&reports[i].User.ID, &reports[i].User.Uuid, &reports[i].User.Username, &reports[i].User.Email, &reports[i].User.Password, &reports[i].User.Role, &reports[i].User.CreationDate, &reports[i].User.Biography, &reports[i].User.LastSeen)
						if err != nil {
							log.Fatal(err)
						}
					}
					row.Close()
					admin.Reports.Users = append(admin.Reports.Users, reports[i])
				} else if reports[i].Type == "post" {
					row, err = db.Query("SELECT * FROM post AS p INNER JOIN user AS u ON p.user_id = u.id WHERE p.id = ?", pids[i])
					if err != nil {
						log.Fatal(err)
					}
					for row.Next() {
						err = row.Scan(&reports[i].Post.ID, &reports[i].Post.Title, &reports[i].Post.Content, &reports[i].Post.CreationDate, &skip, &reports[i].Post.CategoryId, &reports[i].Post.Pinned, &reports[i].Post.LastUpdate, &reports[i].User.ID, &reports[i].User.Uuid, &reports[i].User.Username, &reports[i].User.Email, &reports[i].User.Password, &reports[i].User.Role, &reports[i].User.CreationDate, &reports[i].User.Biography, &reports[i].User.LastSeen)
						if err != nil {
							log.Fatal(err)
						}
					}
					row.Close()
					admin.Reports.Posts = append(admin.Reports.Posts, reports[i])
				} else if reports[i].Type == "comment" {
					row, err = db.Query("SELECT * FROM comment AS c INNER JOIN post AS p ON c.post_id = p.id INNER JOIN user AS u ON c.user_id = u.id WHERE c.id = ?", cids[i])
					if err != nil {
						log.Fatal(err)
					}
					for row.Next() {
						err = row.Scan(&reports[i].Comment.ID, &reports[i].Comment.Content, &reports[i].Comment.CreationDate, &skip, &skip, &reports[i].Comment.Pinned, &reports[i].Post.ID, &reports[i].Post.Title, &reports[i].Post.Content, &reports[i].Post.CreationDate, &skip, &reports[i].Post.CategoryId, &reports[i].Post.Pinned, &reports[i].Post.LastUpdate, &reports[i].User.ID, &reports[i].User.Uuid, &reports[i].User.Username, &reports[i].User.Email, &reports[i].User.Password, &reports[i].User.Role, &reports[i].User.CreationDate, &reports[i].User.Biography, &reports[i].User.LastSeen)
						if err != nil {
							log.Fatal(err)
						}
					}
					row.Close()

					admin.Reports.Comments = append(admin.Reports.Comments, reports[i])
				}
			}

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

func AdminDeleteHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("user")
	if cookie != nil {
		userID := cookie.Value

		if r.URL.Query().Has("id") && r.URL.Query().Has("type") {
			reqID := r.URL.Query().Get("id")
			reqType := r.URL.Query().Get("type")

			db, err := sql.Open("sqlite3", "./forum.db")
			if err != nil {
				log.Fatal(err)
			}

			row, err := db.Query("SELECT role FROM user WHERE uuid = ?", userID)
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

				if reqType == "user" {
					row, err = db.Query("SELECT COUNT(*) FROM user WHERE id = ?", reqID)
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

					if count > 0 {
						_, err = db.Exec("DELETE FROM user WHERE id = ?", reqID)
						if err != nil {
							log.Fatal(err)
						}

						http.Redirect(w, r, "/", http.StatusOK)
					} else {
						http.Redirect(w, r, "/", http.StatusFound)
					}
				} else if reqType == "report" {
					row, err = db.Query("SELECT COUNT(*) FROM report WHERE id = ?", reqID)
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

					if count > 0 {
						_, err = db.Exec("DELETE FROM report WHERE id = ?", reqID)
						if err != nil {
							log.Fatal(err)
						}
						http.Redirect(w, r, "/", http.StatusOK)
					} else {
						http.Redirect(w, r, "/", http.StatusFound)
					}
				}

			} else {
				http.Redirect(w, r, "/", http.StatusFound)
			}

			db.Close()
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
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
			err = row.Scan(&user.ID, &user.Uuid, &user.ProfilePic, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreationDate, &user.Biography, &user.LastSeen)
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
			post.User = User{}
			err = row.Scan(&post.ID, &post.Title, &post.Content, &post.CreationDate, &skip, &post.CategoryId, &post.Pinned, &post.LastUpdate, &post.Category.ID, &post.Category.Name, &post.Category.CreationDate, &post.Category.Pinned, &post.Category.LastUpdate)
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
			err = row.Scan(&comment.ID, &comment.Content, &comment.CreationDate, &skip, &skip, &comment.Pinned, &comment.PostID.ID, &comment.PostID.Title, &comment.PostID.Content, &comment.PostID.CreationDate, &skip, &comment.PostID.CategoryId, &comment.PostID.Pinned, &comment.PostID.LastUpdate)
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
			} else if r.URL.Query().Get("type") == "comment" {
				_, err = db.Exec("DELETE FROM comment WHERE id = ?", r.URL.Query().Get("id"))
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		db.Close()
		http.Redirect(w, r, "/user/profile", http.StatusFound)
	} else {
		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
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
