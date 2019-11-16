package main

import "testing"

func TestReadFile(t *testing.T) {
	f, error := readFile("read_file_test.go")
	if error != nil {
		t.Errorf("File should be readable")
	}
	if len(f) <= 0 {
		t.Errorf("data len should not be 0")
	}
}
