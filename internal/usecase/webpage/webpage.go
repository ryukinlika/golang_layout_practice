package webpage

import (
	"golang_layout/internal/model/page_model"
	model "golang_layout/internal/model/page_model"
	"golang_layout/internal/repo/wiki_db"
	"html/template"
)

type WebPage struct {
	Wiki      wiki_db.WikiRepoInterface
	Templates *template.Template
}

func (web WebPage) Init() {
	web.Templates = template.Must(template.ParseFiles(page_model.Template_lists...))
}

func (web WebPage) LoadPage(id int64) (*model.Page, error) {
	//call db function
	page, err := web.Wiki.GetById(id)

	if err != nil {
		return nil, err
	}

	return page, nil
}

func (web WebPage) LoadHome() (*[]model.Page, error) {
	pages, err := web.Wiki.GetAllTitles()
	if err != nil {
		return nil, err
	}
	return &pages, err
}
func (web WebPage) Insert(title string, body string) (int64, error) {
	page := &model.Page{Title: title, Body: body}
	id, err := web.Wiki.InsertPage(page)

	return id, err
}

func (web WebPage) Update(id int64, title string, body string) error {
	page := &model.Page{Id: id, Title: title, Body: body}
	_, err := web.Wiki.UpdatePage(page)
	return err
}
