package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

type Note struct {
	Text string
}

type Args struct {
	Method string // possible methods: add, delete, read
	Text   string
}

func ParseArgs() (Args, error) {
	args := os.Args
	if len(args) == 1 {
		return Args{}, nil
	}

	method := args[1]
	possibleMethods := []string{"add", "delete", "read"}
	if !slices.Contains(possibleMethods, method) {
		return Args{}, fmt.Errorf("Invalid method provided. Possible methods: add, delete, read")
	}

	if method == "add" && len(args) < 3 {
		return Args{}, fmt.Errorf("No text provided for adding a note")
	}

	return Args{Method: method, Text: strings.Join(args[2:], " ")}, nil
}

type File struct {
	Path string
}

func (f *File) Add(note Note) error {
	file, err := os.OpenFile(f.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("\n" + note.Text)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Read(command ReadCommand) error {
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	var notes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		notes = append(notes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	input := strings.Join(notes, "\n")
	command(input)
	return nil
}

func (f *File) Delete(command DeleteCommand) error {
	file, err := os.OpenFile(f.Path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var notes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		notes = append(notes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	input := strings.Join(notes, "\n")
	output := command(input)

	selectedNotes := strings.Split(output, "\n")
	if len(selectedNotes) == 0 || (len(selectedNotes) == 1 && selectedNotes[0] == "") {
		fmt.Println("No note selected")
		return nil
	}

	var updatedNotes []string
	for _, note := range notes {
		if !slices.Contains(selectedNotes, note) {
			updatedNotes = append(updatedNotes, note)
		}
	}

	//point to the beginning of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	//remove old content
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.WriteString(strings.Join(updatedNotes, "\n"))
	if err != nil {
		return err
	}

	return nil
}

type ReadCommand func(input string)

func FzfReadCommand(input string) {
	var copyCmd string

	if runtime.GOOS == "darwin" {
		copyCmd = "pbcopy"
	} else {
		copyCmd = "xclip -selection clipboard"
	}

	cmd := exec.Command("sh", "-c", `echo "`+input+`" | fzf --ansi --multi --bind "enter:execute(echo {} | `+copyCmd+`)+abort"`)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stderr = os.Stderr

	cmd.Run()
}

type DeleteCommand func(nput string) string

func FzfDeleteCommand(input string) string {
	cmd := exec.Command("fzf", "--ansi", "--multi", "--bind", "ctrl-s:toggle-sort", "--preview", "cat {}")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(out)

}

func main() {
	args, err := ParseArgs()
	if err != nil {
		panic(err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	path := filepath.Join(home, "notes.md")
	file := File{Path: path}
	switch args.Method {
	case "add":
		note := Note{Text: args.Text}
		err := file.Add(note)
		if err != nil {
			panic(err)
		}

	case "read":
		file.Read(FzfReadCommand)

	case "delete":
		file.Delete(FzfDeleteCommand)

	}

}
