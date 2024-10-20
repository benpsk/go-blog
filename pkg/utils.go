package pkg

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type AuthUser struct {
	Id    int
	Name  string
	Email string
}

type PageData struct {
	Data   interface{}
	Error  string
	Errors map[string]string
	User   AuthUser
}

var layoutTmpl *template.Template

func ParseLayoutFiles() {
	layoutTmpl = template.Must(template.ParseFiles(filepath.Join("internal", "templates", "layouts", "layout.html")))
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data *PageData) {
	tmp, err := template.Must(layoutTmpl.Clone()).ParseFiles(templateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, _ := r.Context().Value("user").(AuthUser)
	data.User = user
	log.Println(user)
	if err := tmp.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
