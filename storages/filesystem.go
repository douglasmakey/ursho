package storages

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/douglasmakey/ursho/config"
	"github.com/douglasmakey/ursho/enconding"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Filesystem struct {
	Root string
	sync.RWMutex
	model Model
}

func (s *Filesystem) Init(config config.Config) error {
	// Validate Path if not empty
	if config.Filesystem.Path == "" {
		return errors.New("Filesystem fail config")
	}

	s.Root = config.Filesystem.Path
	return os.MkdirAll(config.Filesystem.Path, 0744)
}

func (s *Filesystem) Code() string {
	s.Lock()
	files, _ := ioutil.ReadDir(s.Root)
	s.Unlock()

	return enconding.Encode(int(len(files) + 1))
}

func (s *Filesystem) Save(url string) (string, error) {
	code := s.Code()

	s.Lock()
	url = fmt.Sprintf(`{ "url": "%s", "visited": %t, "count": %d }`, url, false, 0)
	err := ioutil.WriteFile(filepath.Join(s.Root, code), []byte(url), 0744)
	if err != nil {
		return "", err
	}
	s.Unlock()

	return code, nil
}

func (s *Filesystem) Load(code string) (Model, error) {
	s.Lock()
	urlBytes, err := ioutil.ReadFile(filepath.Join(s.Root, code))
	s.Unlock()
	if err != nil {
		return Model{}, err
	}
	err = json.Unmarshal(urlBytes, &s.model)
	if err != nil {
		panic(err)
	}
	info := fmt.Sprintf(`{ "url": "%s", "visited": %t, "count": %d }`, s.model.Url, true, s.model.Count+1)
	s.Lock()
	err = ioutil.WriteFile(filepath.Join(s.Root, code), []byte(info), 0744)
	s.Unlock()
	if err != nil {
		panic(err)
	}

	return s.model, nil
}

func (s *Filesystem) LoadInfo(code string) (Model, error) {
	s.Lock()
	urlBytes, err := ioutil.ReadFile(filepath.Join(s.Root, code))
	s.Unlock()

	json.Unmarshal(urlBytes, &s.model)

	return s.model, err
}

func (s *Filesystem) Close() {

}
