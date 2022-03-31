package forum

import (
	"log"
	"net/http"
	"text/template"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		IndexHandler(w, r)
		return
	}

	tplt := template.Must(template.ParseFiles("./static/error.html"))

	err := tplt.Execute(w, tplt)
	if err != nil {
		log.Fatal(err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/index.html", "./static/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	err := tplt.ExecuteTemplate(w, "base", tplt)
	if err != nil {
		log.Fatal(err)
	}
}
