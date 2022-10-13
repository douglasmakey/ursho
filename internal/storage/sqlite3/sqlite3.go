package sqlite3

import (
	"database/sql"

	"github.com/douglasmakey/ursho"
	"github.com/douglasmakey/ursho/internal/base62"

	_ "github.com/mattn/go-sqlite3" // sqlite engine
)

type sqlite3 struct {
	db *sql.DB
}

// New returns a sqlite backed storage service.
func New(path string) (ursho.ItemService, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	strQuery := "CREATE TABLE IF NOT EXISTS shortener (uid INTEGER PRIMARY KEY AUTOINCREMENT, url VARCHAR not NULL, " +
		"visited boolean DEFAULT FALSE, count INTEGER DEFAULT 0);"
	_, err = db.Exec(strQuery)
	if err != nil {
		return nil, err
	}

	return &sqlite3{db: db}, nil
}

func (s *sqlite3) Save(url string) (string, error) {
	var id int64

	stmt, err := s.db.Prepare("INSERT INTO shortener(url,visited,count) VALUES(?, ?, ?)")
	if err != nil {
		return "", err
	}

	res, err := stmt.Exec(url, 0, 0)
	if err != nil {
		return "", err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return "", err
	}

	return base62.Encode(id), nil
}

func (s *sqlite3) Load(code string) (string, error) {
	id, err := base62.Decode(code)
	if err != nil {
		return "", err
	}

	var url string
	err = s.db.QueryRow("update shortener set visited=true, count = count + 1 where uid=$1 RETURNING url", id).Scan(&url)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *sqlite3) LoadInfo(code string) (*ursho.Item, error) {
	id, err := base62.Decode(code)
	if err != nil {
		return nil, err
	}

	var item ursho.Item
	query := "SELECT url, visited, count FROM shortener where uid=$1 limit 1"
	err = s.db.QueryRow(query, id).Scan(&item.URL, &item.Visited, &item.Count)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *sqlite3) Close() error { return s.db.Close() }
