## Synopsis

#### URL Shortener Service

## Code Example

Using CURL

Generate shortener\
`curl -H "Content-Type: application/json" -X POST -d '{"url":"www.google.com"}' http://localhost:8080/encode/`

Response:
`{"success":true,"response":"http://localhost:8080/1"}`

Redirect\
Open url in your browser [http://localhost:8080/1](http://localhost:8080/p)

OR\
`curl http://localhost:8080/1  `

Get info for url shortener\
`curl http://localhost:8080/info/1 `

Response:
```json
{
 "success":true,
 "response": {
    "url":"www.google.com",
    "visited":true,
    "count":1
 }
}
```

## Motivation

..

## Installation

We'll use github.com/user as our base path. Create a directory inside your workspace in which to keep source code:

***mkdir -p $GOPATH/src/github.com/douglasmakey cd "$_"***

Clone repository or download and unrar in folder\
```git clone https://github.com/douglasmakey/ursho.git```


Use GLIDE Package Management for Golang, for installation all packages 

https://github.com/Masterminds/glide

Run `glide install` in the folder.

If selected Postgresql as Storage, create database
```sql
CREATE DATABASE ursho_db;
```

select the method of persistence, in which you going to work.\
`storage := &storages.Postgres{}`

Add your config for the method of persistence and other options in file config.json\
```json
{
  "server": {
    "host": "0.0.0.0",
    "port": "8080"
  },
  "options": {
    "prefix": "http://localhost:8080/"
  },
  "posgres": {
    "user": "postgres",
    "password": "mysecretpassword",
    "db": "ursho_db"
  }
}
```
## API Reference

..

## Tests

..

## Contributors

..

## License

A short snippet describing the license (MIT, Apache, etc.)
