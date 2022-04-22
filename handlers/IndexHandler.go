package forum

import (
	f "forum"
	"html/template"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/index.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page f.Page
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
