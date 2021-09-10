package webpage

import (
	"golang_layout/internal/model/page_model"
	"golang_layout/internal/repo/wiki_db"
	"html/template"
	"io"
)

var wiki wiki_db.WikiRepoInterface
var templates *template.Template

type WebPage struct {
	// Wiki      wiki_db.WikiRepoInterface
	// Templates *template.Template
}

type WebPageInterface interface {
	Init()
	LoadPage(int64) (*page_model.Page, error)
	LoadHome() (*[]page_model.Page, error)
	Insert(string, string) (int64, error)
	Update(int64, string, string) error
	AddWiki(wiki_db.WikiRepo)
	Open()
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func (web WebPage) Init() {
	templates = template.Must(template.ParseFiles(page_model.Template_lists...))
}

func (web WebPage) LoadPage(id int64) (*page_model.Page, error) {
	//call db function
	page, err := wiki.GetById(id)

	if err != nil {
		return nil, err
	}

	return page, nil
}

func (web WebPage) LoadHome() (*[]page_model.Page, error) {
	pages, err := wiki.GetAllTitles()
	if err != nil {
		return nil, err
	}
	return &pages, err
}
func (web WebPage) Insert(title string, body string) (int64, error) {
	page := &page_model.Page{Title: title, Body: body}
	id, err := wiki.InsertPage(page)

	return id, err
}

func (web WebPage) Update(id int64, title string, body string) error {
	page := &page_model.Page{Id: id, Title: title, Body: body}
	_, err := wiki.UpdatePage(page)
	return err
}

func (web WebPage) AddWiki(w wiki_db.WikiRepo) {
	wiki = w
}

func (web WebPage) Open() {
	wiki.Open()
}

func (web WebPage) ExecuteTemplate(w io.Writer, tmpl string, p interface{}) error {
	return templates.ExecuteTemplate(w, tmpl, p)
}
