package main

import (
	"os"
	"testing"
)

func TestParseArgs(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }() // Restore original os.Args after test

	os.Args = []string{"main"}
	args, err := ParseArgs()
	if args != (Args{}) || err != nil {
		t.Errorf("Expected nil for no arguments, got %v", args)
	}

	os.Args = []string{"main", "invalidMethod"}
	_, err = ParseArgs()
	if err == nil {
		t.Errorf("Expected error for invalid method, got nil")
	}

	os.Args = []string{"main", "add"}
	_, err = ParseArgs()
	if err == nil {
		t.Errorf("Expected error for no text provided, got nil")
	}

	os.Args = []string{"main", "add", "test note"}
	args, err = ParseArgs()
	if err != nil {
		t.Errorf("Expected no error for valid arguments, got %v", err)
	}
	if args.Method != "add" || args.Text != "test note" {
		t.Errorf("Expected method 'add' and text 'test note', got method '%s' and text '%s'", args.Method, args.Text)
	}
}

func TestFileAdd(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	f := File{Path: tmpfile.Name()}
	note := Note{Text: "test note"}

	err = f.Add(note)
	if err != nil {
		t.Errorf("Expected no error for Add, got %v", err)
	}

	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "test note\n" {
		t.Errorf("Expected 'test note\\n' in file, got '%s'", string(content))
	}
}

func TestFileDelete(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	f := File{Path: tmpfile.Name()}
	note := Note{Text: "test note"}

	err = f.Add(note)
	if err != nil {
		t.Errorf("Expected no error for Add, got %v", err)
	}

	note2 := Note{Text: "test note2"}
	err = f.Add(note2)
	if err != nil {
		t.Errorf("Expected no error for Add, got %v", err)
	}

	mockCommand := func(input string) string {
		return note2.Text

	}

	f.Delete(mockCommand)

	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	got := string(content)
	if got != note.Text {
		t.Errorf("Expected %v in the file, got %v", note.Text, string(content))
	}

}
