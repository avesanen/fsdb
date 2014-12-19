package fsdb

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Key is a Key in collection
type key struct {
	Name       string      `json:"name"`
	Path       string      `json:"-"`
	File       os.File     `json:"-"`
	LastAccess time.Time   `json:"lastAccess"`
	Content    interface{} `json:"content"`
	*sync.RWMutex
}

// read will return a interface{} or error
func (k *key) read() (interface{}, error) {
	k.RLock()
	defer k.RUnlock()
	k.LastAccess = time.Now()

	if k.Content != nil {
		return k.Content, nil
	}
	k.LastAccess = time.Now()

	file, err := os.Open(k.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	if err = dec.Decode(&k.Content); err != nil {
		return nil, err
	}

	return k.Content, nil
}

// Write will write a interface{} to a file, or error.
func (k *key) write(v interface{}) error {
	k.Lock()
	defer k.Unlock()
	k.LastAccess = time.Now()

	// Open file for writing
	file, err := os.OpenFile(k.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(v); err != nil {
		return err
	}
	k.Content = v
	return nil
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
