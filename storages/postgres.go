package storages

import (
	"database/sql"
	"fmt"

	"github.com/douglasmakey/ursho/config"
	_ "github.com/lib/pq"

	"github.com/douglasmakey/ursho/enconding"
)

type Postgres struct {
	DB    *sql.DB
	model Model
}

func (p *Postgres) Init(config config.Config) error {
	// Coonect postgres
	strConnect := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.Postgres.User, config.Postgres.Password, config.Postgres.DB)
	db, err := sql.Open("postgres", strConnect)
	if err != nil {
		return err
	}

	// Ping to connection
	err = db.Ping()
	if err != nil {
		return err
	}

	// Create table if not exists
	strQuery := "CREATE TABLE IF NOT EXISTS shortener (uid serial NOT NULL, url VARCHAR not NULL, " +
		"visited boolean DEFAULT FALSE, count INTEGER DEFAULT 0);"

	_, err = db.Exec(strQuery)
	if err != nil {
		return err
	}
	// Set db in Model
	p.DB = db

	return nil
}

func (p *Postgres) Save(url string) (string, error) {

	var lastInsertId int
	err := p.DB.QueryRow("INSERT INTO shortener(url,visited,count) VALUES($1,$2,$3) returning uid;", url, false, 0).Scan(&lastInsertId)
	if err != nil {
		return "", err
	}
	fmt.Println("last inserted id =", lastInsertId)

	return enconding.Encode(lastInsertId), nil
}

func (p *Postgres) Load(code string) (Model, error) {
	// Decode code
	decodeID := enconding.Decode(code)

	// Query select
	err := p.DB.QueryRow("SELECT url, visited, count FROM shortener where uid=$1 limit 1",
		decodeID).Scan(&p.model.Url, &p.model.Visited, &p.model.Count)
	if err != nil {
		return Model{}, err
	}

	// Query update
	stmt, err := p.DB.Prepare("update shortener set visited=$1, count=$2 where uid=$3")
	if err != nil {
		return Model{}, err
	}
	_, err = stmt.Exec(true, p.model.Count+1, decodeID)
	return p.model, err
}

func (p *Postgres) LoadInfo(code string) (Model, error) {
	// Decode code
	decodeID := enconding.Decode(code)

	// Query select
	err := p.DB.QueryRow("SELECT url, visited, count FROM shortener where uid=$1 limit 1",
		decodeID).Scan(&p.model.Url, &p.model.Visited, &p.model.Count)
	if err != nil {
		return Model{}, err
	}

	return p.model, err
}

func (p *Postgres) Close() {
	p.DB.Close()
}