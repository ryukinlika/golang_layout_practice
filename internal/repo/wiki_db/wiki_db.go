package wiki_db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"golang_layout/internal/model/page_model"
)

type WikiRepoInterface interface {
	GetAllTitles() ([]page_model.Page, error)
	GetById(int64) (*page_model.Page, error)
	InsertPage(page *page_model.Page) (int64, error)
	UpdatePage(page *page_model.Page) (int64, error)
	DeletePage(int64) (int64, error)
	Close()
	Open() error
}

type DBInterface interface {
	Ping() error
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
	Close() error
}

type SQLInterface interface {
	// Open(string, string) (DBInterface, error)
	Open(string, string) (*sql.DB, error)
}

// var db DBInterface
var db *sql.DB
var sqlObject SQLInterface

type SQLStruct struct {
	openRet func(string, string) (*sql.DB, error)
}

// func (sql SQLStruct) Open(s1 string, s2 string) (DBInterface, error) {
// 	db, err := sql.openRet(s1, s2)
// 	return DBInterface(db), err
// }
func (sql SQLStruct) Open(s1 string, s2 string) (*sql.DB, error) {
	db, err := sql.openRet(s1, s2)
	return db, err
}

type WikiRepo struct {
}

func (w WikiRepo) Open() error {
	if sqlObject == nil {
		sqlObject = SQLStruct{
			openRet: sql.Open,
		}
	}
	// Get a database handle.
	var err error
	if db == nil { //not initiated
		cfg := mysql.Config{
			User:                 "root",
			Passwd:               "root",
			Net:                  "tcp",
			Addr:                 "127.0.0.1:3306",
			DBName:               "wikis",
			AllowNativePasswords: true,
		}
		db, err = sqlObject.Open("mysql", cfg.FormatDSN())
		if err != nil {
			return err
		}

		pingErr := db.Ping()
		if pingErr != nil {
			return pingErr
		}
	}
	fmt.Println("Connected!")
	return nil
}

func (w WikiRepo) GetAllTitles() ([]page_model.Page, error) {
	w.Open()
	var pages []page_model.Page
	rows, err := db.Query("SELECT id, title FROM pages")

	if err != nil {
		return nil, fmt.Errorf("error in select operation: %v", err)
	}
	for rows.Next() {
		var p page_model.Page
		if err := rows.Scan(&p.Id, &p.Title); err != nil {
			return nil, fmt.Errorf("error in row scan")
		}
		pages = append(pages, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %v", err)
	}

	return pages, err
}

func (w WikiRepo) GetById(id int64) (*page_model.Page, error) {
	w.Open()
	var page page_model.Page

	row := db.QueryRow("SELECT * FROM pages WHERE id = ?", id)
	if err := row.Scan(&page.Id, &page.Title, &page.Body); err != nil {
		if err == sql.ErrNoRows {
			return &page, fmt.Errorf("pageId %d: not found", id)
		}
		return &page, fmt.Errorf("pageId %d: %v", id, err)
	}
	return &page, nil
}

func (w WikiRepo) InsertPage(page *page_model.Page) (int64, error) {
	w.Open()
	result, err := db.Exec("INSERT INTO pages (title, body) VALUES (?, ?)", page.Title, page.Body)
	if err != nil {
		return 0, fmt.Errorf("error insert")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error insert")
	}
	return id, nil

}

func (w WikiRepo) UpdatePage(page *page_model.Page) (int64, error) {
	w.Open()
	result, err := db.Exec("UPDATE pages SET title = ?, body=? WHERE id = ?", page.Title, page.Body, page.Id)
	if err != nil {
		return 0, fmt.Errorf("error update")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error update")
	}
	return id, nil

}
func (w WikiRepo) DeletePage(id int64) (int64, error) {
	w.Open()
	result, err := db.Exec("DELETE from pages where id = ?", id)
	if err != nil {
		return 0, fmt.Errorf("error delete")
	}
	processed_id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error update")
	}
	return processed_id, nil

}

func (w WikiRepo) Close() {
	db.Close()
}
