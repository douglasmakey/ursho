package sqlite3

import (
	"database/sql"
	"log"

	"github.com/douglasmakey/ursho/base62"
	"github.com/douglasmakey/ursho/storage"
	_ "github.com/mattn/go-sqlite3" // sqlite engine
)

// New returns a sqlite backed storage service.
func New(FilePath string) (storage.Service, error) {
	db, err := sql.Open("sqlite3", FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Create table if not exists
	strQuery := "CREATE TABLE IF NOT EXISTS shortener (uid INTEGER PRIMARY KEY AUTOINCREMENT, url VARCHAR not NULL, " +
		"visited boolean DEFAULT FALSE, count INTEGER DEFAULT 0);"
	_, err = db.Exec(strQuery)
	if err != nil {
		return nil, err
	}
	return &sqlite3{FilePath}, nil
}

type sqlite3 struct {
	filePath string
}

func (s *sqlite3) Save(url string) (string, error) {
	var id int64
	db, err := sql.Open("sqlite3", s.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO shortener(url,visited,count) VALUES(?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(url, 0, 0)
	id, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return base62.Encode(id), nil
}

func (s *sqlite3) Load(code string) (*storage.Item, error) {
	db, err := sql.Open("sqlite3", s.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	id, err := base62.Decode(code)
	if err != nil {
		return nil, err
	}

	item, err := s.LoadInfo(code)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("update shortener set visited=$1, count=$2 where uid=$3", true, item.Count+1, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *sqlite3) LoadInfo(code string) (*storage.Item, error) {
	db, err := sql.Open("sqlite3", s.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	id, err := base62.Decode(code)
	if err != nil {
		return nil, err
	}

	var item storage.Item
	err = db.QueryRow("SELECT url, visited, count FROM shortener where uid=$1 limit 1", id).
		Scan(&item.URL, &item.Visited, &item.Count)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *sqlite3) Close() error { return nil }
