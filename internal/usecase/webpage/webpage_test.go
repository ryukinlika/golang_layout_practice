package webpage

import (
	"fmt"
	"golang_layout/internal/model/page_model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// type WikiRepoInterface interface {
// 	GetAllTitles() ([]page_model.Page, error)
// 	GetById(int64) (*page_model.Page, error)
// 	InsertPage(*page_model.Page) (int64, error)
// 	UpdatePage(page *page_model.Page) (int64, error)
// 	Close()
// 	Open()
// }

type WikiRepoMock struct {
	mock.Mock
	titleRet  func() ([]page_model.Page, error)
	idRet     func(int64) (*page_model.Page, error)
	insertRet func(*page_model.Page) (int64, error)
	updateRet func(*page_model.Page) (int64, error)
}

func (w WikiRepoMock) GetAllTitles() ([]page_model.Page, error) {
	return w.titleRet()
}
func (w WikiRepoMock) GetById(id int64) (*page_model.Page, error) {
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
	*page_model.Page
	error
}

type IdAnswer struct {
	int64
	error
}

type ArrayAnswer struct {
	objects []page_model.Page
	err     error
}

func TestLoadPage(t *testing.T) {
	web := WebPage{}
	web.Wiki = WikiRepoMock{
		idRet: func(id int64) (*page_model.Page, error) {
			return &page_model.Page{Id: id, Title: "title_test", Body: "test"}, nil
		},
	}

	expected := Answer{&page_model.Page{Id: 1, Title: "title_test", Body: "test"}, nil}

	a, b := web.LoadPage(1)
	actual := Answer{a, b}

	assert.Equal(t, expected, actual, "check load page")
}

func TestLoadPageId_Failed(t *testing.T) {
	web := WebPage{}
	web.Wiki = WikiRepoMock{
		idRet: func(id int64) (*page_model.Page, error) {
			return nil, fmt.Errorf("error in select operation")
		},
	}

	expected := Answer{nil, fmt.Errorf("error in select operation")}

	a, b := web.LoadPage(1)
	actual := Answer{a, b}

	assert.Equal(t, expected, actual, "check load fails")
}

func TestLoadHome_Success(t *testing.T) {
	web := WebPage{}
	web.Wiki = WikiRepoMock{
		titleRet: func() ([]page_model.Page, error) {
			return []page_model.Page{
				{Id: 1, Title: "a", Body: "aabc"},
				{Id: 2, Title: "b", Body: "babc"},
			}, nil
		},
	}

	expected := ArrayAnswer{[]page_model.Page{
		{Id: 1, Title: "a", Body: "aabc"},
		{Id: 2, Title: "b", Body: "babc"},
	}, nil}

	a, b := web.LoadHome()
	actual := ArrayAnswer{*a, b}

	assert.Equal(t, expected, actual, "check load home")

}

func TestLoadHome_Fail(t *testing.T) {
	web := WebPage{}
	web.Wiki = WikiRepoMock{
		titleRet: func() ([]page_model.Page, error) {
			return []page_model.Page{}, fmt.Errorf("Error in select operation")
		},
	}

	expected := ArrayAnswer{[]page_model.Page{}, fmt.Errorf("Error in select operation")}
	actual := ArrayAnswer{}
	a, b := web.LoadHome()
	if a == nil {
		actual = ArrayAnswer{[]page_model.Page{}, b}
	} else {
		actual = ArrayAnswer{*a, b}
	}

	assert.Equal(t, expected, actual, "check load home fails")

}

func TestInsert_Success(t *testing.T) {
	web := WebPage{}
	web.Wiki = WikiRepoMock{
		insertRet: func(*page_model.Page) (int64, error) {
			return 4, nil
		},
	}

	expected := IdAnswer{4, nil}

	a, b := web.Insert("abc", "abc")
	actual := IdAnswer{a, b}

	assert.Equal(t, expected, actual, "check insert success")

}

func TestInsert_Fail(t *testing.T) {
	web := WebPage{}
	web.Wiki = WikiRepoMock{
		insertRet: func(*page_model.Page) (int64, error) {
			return 0, fmt.Errorf("addPage error")
		},
	}

	expected := IdAnswer{0, fmt.Errorf("addPage error")}

	a, b := web.Insert("abc", "abc")
	actual := IdAnswer{a, b}

	assert.Equal(t, expected, actual, "check insert fails")

}

func TestUpdate_Fail(t *testing.T) {
	web := WebPage{}
	web.Wiki = WikiRepoMock{
		updateRet: func(*page_model.Page) (int64, error) {
			return 0, fmt.Errorf("updatePage error")
		},
	}

	actual := web.Update(1, "a", "a")
	assert.Equal(t, fmt.Errorf("updatePage error"), actual, "check update fails")

}

func TestUpdate_Success(t *testing.T) {
	web := WebPage{}
	web.Wiki = WikiRepoMock{
		updateRet: func(*page_model.Page) (int64, error) {
			return 1, nil
		},
	}

	actual := web.Update(1, "a", "a")
	assert.Equal(t, nil, actual, "check update fails")

}
