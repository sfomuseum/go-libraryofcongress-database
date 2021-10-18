package www

import (
	"html/template"
	"net/http"
)

func renderTemplate(rsp http.ResponseWriter, t *template.Template, vars interface{}) {

	rsp.Header().Set("Content-type", "text/html")

	err := t.Execute(rsp, vars)

	if err != nil {
		http.Error(rsp, err.Error(), http.StatusInternalServerError)
	}
}
