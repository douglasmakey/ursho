package postgres

import (
	"database/sql"
	"fmt"

	// This loads the postgres drivers.
	_ "github.com/lib/pq"

	"github.com/douglasmakey/ursho/enconding"
	"github.com/douglasmakey/ursho/storage"
)

// New returns a postgres backed storage service.
func New(user, password, dbName string) (storage.Service, error) {
	// Coonect postgres
	connect := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		user, password, dbName)
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
	var id int
	err := p.db.QueryRow("INSERT INTO shortener(url,visited,count) VALUES($1,$2,$3) returning uid;", url, false, 0).Scan(&id)
	if err != nil {
		return "", err
	}
	return enconding.Encode(id), nil
}

func (p *postgres) Load(code string) (*storage.Item, error) {
	id := enconding.Decode(code)

	item, err := p.LoadInfo(code)
	if err != nil {
		return nil, err
	}

	_, err = p.db.Exec("update shortener set visited=$1, count=$2 where uid=$3", true, item.Count+1, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (p *postgres) LoadInfo(code string) (*storage.Item, error) {
	id := enconding.Decode(code)

	var item storage.Item
	err := p.db.QueryRow("SELECT url, visited, count FROM shortener where uid=$1 limit 1", id).
		Scan(&item.URL, &item.Visited, &item.Count)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (p *postgres) Close() error { return p.db.Close() }
