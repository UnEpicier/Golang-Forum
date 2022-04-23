package forum

import (
	"database/sql"
	f "forum"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/post.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page f.Page
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

		if f.Contains(uuid_list, uuid) {
			var userID int
			if cookie != nil {
				row, err = db.Query("SELECT id FROM user WHERE uuid = ?", cookie.Value)
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					err = row.Scan(&userID)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()
			}

			if r.Method == "POST" {
				err = r.ParseForm()
				if err != nil {
					log.Fatal(err)
				}

				if r.FormValue("form") == "createcomment" {
					content := r.FormValue("content")
					content = strings.Replace(content, "\r\n", "<br/>", -1)

					_, err = db.Exec("INSERT INTO comment (content, creation_date, user_id, post_id) VALUES (?, ?, ?, ?)", content, time.Now(), userID, uuid)
					if err != nil {
						log.Fatal(err)
					}

					_, err = db.Exec("UPDATE post SET last_update = ? WHERE id = ?", time.Now(), uuid)
					if err != nil {
						log.Fatal(err)
					}
					http.Redirect(w, r, "/post?id="+uuid, http.StatusFound)
				} else if r.FormValue("form") == "like" {
					type_ := r.FormValue("type")
					id := r.FormValue("id")
					like := r.FormValue("like")
					dislike := r.FormValue("dislike")

					row, err := db.Query("SELECT id FROM user WHERE uuid = ?", cookie.Value)
					if err != nil {
						log.Fatal(err)
					}
					var userID string
					for row.Next() {
						err = row.Scan(&userID)
						if err != nil {
							log.Fatal(err)
						}
					}

					if page.Logged {
						if type_ == "comment" {
							if like == "" && dislike == "" {
								_, err = db.Exec("DELETE FROM vote WHERE comment_id = ? AND user_id = ?", id, userID)
								if err != nil {
									log.Fatal(err)
								}
								http.Redirect(w, r, "/post?id="+uuid, http.StatusFound)
							} else {
								if like == "on" && dislike == "" {
									_, err = db.Exec("INSERT INTO vote (type, user_id, comment_id) VALUES ('like', ?, ?)", userID, id)
									if err != nil {
										log.Fatal(err)
									}
									_, err = db.Exec("DELETE FROM vote WHERE comment_id = ? AND user_id = ? AND type = 'dislike'", id, userID)
									if err != nil {
										log.Fatal(err)
									}
									http.Redirect(w, r, "/post?id="+uuid, http.StatusFound)
								} else if like == "" && dislike == "on" {
									_, err = db.Exec("INSERT INTO vote (type, user_id, comment_id) VALUES ('dislike', ?, ?)", userID, id)
									if err != nil {
										log.Fatal(err)
									}
									_, err = db.Exec("DELETE FROM vote WHERE comment_id = ? AND user_id = ? AND type = 'like'", id, userID)
									if err != nil {
										log.Fatal(err)
									}
									http.Redirect(w, r, "/post?id="+uuid, http.StatusFound)
								}
							}
						} else if type_ == "post" {
							if like == "" && dislike == "" {
								_, err = db.Exec("DELETE FROM vote WHERE post_id = ? AND user_id = ?", id, userID)
								if err != nil {
									log.Fatal(err)
								}
								http.Redirect(w, r, "/post?id="+uuid, http.StatusFound)
							} else {
								if like == "on" && dislike == "" {
									_, err = db.Exec("INSERT INTO vote (type, user_id, post_id) VALUES ('like', ?, ?)", userID, id)
									if err != nil {
										log.Fatal(err)
									}
									_, err = db.Exec("DELETE FROM vote WHERE post_id = ? AND user_id = ? AND type = 'dislike'", id, userID)
									if err != nil {
										log.Fatal(err)
									}
									http.Redirect(w, r, "/post?id="+uuid, http.StatusFound)
								} else if like == "" && dislike == "on" {
									_, err = db.Exec("INSERT INTO vote (type, user_id, post_id) VALUES ('dislike', ?, ?)", userID, id)
									if err != nil {
										log.Fatal(err)
									}
									_, err = db.Exec("DELETE FROM vote WHERE post_id = ? AND user_id = ? AND type = 'like'", id, userID)
									if err != nil {
										log.Fatal(err)
									}
									http.Redirect(w, r, "/post?id="+uuid, http.StatusFound)
								}
							}
						}

					}
				}
			}

			var content f.PostPage

			var post f.Post

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

			row, err = db.Query("SELECT COUNT(*) FROM vote WHERE post_id = ? AND type = 'dislike'", post.ID)
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

			if cookie != nil {
				row, err = db.Query("SELECT COUNT(*) FROM vote WHERE post_id = ? AND user_id = ? AND type = 'like'", post.ID, userID)
				if err != nil {
					log.Fatal(err)
				}
				var userLikes int
				for row.Next() {
					err = row.Scan(&userLikes)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()
				if userLikes > 0 {
					post.Vote = "like"
				}

				row, err = db.Query("SELECT COUNT(*) FROM vote WHERE post_id = ? AND user_id = ? AND type = 'dislike'", post.ID, userID)
				if err != nil {
					log.Fatal(err)
				}
				var userDislikes int
				for row.Next() {
					err = row.Scan(&userDislikes)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()
				if userDislikes > 0 {
					post.Vote = "dislike"
				}
			}

			row, err = db.Query("SELECT * FROM comment AS c INNER JOIN user AS u ON c.user_id = u.id WHERE post_id = ? ORDER BY c.creation_date ASC", post.ID)
			if err != nil {
				log.Fatal(err)
			}
			var c []f.Comment
			var uid string
			for row.Next() {
				var comment f.Comment
				comment.User = f.User{}
				err = row.Scan(&comment.ID, &comment.Content, &comment.CreationDate, &uid, &skip, &comment.Pinned, &comment.User.ID, &comment.User.Uuid, &comment.User.ProfilePic, &comment.User.Username, &comment.User.Email, &comment.User.Password, &comment.User.Role, &comment.User.CreationDate, &comment.User.Biography, &comment.User.LastSeen)
				if err != nil {
					log.Fatal(err)
				}

				comment.CreationDate = comment.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				comment.User.CreationDate = comment.User.CreationDate.(time.Time).Format("01/02/2006 15:04:05")
				comment.User.LastSeen = comment.User.LastSeen.(time.Time).Format("01/02/2006 15:04:05")

				comment.Content = strings.Replace(comment.Content, "<br/>", "\r\n", -1)
				comment.Vote = ""

				c = append(c, comment)
			}
			row.Close()

			// Move comments that are pinned to the top
			var pinnedComments []f.Comment
			var unpinnedComments []f.Comment
			for _, comment := range c {
				if comment.Pinned == 1 {
					pinnedComments = append(pinnedComments, comment)
				} else {
					unpinnedComments = append(unpinnedComments, comment)
				}
			}

			comments := []f.Comment{}
			comments = append(comments, pinnedComments...)
			comments = append(comments, unpinnedComments...)

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

				row, err = db.Query("SELECT COUNT(*) FROM vote WHERE comment_id = ? AND type = 'dislike'", comments[i].ID)
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

				if cookie != nil {
					row, err = db.Query("SELECT COUNT(*) FROM vote WHERE comment_id = ? AND user_id = ? AND type = 'like'", comments[i].ID, userID)
					if err != nil {
						log.Fatal(err)
					}
					var userLikes int
					for row.Next() {
						err = row.Scan(&userLikes)
						if err != nil {
							log.Fatal(err)
						}
					}
					row.Close()
					if userLikes > 0 {
						comments[i].Vote = "like"
					}

					row, err = db.Query("SELECT COUNT(*) FROM vote WHERE comment_id = ? AND user_id = ? AND type = 'dislike'", comments[i].ID, userID)
					if err != nil {
						log.Fatal(err)
					}
					var userDislikes int
					for row.Next() {
						err = row.Scan(&userDislikes)
						if err != nil {
							log.Fatal(err)
						}
					}
					row.Close()
					if userDislikes > 0 {
						comments[i].Vote = "dislike"
					}
				}
			}

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
