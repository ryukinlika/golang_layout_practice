package page_handler

import (
	"fmt"
	"golang_layout/internal/model/page_model"
	"golang_layout/internal/repo/wiki_db"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

type WebPageMock struct {
	mock.Mock
}

func (w WebPageMock) Init() {
}

func (web WebPageMock) LoadPage(id int64) (*page_model.Page, error) {
	args := web.Called(id)
	return args.Get(0).(*page_model.Page), args.Error(1)
}

func (web WebPageMock) LoadHome() (*[]page_model.Page, error) {
	args := web.Called()
	return args.Get(0).(*[]page_model.Page), args.Error(1)
}
func (web WebPageMock) Insert(title string, body string) (int64, error) {
	args := web.Called(title, body)
	return args.Get(0).(int64), args.Error(1)
}

func (web WebPageMock) Update(id int64, title string, body string) error {
	args := web.Called(id, title, body)
	return args.Error(0)
}

func (web WebPageMock) Delete(id int64) error {
	args := web.Called(id)
	return args.Error(0)
}

func (web WebPageMock) AddWiki(w wiki_db.WikiRepo) {

}

func (web WebPageMock) Open() {
}

func (web WebPageMock) ExecuteTemplate(w io.Writer, tmpl string, p interface{}) error {
	args := web.Called(w, tmpl, p)
	return args.Error(0)
}

func TestViewHandler_Success(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("LoadPage", int64(1)).Return(&page_model.Page{Id: 1, Title: "Title", Body: "Body"}, nil)
	webMock.On("ExecuteTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	webpage = webMock

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/view/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	viewHandler(rr, req, "1")

	webMock.AssertExpectations(t)

}

func TestViewHandler_InvalidInput(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/view/1abc", nil)
	if err != nil {
		t.Fatal(err)
	}

	viewHandler(rr, req, "1abc")

	assert.Equal(t, http.StatusBadRequest, rr.Code)

}

func TestViewHandler_NotFound(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("LoadPage", int64(99)).Return(&page_model.Page{}, fmt.Errorf("pageId 99: not found"))
	webpage = webMock

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/view/99", nil)
	if err != nil {
		t.Fatal(err)
	}

	viewHandler(rr, req, "99")

	webMock.AssertExpectations(t)
}

func TestEditHandler_Success(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("LoadPage", int64(1)).Return(&page_model.Page{Id: 1, Title: "Title", Body: "Body"}, nil)
	webMock.On("ExecuteTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	webpage = webMock

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ediit/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	editHandler(rr, req, "1")

	webMock.AssertExpectations(t)

}

func TestEditHandler_InvalidInput(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/edit/1abc", nil)
	if err != nil {
		t.Fatal(err)
	}

	editHandler(rr, req, "1abc")

	assert.Equal(t, http.StatusBadRequest, rr.Code)

}

func TestEditHandler_NotFound(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("LoadPage", int64(99)).Return(&page_model.Page{}, fmt.Errorf("pageId 99: not found"))
	webpage = webMock

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/edit/99", nil)
	if err != nil {
		t.Fatal(err)
	}

	editHandler(rr, req, "99")

	webMock.AssertExpectations(t)
}

func TestUdpateHandler_Success(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("Update", int64(1), "new_title", "new_body").Return(nil)
	webpage = webMock

	rr := httptest.NewRecorder()

	form := url.Values{}
	form.Add("title", "new_title")
	form.Add("body", "new_body")

	req, err := http.NewRequest("POST", "/update/1", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	updateHandler(rr, req, "1")

	assert.Equal(t, http.StatusFound, rr.Code)

}

func TestUdpateHandler_InvalidInput_Id(t *testing.T) {
	rr := httptest.NewRecorder()

	form := url.Values{}
	form.Add("title", "new_title")
	form.Add("body", "new_body")

	req, err := http.NewRequest("POST", "/update/1", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	updateHandler(rr, req, "1abcsdwdwdas")

	assert.Equal(t, http.StatusBadRequest, rr.Code)

}
func TestUdpateHandler_InvalidInput_Form(t *testing.T) {
	rr := httptest.NewRecorder()

	form := url.Values{}
	form.Add("title", "")
	form.Add("body", "")

	req, err := http.NewRequest("POST", "/update/1", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	updateHandler(rr, req, "1")

	assert.Equal(t, http.StatusBadRequest, rr.Code)

}

func TestUdpateHandler_DatabaseError(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("Update", int64(1), "new_title", "new_body").Return(fmt.Errorf("updatePage: internal error"))
	webpage = webMock
	rr := httptest.NewRecorder()

	form := url.Values{}
	form.Add("title", "new_title")
	form.Add("body", "new_body")

	req, err := http.NewRequest("POST", "/update/1", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	updateHandler(rr, req, "1")

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

}

func TestInsertHandler_Success(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("Insert", "new_title", "new_body").Return(int64(5), nil)
	webpage = webMock

	rr := httptest.NewRecorder()

	form := url.Values{}
	form.Add("title", "new_title")
	form.Add("body", "new_body")

	req, err := http.NewRequest("POST", "/insert/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	insertHandler(rr, req, "1")

	assert.Equal(t, http.StatusFound, rr.Code)

}

func TestInsertHandler_DatabaseError(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("Insert", "new_title", "new_body").Return(int64(0), fmt.Errorf("addPage: Error"))
	webpage = webMock

	rr := httptest.NewRecorder()

	form := url.Values{}
	form.Add("title", "new_title")
	form.Add("body", "new_body")

	req, err := http.NewRequest("POST", "/insert/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	insertHandler(rr, req, "1")

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

}
func TestInsertHandler_InvalidForm(t *testing.T) {
	rr := httptest.NewRecorder()

	form := url.Values{}
	form.Add("title", "")
	form.Add("body", "")

	req, err := http.NewRequest("POST", "/insert/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	insertHandler(rr, req, "1")

	assert.Equal(t, http.StatusBadRequest, rr.Code)

}

func TestHomeHandler_Success(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("LoadHome").Return(&[]page_model.Page{{Id: 1, Title: "Title", Body: ""}}, nil)
	webMock.On("ExecuteTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	webpage = webMock

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/home/", nil)
	if err != nil {
		t.Fatal(err)
	}

	homeHandler(rr, req, "home")

	webMock.AssertExpectations(t)

}

func TestHomeHandler_DatabaseError(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("LoadHome").Return(&[]page_model.Page{}, fmt.Errorf("row error: error"))
	webpage = webMock

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/home/", nil)
	if err != nil {
		t.Fatal(err)
	}

	homeHandler(rr, req, "home")

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

}

func TestAddHandler(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("ExecuteTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	webpage = webMock

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/add/", nil)
	if err != nil {
		t.Fatal(err)
	}
	addHandler(rr, req, "add")

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAddHandler_TemplateFails(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("ExecuteTemplate", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("template error"))
	webpage = webMock

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/add/", nil)
	if err != nil {
		t.Fatal(err)
	}
	addHandler(rr, req, "add")

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHomeHandler_TemplateFails(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("LoadHome").Return(&[]page_model.Page{{Id: 1, Title: "Title", Body: ""}}, nil)
	webMock.On("ExecuteTemplate", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("template error"))
	webpage = webMock

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/add/", nil)
	if err != nil {
		t.Fatal(err)
	}
	homeHandler(rr, req, "add")

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandlerAssignment_Home(t *testing.T) {
	webMock := WebPageMock{}
	webMock.On("LoadHome").Return(&[]page_model.Page{{Id: 1, Title: "Title", Body: ""}}, nil)
	webMock.On("ExecuteTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	webpage = webMock

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/home/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := makeHandler(homeHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
