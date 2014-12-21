# FsDb

FsDb is a small pet project that uses filesystem with json files as database.


## Installation

```bash
go get github.com/avesanen/fsdb
```

## Usage


```golang
package main

import (
	"github.com/avesanen/fsdb"
	"log"
)

type User struct {
	Id   string
	Name string
}

func main() {
	u1 := &User{}
	u1.Id = "0001"
	u1.Name = "admin"

	// Create or open a new database at "./db"
	db, _ := fsdb.NewFsDb("./db")

	// Write u1 to db
	if err := db.Write("users", u1.Id, &u1); err != nil {
		panic(err)
	}

	// Read key "0001" from collection "users".
	var u2 User
	if err := db.Read("users", "0001", &u2); err != nil {
		panic(err)
	}

	log.Println(u2)
}
```