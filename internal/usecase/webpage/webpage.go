package webpage

import (
	model "golang_layout/internal/model/page_model"
	"golang_layout/internal/repo/wiki_db"
	"html/template"
	"net/http"
)

type WebPage struct {
	wiki      wiki_db.WikiRepoInterface
	templates *template.Template
}

func (web WebPage) LoadPage(id int64) (*model.Page, error) {
	//call db function
	page, err := web.wiki.GetById(id)

	if err != nil {
		return nil, err
	}

	return &page, nil
}

func (web WebPage) LoadHome() (*[]model.Page, error) {
	pages, err := web.wiki.GetAllTitles()
	if err != nil {
		return nil, err
	}
	return &pages, err
}
func (web WebPage) Insert(title string, body string) (int64, error) {
	page := &model.Page{Title: title, Body: body}
	id, err := web.wiki.InsertPage(page)

	return id, err
}

func (web WebPage) Update(id int64, title string, body string) error {
	page := &model.Page{Id: id, Title: title, Body: body}
	_, err := web.wiki.UpdatePage(page)
	return err
}

func (web WebPage) RenderTemplate(w http.ResponseWriter, tmpl string, p *model.Page) {
	web.templates = template.Must(template.ParseFiles(
		model.Template_lists...,
	))
	err := web.templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (web WebPage) RenderHome(w http.ResponseWriter, p *[]model.Page) {
	web.templates = template.Must(template.ParseFiles(
		model.Template_lists...,
	))
	err := web.templates.ExecuteTemplate(w, "home.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
