package forum

import (
	"bytes"
	"database/sql"
	f "forum"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/user/profile.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page f.Page
	page.Logged = false
	page.Error = ""

	cookie, _ := r.Cookie("user")

	if cookie != nil {
		page.Logged = true
		var user f.User

		user.Posts = []f.Post{}
		user.Comments = []f.Comment{}

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
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
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
			var post f.Post
			post.User = f.User{}
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
		row, err = db.Query("SELECT * FROM comment AS c INNER JOIN post AS p ON c.post_id = p.id WHERE c.user_id = ? ORDER BY c.creation_date DESC", user.ID)
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			var skip int
			var comment f.Comment
			err = row.Scan(&comment.ID, &comment.Content, &comment.CreationDate, &skip, &skip, &comment.Pinned, &comment.PostID.ID, &comment.PostID.Title, &comment.PostID.Content, &comment.PostID.CreationDate, &skip, &comment.PostID.CategoryId, &comment.PostID.Pinned, &comment.PostID.LastUpdate)
			if err != nil {
				log.Fatal(err)
			}

			comment.CreationDate = comment.CreationDate.(time.Time).Format("02/01/2006 15:04:05")

			user.Comments = append(user.Comments, comment)
		}
		row.Close()

		if r.Method == "POST" {

			if r.FormValue("form") == "profilepic" {
				r.ParseMultipartForm(32 << 20)

				file, header, err := r.FormFile("pic")
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, file); err != nil {
					log.Fatal(err)
				}

				row, err = db.Query("SELECT id FROM user WHERE uuid = ?", cookie.Value)
				if err != nil {
					log.Fatal(err)
				}
				var id int
				for row.Next() {
					err = row.Scan(&id)
					if err != nil {
						log.Fatal(err)
					}
				}
				row.Close()

				matches, err := filepath.Glob("./static/assets/profile/" + strconv.Itoa(id) + ".*")
				if err != nil {
					log.Fatal(err)
				}

				if len(matches) > 0 {
					os.Remove(matches[0])
				}

				err = os.WriteFile("./static/assets/profile/"+strconv.Itoa(id)+filepath.Ext(header.Filename), buf.Bytes(), 0666)
				if err != nil {
					log.Fatal(err)
				}

				_, err = db.Exec("UPDATE user SET profile_pic = ? WHERE uuid = ?", "/assets/profile/"+strconv.Itoa(id)+filepath.Ext(header.Filename), cookie.Value)
				if err != nil {
					log.Fatal(err)
				}
				http.Redirect(w, r, "/user/profile", http.StatusFound)

			} else {
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

						if f.CheckPasswordhash(password, db_pass) {
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

							if f.CheckPasswordhash(password, db_pass) {
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

						if f.CheckPasswordhash(current, db_pass) {
							hash, err := f.HashPassword(new)
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

					if f.CheckPasswordhash(password, db_pass) {
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
