package forum

import (
	"html/template"
	"log"
	"net/http"
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
