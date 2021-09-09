package page_handler

import (
	model "golang_layout/internal/model/page_model"
	webpage_lib "golang_layout/internal/usecase/webpage"
	"net/http"
	"regexp"
	"strconv"
)

var webpage = webpage_lib.WebPage{}

func viewHandler(w http.ResponseWriter, r *http.Request, id string) {
	nId, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		http.Error(w, "Id must be int", http.StatusBadRequest)
		return
	}
	p, err := webpage.LoadPage(nId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound) //shortcut to add new entry
	}
	webpage.RenderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, id string) {
	nId, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		http.Error(w, "Id must be int", http.StatusBadRequest)
		return
	}
	p, err := webpage.LoadPage(nId)
	if err != nil {
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}
	webpage.RenderTemplate(w, "edit", p)

}

func updateHandler(w http.ResponseWriter, r *http.Request, id string) { //verifies data
	nId, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		http.Error(w, "Id must be int", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	title := r.FormValue("title")
	body := r.FormValue("body")
	if len(body) == 0 || len(title) == 0 {
		http.Error(w, "Content must not be empty", http.StatusBadRequest)
		return
	}
	err = webpage.Update(nId, title, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+id, http.StatusFound)
}

func insertHandler(w http.ResponseWriter, r *http.Request, placeholder string) {
	r.ParseForm()
	title := r.FormValue("title")
	body := r.FormValue("body")
	if len(body) == 0 || len(title) == 0 {
		http.Error(w, "Content must not be empty", http.StatusBadRequest)
		return
	}
	id, err := webpage.Insert(title, body)
	strId := strconv.FormatInt(id, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+strId, http.StatusFound)
}

func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := webpage.LoadHome()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	webpage.RenderHome(w, p)
}

func addHandler(w http.ResponseWriter, r *http.Request, title string) {
	p := &model.Page{Title: title}
	webpage.RenderTemplate(w, "add", p)
}

var validPath = regexp.MustCompile("^/(edit|view|update)/([0-9]+)$") //regex for crud path

var homePath = regexp.MustCompile("^/(home|add|insert)/$") //regex for home and add path

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		//fmt.Println(r.URL.Path)
		if m == nil {
			n := homePath.FindStringSubmatch(r.URL.Path)
			if n == nil {
				if r.URL.Path == "/" {
					http.Redirect(w, r, "/home", http.StatusFound) //redirect "/" to "/home"
				} else {
					http.NotFound(w, r)
				}
				return
			}
			fn(w, r, "Title")
			return
		}
		fn(w, r, m[2])
	}
}

func CreateHandlers() {
	http.HandleFunc("/", makeHandler(homeHandler))
	http.HandleFunc("/home/", makeHandler(homeHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/update/", makeHandler(updateHandler))
	http.HandleFunc("/insert/", makeHandler(insertHandler))
	http.HandleFunc("/add/", makeHandler(addHandler))

}