package page_model

type Page struct {
	Id    int64
	Title string
	Body  string
}

var Template_lists = []string{
	"web/template/edit.html",
	"web/template/home.html",
	"web/template/view.html",
	"web/template/add.html",
}
