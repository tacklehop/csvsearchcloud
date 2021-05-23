package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Preprocessing")
	uri = "https://raw.githubusercontent.com/tacklehop/csvsearchcloud/main/sample.csv"
	//go cscMain(uri)

	status := m.Run()

	fmt.Println("Postprocessing")
	os.Exit(status)
}

func TestSearch(t *testing.T) {
	cases := []struct {
		name string
		in   string
		hit  bool
	}{
		{"case 1", "Key1", true},
		{"case 2", "Word3", true},
		{"case 3", "None", false},
	}

	for _, c := range cases {
		result, err := searchCsvFromHttp(uri, c.in)
		if err != nil {
			t.Errorf("Error: %v results in %v", c.in, err)
		}
		if (result == "" && c.hit) || (result != "" && !c.hit) {
			t.Errorf("%v results in %v", c.in, result)
		}
	}
}
