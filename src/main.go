package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/SeongJuMoon/tutorial-gowiki/pkg/domain"
)

var (
	templates = template.Must(template.ParseFiles("./template/html/edit.html", "./template/html/view.html"))
	vaildPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi There, I Love %s", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := domain.LoadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderStaticTemplateHtml(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := domain.LoadPage(title)
	if err != nil {
		p = &domain.Page{Title: title}
	}
	renderStaticTemplateHtml(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	content := r.FormValue("content")
	p := &domain.Page{Title: title, Content: []byte(content)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := vaildPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// helper:
func renderStaticTemplateHtml(w http.ResponseWriter, filename string, p *domain.Page) {
	err := templates.ExecuteTemplate(w, fmt.Sprintf("%s.html", filename), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	p1 := &domain.Page{Title: "TestPage", Content: []byte("아니 이게 맞네?")}
	p1.Save()
	p2, _ := domain.LoadPage("TestPage")
	fmt.Println(string(p2.Content))

	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
