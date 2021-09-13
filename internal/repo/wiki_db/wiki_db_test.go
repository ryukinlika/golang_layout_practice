package wiki_db

import (
	"database/sql"
	"fmt"
	"golang_layout/internal/model/page_model"
	"testing"

	"github.com/stretchr/testify/assert"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
)

type DBInterfaceMock struct {
	mock.Mock
	pingRet     func() error
	queryRet    func() (*sql.Rows, error)
	queryRowRet func() *sql.Row
	execRet     func() (sql.Result, error)
	closeRet    func() error
}

func (db DBInterfaceMock) Ping() error {
	return db.pingRet()
}
func (db DBInterfaceMock) Query(string, ...interface{}) (*sql.Rows, error) {
	return db.queryRet()
}
func (db DBInterfaceMock) QueryRow(string, ...interface{}) *sql.Row {
	return db.queryRowRet()
}
func (db DBInterfaceMock) Exec(string, ...interface{}) (sql.Result, error) {
	return db.execRet()
}
func (db DBInterfaceMock) Close() error {
	return db.closeRet()
}

type SQLInterfaceMock struct {
	mock.Mock
	openRet func() (*sql.DB, error)
}

func (sql SQLInterfaceMock) Open(string, string) (*sql.DB, error) {
	return sql.openRet()
}

type TestCase struct {
	test_num  int
	test_name string
	dbMock    sqlmock.Sqlmock
	sqlMock   SQLInterfaceMock
}

var test_case_open = []TestCase{
	{
		test_num:  1,
		test_name: "Error opening connection",
		sqlMock: SQLInterfaceMock{
			openRet: func() (*sql.DB, error) { return nil, fmt.Errorf("error") },
		},
	},
	{
		test_num:  2,
		test_name: "Success opening connection",
		sqlMock: SQLInterfaceMock{
			openRet: func() (*sql.DB, error) {
				db, _, _ := sqlmock.New()
				return db, nil
			},
		},
	},
}

func TestDatabaseOpenFunction(t *testing.T) {
	for _, tc := range test_case_open {
		db = nil
		sqlObject = tc.sqlMock
		wiki := WikiRepo{}

		_, expected_err := tc.sqlMock.openRet()

		err := wiki.Open()

		assert.Equal(t, expected_err, err, tc.test_name)
	}
}

func TestDatabaseGetAllTitles_ErrorQuery(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("error select query"))

	expected_err := fmt.Errorf("error in select operation: error select query")

	db = db_mock
	_, err = wiki.GetAllTitles()

	assert.Equal(t, expected_err, err)

}
func TestDatabaseGetAllTitles_RowError(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(int64(1), "title").
		AddRow(int64(2), "title").
		RowError(1, fmt.Errorf("error"))

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	expected_err := fmt.Errorf("row error: error")

	db = db_mock
	_, err = wiki.GetAllTitles()

	assert.Equal(t, expected_err, err)

}

func TestDatabaseGetAllTitles_InvalidType(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow("eeeeeee", "title")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	expected_err := fmt.Errorf("error in row scan")

	db = db_mock
	_, err = wiki.GetAllTitles()

	assert.Equal(t, expected_err, err)

}

func TestDatabaseGetAllTitles_Success(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(int64(1), "title").
		AddRow(int64(2), "title")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	db = db_mock
	_, err = wiki.GetAllTitles()

	assert.Equal(t, nil, err)

}

func TestDatabaseGetById_Success(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(int64(1), "title", "body")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	db = db_mock
	_, err = wiki.GetById(int64(1))

	assert.Equal(t, nil, err)

}

func TestDatabaseGetById_NoRow(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "body"})

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	db = db_mock
	_, err = wiki.GetById(int64(1))

	assert.Equal(t, fmt.Errorf("pageId 1: not found"), err)

}

func TestDatabaseInsertPage_Success(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	mock.ExpectExec("INSERT INTO pages").WithArgs("title", "body").WillReturnResult(sqlmock.NewResult(2, 1))

	db = db_mock
	id, _ := wiki.InsertPage(&page_model.Page{
		Title: "title",
		Body:  "body",
	})

	assert.Equal(t, int64(2), id)

}

func TestDatabaseInsertPage_ErrorInsert(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	mock.ExpectExec("INSERT INTO pages").WithArgs("title", "body").WillReturnError(fmt.Errorf("error insert"))

	db = db_mock
	_, err = wiki.InsertPage(&page_model.Page{
		Title: "",
		Body:  "",
	})

	assert.Equal(t, fmt.Errorf("error insert"), err)

}

func TestDatabaseUpdatePage_Success(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	mock.ExpectExec("UPDATE pages").WithArgs("title", "body", int64(1)).WillReturnResult(sqlmock.NewResult(2, 1))

	db = db_mock
	_, err = wiki.UpdatePage(&page_model.Page{
		Id:    int64(1),
		Title: "title",
		Body:  "body",
	})

	assert.Equal(t, nil, err)

}

func TestDatabaseUpdatePage_ErrorInsert(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	mock.ExpectExec("UPDATE pages").WithArgs("title", "body", int64(1)).WillReturnError(fmt.Errorf("error update"))

	db = db_mock
	_, err = wiki.InsertPage(&page_model.Page{
		Id:    int64(1),
		Title: "123123",
		Body:  "213123",
	})

	assert.Equal(t, fmt.Errorf("error insert"), err)

}

func TestDatabaseDeletePage_Success(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	mock.ExpectExec("DELETE from pages").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))

	db = db_mock
	_, err = wiki.DeletePage(int64(1))

	assert.Equal(t, nil, err)

}
func TestDatabaseDeletePage_ErrorInsert(t *testing.T) {
	wiki := WikiRepo{}
	db_mock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db_mock.Close()

	mock.ExpectExec("DELETE from pages").WithArgs(int64(1)).WillReturnError(fmt.Errorf("error delete"))

	db = db_mock
	_, err = wiki.DeletePage(int64(1))

	assert.Equal(t, fmt.Errorf("error delete"), err)

}
