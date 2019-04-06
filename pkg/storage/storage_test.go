package storage

import (
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	os.Remove("storage-test.db")

	db, err := Open("storage-test.db")
	if err != nil {
		t.Fatal(err)
	}

	// put smth
	if err := db.Put("little", "molly"); err != nil {
		t.Fatal(err)
	}

	// get it out
	var val string
	if err := db.Get("little", &val); err != nil {
		t.Fatal(err)
	} else if val != "molly" {
		t.Fatalf("got %q, expected \"molly\"", val)
	}

	// put smth else with the same key
	if err := db.Put("little", "big"); err != nil {
		t.Fatal(err)
	}

	// get it
	if err := db.Get("little", &val); err != nil {
		t.Fatal(err)
	} else if val != "big" {
		t.Fatalf("got %q, expected \"big\"", val)
	}

	// get something nonexistent
	if err := db.Get("worldPeace", &val); err != ErrNotFound {
		if err != nil {
			t.Fatal(err)
		}

		t.Fatalf("got %q, expected absence", val)
	}

	// delete a key
	if err := db.Delete("little"); err != nil {
		t.Fatal(err)
	}

	// delete it again
	if err := db.Delete("little"); err != ErrNotFound {
		if err != nil {
			t.Fatal(err)
		}

		t.Fatalf("got %q, expected absence", val)
	}

	// close the db
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	os.Remove("storage-test.db")
}
