package storages

type IFStorage interface {
	Save(string) (string, error)
	Load(string) (*Model, error)
	LoadInfo(string) (*Model, error)
	Close()
}

type Model struct {
	Url     string `json:"url"`
	Visited bool   `json:"visited"`
	Count   int    `json:"count"`
}
