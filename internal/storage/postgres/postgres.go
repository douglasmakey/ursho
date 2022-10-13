package postgres

import (
	"database/sql"
	"fmt"

	"github.com/douglasmakey/ursho"
	"github.com/douglasmakey/ursho/internal/base62"

	// This loads the postgres drivers.
	_ "github.com/lib/pq"
)

// New returns a postgres backed storage service.
func New(host, port, user, password, dbName string) (ursho.ItemService, error) {
	// Connect postgres
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", connect)
	if err != nil {
		return nil, err
	}

	// Ping to connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	strQuery := "CREATE TABLE IF NOT EXISTS shortener (uid serial NOT NULL, url VARCHAR not NULL, " +
		"visited boolean DEFAULT FALSE, count INTEGER DEFAULT 0);"

	_, err = db.Exec(strQuery)
	if err != nil {
		return nil, err
	}
	return &postgres{db}, nil
}

type postgres struct{ db *sql.DB }

func (p *postgres) Save(url string) (string, error) {
	var id int64
	err := p.db.QueryRow("INSERT INTO shortener(url,visited,count) VALUES($1,$2,$3) returning uid;", url, false, 0).Scan(&id)
	if err != nil {
		return "", err
	}
	return base62.Encode(id), nil
}

func (p *postgres) Load(code string) (string, error) {
	id, err := base62.Decode(code)
	if err != nil {
		return "", err
	}

	var url string
	err = p.db.QueryRow("update shortener set visited=true, count = count + 1 where uid=$1 RETURNING url", id).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (p *postgres) LoadInfo(code string) (*ursho.Item, error) {
	id, err := base62.Decode(code)
	if err != nil {
		return nil, err
	}

	var item ursho.Item
	err = p.db.QueryRow("SELECT url, visited, count FROM shortener where uid=$1 limit 1", id).
		Scan(&item.URL, &item.Visited, &item.Count)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (p *postgres) Close() error { return p.db.Close() }
