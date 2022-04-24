package forum

import (
	"database/sql"
	"fmt"
	f "forum"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/admin/admin.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page f.Page
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
			admin := f.Admin{}

			/*
				STATS
			*/
			stats := f.Stats{}

			// Globals
			forum := f.AD_Forum{}
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
				c := f.AD_Categories{}
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
			sum_catAct := []f.AD_CatActivities{}
			for i := 0; i < len(CAT_ids); i++ {
				catAct := f.AD_CatActivities{}
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
			users := f.AD_Users{}
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
				inscr := f.AD_Inscription{}
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
				user := f.User{}
				err = row.Scan(&user.ID, &user.Uuid, &user.ProfilePic, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreationDate, &user.Biography, &user.LastSeen)
				if err != nil {
					log.Fatal(err)
				}
				user.CreationDate = user.CreationDate.(time.Time).Format("2006/01/02 15:04:05")
				user.LastSeen = user.LastSeen.(time.Time).Format("2006/01/02 15:04:05")

				admin.Users = append(admin.Users, user)
			}
			row.Close()

			for i := 0; i < len(admin.Users); i++ {
				row, err = db.Query("SELECT * FROM post WHERE user_id = ?", admin.Users[i].ID)
				if err != nil {
					log.Fatal(err)
				}
				for row.Next() {
					var p f.Post
					var skip interface{}
					err = row.Scan(&p.ID, &p.Title, &p.Content, &p.CreationDate, &skip, &p.CategoryId, &p.Pinned, &p.LastUpdate)
					if err != nil {
						fmt.Println("Posts")
						log.Fatal(err)
					}
					admin.Users[i].Posts = append(admin.Users[i].Posts, p)
				}
				row.Close()

				row, err = db.Query("SELECT * FROM comment WHERE user_id = ?", admin.Users[i].ID)
				if err != nil {
					log.Fatal(err)
					fmt.Println("Posts")
				}
				for row.Next() {
					var p f.Comment
					var skip interface{}
					err = row.Scan(&p.ID, &p.Content, &p.CreationDate, &skip, &skip, &p.Pinned)
					if err != nil {
						log.Fatal(err)
					}
					admin.Users[i].Comments = append(admin.Users[i].Comments, p)
				}
				row.Close()
			}

			if r.Method == "POST" {
				err := r.ParseForm()
				if err != nil {
					log.Fatal(err)
				}
				form := r.FormValue("form")

				if form == "role" {
					role := r.FormValue("promote")
					userID := r.FormValue("userID")
					_, err = db.Exec("UPDATE user SET role = ? WHERE id = ?", role, userID)
					if err != nil {
						log.Fatal(err)
					}
					http.Redirect(w, r, "/admin#users", http.StatusFound)
				}

			}
			/* REPORTS TAB */
			reports := []f.Report{}

			row, err = db.Query("SELECT * FROM report")
			if err != nil {
				log.Fatal(err)
			}

			uids := []interface{}{}
			pids := []interface{}{}
			cids := []interface{}{}
			for row.Next() {
				report := f.Report{}
				var userID interface{}
				var postID interface{}
				var commentID interface{}
				err = row.Scan(&report.ID, &report.Type, &report.Reason, &report.CreationDate, &userID, &postID, &commentID)
				if err != nil {
					log.Fatal(err)
				}
				report.CreationDate = report.CreationDate.(time.Time).Format("02/01/2006 15:04:05")

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
						err = row.Scan(&reports[i].User.ID, &reports[i].User.Uuid, &reports[i].User.Username, &reports[i].User.ProfilePic, &reports[i].User.Email, &reports[i].User.Password, &reports[i].User.Role, &reports[i].User.CreationDate, &reports[i].User.Biography, &reports[i].User.LastSeen)
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
						err = row.Scan(&reports[i].Post.ID, &reports[i].Post.Title, &reports[i].Post.Content, &reports[i].Post.CreationDate, &skip, &reports[i].Post.CategoryId, &reports[i].Post.Pinned, &reports[i].Post.LastUpdate, &reports[i].User.ID, &reports[i].User.Uuid, &reports[i].User.ProfilePic, &reports[i].User.Username, &reports[i].User.Email, &reports[i].User.Password, &reports[i].User.Role, &reports[i].User.CreationDate, &reports[i].User.Biography, &reports[i].User.LastSeen)
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
						err = row.Scan(&reports[i].Comment.ID, &reports[i].Comment.Content, &reports[i].Comment.CreationDate, &skip, &skip, &reports[i].Comment.Pinned, &reports[i].Post.ID, &reports[i].Post.Title, &reports[i].Post.Content, &reports[i].Post.CreationDate, &skip, &reports[i].Post.CategoryId, &reports[i].Post.Pinned, &reports[i].Post.LastUpdate, &reports[i].User.ID, &reports[i].User.Uuid, &reports[i].User.Username, &reports[i].User.ProfilePic, &reports[i].User.Email, &reports[i].User.Password, &reports[i].User.Role, &reports[i].User.CreationDate, &reports[i].User.Biography, &reports[i].User.LastSeen)
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
