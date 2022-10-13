package ursho

type Item struct {
	URL     string `json:"url"`
	Visited bool   `json:"visited"`
	Count   int    `json:"count"`
}

type ItemService interface {
	Save(string) (string, error)
	Load(string) (string, error)
	LoadInfo(string) (*Item, error)
	Close() error
}
