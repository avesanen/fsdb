package fsdb

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// key struct stores a key and manages loading, saving and deleting it
type key struct {
	Name       string    `json:"name"`
	Path       string    `json:"-"`
	File       os.File   `json:"-"`
	LastAccess time.Time `json:"lastAccess"`
	Content    *[]byte   `json:"content"`
	*sync.RWMutex
}

// read will read the key to v
func (k *key) read(v interface{}) error {
	k.RLock()
	defer k.RUnlock()
	k.LastAccess = time.Now()

	b, err := ioutil.ReadFile(k.Path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &v)
}

// write will write the v to key
func (k *key) write(v interface{}) error {
	k.Lock()
	defer k.Unlock()
	k.LastAccess = time.Now()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(k.Path, b, 0666)
}

// delete will lock the file so nothing will read it, delete the file and
// retun success so collection can remove it.
func (k *key) delete() error {
	k.Lock()
	defer k.Unlock()
	if err := os.Remove(k.Path); err != nil {
		return err
	}
	return nil
}
