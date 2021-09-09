package webpage

import (
	"golang_layout/internal/model/page_model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// type WikiRepoInterface interface {
// 	GetAllTitles() ([]page_model.Page, error)
// 	GetById(int64) (page_model.Page, error)
// 	InsertPage(*page_model.Page) (int64, error)
// 	UpdatePage(page *page_model.Page) (int64, error)
// 	Close()
// 	Open()
// }

type WikiRepoMock struct {
	mock.Mock
	titleRet  func() ([]page_model.Page, error)
	idRet     func(int64) (page_model.Page, error)
	insertRet func(*page_model.Page) (int64, error)
	updateRet func(*page_model.Page) (int64, error)
}

func (w WikiRepoMock) GetAllTitles() ([]page_model.Page, error) {
	return w.titleRet()
}
func (w WikiRepoMock) GetById(id int64) (page_model.Page, error) {
	return w.idRet(id)
}
func (w WikiRepoMock) InsertPage(p *page_model.Page) (int64, error) {
	return w.insertRet(p)
}
func (w WikiRepoMock) UpdatePage(p *page_model.Page) (int64, error) {
	return w.updateRet(p)
}
func (w WikiRepoMock) Open() {
}
func (w WikiRepoMock) Close() {
}

type Answer struct {
	page_model.Page
	error
}

type IdAnswer struct {
	int64
	error
}

type ArrayAnswer struct {
	[]page_model.Page
	error
}

func TestLoadPage(t *testing.T) {
	web := WebPage{}
	web.wiki = WikiRepoMock{
		idRet: func(id int64) (page_model.Page, error) {
			return page_model.Page{Id: id, Title: "title_test", Body: "test"}, nil
		},
	}

	expected := Answer{page_model.Page{Id: 1, Title: "title_test", Body: "test"}, nil}
	
	a, b := web.LoadPage(1)
	actual := Answer{*a, b}

	assert.Equal(t, expected, actual, "check load page")
}

func TestGetAll_Success(t *testing.T) {
	web := WebPage{}
	web.wiki = WikiRepoMock{
		titleRet: func() ([]page_model.Page, error) { return []page_model.Page{
			page_model.Page{Id: 1, Title:"a", Body:"aabc"},
			page_model.Page{Id: 2, Title:"b", Body:"babc"},
		}, nil },
	}

	
	expected := IdAnswer{int64(id_case), nil}
	
	a, b := web.Insert("", "")
	actual := IdAnswer{a, b}

}	
