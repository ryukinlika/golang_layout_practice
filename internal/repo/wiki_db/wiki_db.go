package wiki_db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"

	"golang_layout/internal/model/page_model"
)

type WikiRepoInterface interface {
	GetAllTitles() ([]page_model.Page, error)
	GetById(int64) (*page_model.Page, error)
	InsertPage(page *page_model.Page) (int64, error)
	UpdatePage(page *page_model.Page) (int64, error)
	Close()
	Open()
}

var db *sql.DB

type WikiRepo struct {
}

func (w WikiRepo) Open() {

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
		db, err = sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			log.Fatal(err)
		}

		pingErr := db.Ping()
		if pingErr != nil {
			log.Fatal(pingErr)
		}
	}
	fmt.Println("Connected!")

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
			return nil, fmt.Errorf("error in row scan: %v", err)
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
		return 0, fmt.Errorf("addPage: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addPage: %v", err)
	}
	return id, nil

}

func (w WikiRepo) UpdatePage(page *page_model.Page) (int64, error) {
	w.Open()
	result, err := db.Exec("UPDATE pages SET title = ?, body=? WHERE id = ?", page.Title, page.Body, page.Id)
	if err != nil {
		return 0, fmt.Errorf("updatePage: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("updatePage: %v", err)
	}
	return id, nil

}
func (w WikiRepo) Close() {
	db.Close()
}
