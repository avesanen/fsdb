package fsdb

import "testing"
import "os"

type TestType struct {
	Name string `json:"name"`
}

func TestOpen(t *testing.T) {
	dbPath := os.TempDir() + string(os.PathSeparator) + "db"
	if err := os.RemoveAll(dbPath); err != nil {
		t.Fatal(err.Error())
	}

	db, err := NewFsDb("./db")
	if err != nil {
		t.Fatal("should not fail if db does exist.")
	}

	var tt1 TestType
	if err := db.Read("col", "key", tt1); err == nil {
		t.Fatal("read should have failed", err.Error())
	}

	tt2 := &TestType{}
	tt2.Name = "Hello world!"
	if err := db.Write("col", "key", &tt2); err != nil {
		t.Fatal("read failed", err.Error())
	}

	var tt3 TestType
	if err := db.Read("col", "key", &tt3); err != nil {
		t.Fatal("read failed", err.Error())
	}
	t.Log(tt3)

}
