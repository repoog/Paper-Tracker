package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Paper struct {
	Title      string
	Title_CN   string
	Published  string
	Updated    string
	Link       string
	Summary    string
	Summary_CN string
}

type Database struct {
	*sql.DB
}

func ConnectDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) CreateTable() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS papers (
			title TEXT PRIMARY KEY,
			title_cn TEXT,
			published TEXT,
			updated TEXT,
			link TEXT,
			summary TEXT,
			summary_cn TEXT
		)
	`)

	return err
}

func (db *Database) InsertPaper(paper *Paper) error {
	_, err := db.Exec(`
		INSERT OR REPLACE INTO papers (title, title_cn, published, updated, link, summary, summary_cn)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
		paper.Title, paper.Title_CN, paper.Published, paper.Updated, paper.Link, paper.Summary, paper.Summary_CN)

	return err
}

func (db *Database) PaperExists(title string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM papers WHERE title LIKE ?", title).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
