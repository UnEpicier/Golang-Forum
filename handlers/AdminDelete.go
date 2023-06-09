package forum

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

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
						row, err := db.Query("SELECT id FROM post WHERE user_id = ?", reqID)
						if err != nil {
							log.Fatal(err)
						}
						posts := []int{}
						for row.Next() {
							var postID int
							err = row.Scan(&postID)
							if err != nil {
								log.Fatal(err)
							}
							posts = append(posts, postID)
						}

						comments := []int{}
						for _, postID := range posts {
							_, err = db.Exec("DELETE FROM vote WHERE post_id = ?", postID)
							if err != nil {
								log.Fatal(err)
							}

							row, err = db.Query("SELECT id FROM comment WHERE post_id = ?", postID)
							if err != nil {
								log.Fatal(err)
							}
							for row.Next() {
								var commentID int
								err = row.Scan(&commentID)
								if err != nil {
									log.Fatal(err)
								}
								comments = append(comments, commentID)
							}

							_, err = db.Exec("DELETE FROM vote WHERE comment_id = ?", postID)
							if err != nil {
								log.Fatal(err)
							}
						}

						for _, commentID := range comments {
							_, err = db.Exec("DELETE FROM vote WHERE comment_id = ?", commentID)
							if err != nil {
								log.Fatal(err)
							}

							_, err = db.Exec("DELETE FROM comment WHERE id = ?", commentID)
							if err != nil {
								log.Fatal(err)
							}
						}

						_, err = db.Exec("DELETE FROM user WHERE id = ?", reqID)
						if err != nil {
							log.Fatal(err)
						}

						http.Redirect(w, r, "/", http.StatusFound)
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
