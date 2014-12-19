package main

import (
	"encoding/json"
	"io/ioutil"
	//"log"
	"errors"
	"os"
	"sync"
	"time"
)

// Key is a Key in collection
type key struct {
	Name       string                 `json:"name"`
	Path       string                 `json:"-"`
	File       os.File                `json:"-"`
	LastAccess time.Time              `json:"lastAccess"`
	Content    map[string]interface{} `json:"content"`
	*sync.RWMutex
}

// Key return a new key
func (c *collection) newKey(name string) (*key, error) {
	path := c.Path + string(os.PathSeparator) + name
	k := &key{}
	k.Name = name
	k.Path = path
	k.RWMutex = &sync.RWMutex{}
	return k, nil
}

// collection is a subdirectory of fsdb path
type collection struct {
	Name string          `json:"name"`
	Path string          `json:"-"`
	Keys map[string]*key `json:"keys"`
}

// newCollection Returns a collection, populated with Items if the given path already exists.
func (db *FsDb) newCollection(name string) (*collection, error) {
	path := db.Path + string(os.PathSeparator) + name
	// If the path doesn't exist, create it. Return error if that fails.
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// The path didn't exist, create it.
			if err := os.Mkdir(path, 0777); err != nil {
				return nil, err
			}
		} else {
			// Some other error, return it.
			return nil, err
		}
	}

	// Create the collection
	c := &collection{}
	c.Path = path
	c.Keys = make(map[string]*key)

	// If the directory did exist before, create keys.
	if stat != nil && stat.IsDir() {
		collectionKeys, err := ioutil.ReadDir(c.Path)
		if err != nil {
			return nil, err
		}

		for _, f := range collectionKeys {
			if !f.IsDir() {
				k, err := c.newKey(f.Name())
				if err != nil {
					return nil, err
				}
				k.Name = f.Name()
				c.Keys[k.Name] = k
			}
		}
	}
	return c, nil
}

// FsDb is a filesystem Database that automatically caches files for a duration of time.
type FsDb struct {
	Path        string                 `json:"-"`
	Collections map[string]*collection `json:"collections"`
}

// NewFsDb will return a new FsDb, or error if it can not be opened
func NewFsDb(path string) (*FsDb, error) {
	db := &FsDb{}
	db.Path = path
	db.Collections = make(map[string]*collection)

	dbCollections, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range dbCollections {
		if f.IsDir() {
			c, err := db.newCollection(f.Name())
			if err != nil {
				return nil, err
			}
			c.Name = f.Name()
			db.Collections[f.Name()] = c
		}
	}

	return db, nil
}

// String will return the whole loaded database in json format (for debugging)
func (db *FsDb) String() string {
	b, err := json.Marshal(db)
	if err != nil {
		return ""
	}
	return "FSDB: " + string(b)
}

// read will return a map[string]interface{} or error
func (k *key) read() (map[string]interface{}, error) {
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

func (c *collection) read(key string) (map[string]interface{}, error) {
	if c.Keys[key] == nil {
		return nil, errors.New("key does not exist: " + key)
	}
	return c.Keys[key].read()
}

// Write will write a map[string]interface{} to a file, or error.
func (k *key) write(v map[string]interface{}) error {
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

// Write will write a map[string]interface{} to a file, or error.
func (c *collection) write(key string, v map[string]interface{}) error {
	if c.Keys[key] == nil {
		k, err := c.newKey(key)
		if err != nil {
			return err
		}
		c.Keys[key] = k
	}

	if err := c.Keys[key].write(v); err != nil {
		return err
	}
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

// delete will delete the key if it exist in this collection
func (c *collection) delete(key string) error {
	if c.Keys[key] == nil {
		return errors.New("key does not exist: " + key)
	}
	err := c.Keys[key].delete()
	if err != nil {
		return err
	}
	delete(c.Keys, key)
	return nil
}

// Write to key in collection
func (db *FsDb) Write(col, key string, v map[string]interface{}) error {
	if db.Collections[col] == nil {
		c, err := db.newCollection(col)
		if err != nil {
			return err
		}
		db.Collections[col] = c
	}
	return db.Collections[col].write(key, v)
}

// Read key from collection
func (db *FsDb) Read(col string, key string) (map[string]interface{}, error) {
	if db.Collections[col] == nil {
		return nil, errors.New("collection does not exist: " + col)
	}
	return db.Collections[col].read(key)
}

// Delete the key from collection
func (db *FsDb) Delete(col string, key string) error {
	if db.Collections[col] == nil {
		return errors.New("collection does not exist: " + col)
	}
	return db.Collections[col].delete(key)
}
