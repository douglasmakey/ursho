## Synopsis

..

## Code Example

Using CURL

Generate shortener\
`curl -H "Content-Type: application/json" -X POST -d '{"url":"www.google.com"}' http://localhost:8080/encode/`

Response:
`{"success":true,"response":"http://localhost:8080/p"}`

Redirect\
`curl http://localhost:8080/p  `

Get info for url shortener\
`curl http://localhost:8080/info/p  `

Response:
```
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

Use GLIDE Package Management for Golang, for installation all packages 

https://github.com/Masterminds/glide

Run `glide install` in the folder.

select the method of persistence, in which you going to work.\
`storage := &storages.Postgres{}`

Add your config for the method of persistence and other options in file config.json\
```
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