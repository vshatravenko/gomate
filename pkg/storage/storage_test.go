package storage

import (
	"bytes"
	"encoding/json"
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

type BasicStruct struct {
	Content []string
}

func TestTypes(t *testing.T) {
	var input1 = map[string]string{
		"Pink Floyd": "Time",
		"The Cure":   "Faith",
		"Deftones":   "MX",
		"Placebo":    "Lady of the Flowers",
	}

	output1 := make(map[string]string)

	testGetPut(t, input1, &output1)

	var input2 = BasicStruct{
		[]string{"there's", "more", "things", "to", "life"},
	}

	output2 := BasicStruct{}

	testGetPut(t, input2, &output2)
}

func testGetPut(t *testing.T, inVal interface{}, outVal interface{}) {
	// Clean up the test storage file
	os.Remove("storage-test.db")

	// Open the storage
	db, err := Open("storage-test.db")
	if err != nil {
		t.Fatal(err)
	}

	// Marshal the inVal
	input, err := json.Marshal(inVal)
	if err != nil {
		t.Fatal(err)
	}

	// Put inVal into the storage
	if err := db.Put("nyess", inVal); err != nil {
		t.Fatal(err)
	}

	// Get it out
	if err := db.Get("nyess", outVal); err != nil {
		t.Fatal(err)
	}

	// Marshal the outVal
	output, err := json.Marshal(outVal)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the inVal and outVal are equal
	if !bytes.Equal(input, output) {
		t.Fatal("input and output are not equal")
	}

	// Close the db
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	// Clean up the test storage file
	os.Remove("storage-test.db")
}
