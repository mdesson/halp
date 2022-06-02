package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Halp struct {
	Program     string `json:"program"`
	Command     string `json:"command"`
	Explanation string `json:"explanation"`
}

func (h Halp) String() string {
	return fmt.Sprintf("[%s] %s: %s", h.Program, h.Command, h.Explanation)
}

var halpers []Halp

func main() {
	// Read (or create) file to init halpers
	err := initJSON()
	if err != nil {
		log.Fatal(err)
	}

	// Get args
	program, keywords := parseArgs()

	if program == "" {
		// print all
		for _, h := range halpers {
			fmt.Println(h)
		}
		return
	} else if program == "add" {
		// add new entry
		add()
	} else if keywords == "" {
		// print all by program
		byProgram(program)
	} else {
		find(program, keywords)
	}
}

func parseArgs() (program, keywords string) {
	args := os.Args[1:]
	if len(args) > 0 {
		program = args[0]
	}
	if len(args) > 1 {
		keywords = strings.ToLower(strings.Join(args[1:], " "))
	}
	return program, keywords
}

func initJSON() error {
	configPath := getConfigPath()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
		return err
	}

	// Open file, create it it doesn't exist
	f, err := os.OpenFile(configPath+"/config.json", os.O_CREATE|os.O_RDONLY, 0644)
	defer f.Close()
	if err != nil {
		return err
	}

	// Return early if file is empty
	info, err := f.Stat()
	if err != nil {
		return err
	} else if info.Size() == 0 {
		return nil
	}

	// Read it into memory and store in halpers
	b := make([]byte, info.Size())
	if _, err = f.Read(b); err != nil {
		return err
	}

	if err = json.Unmarshal(b, &halpers); err != nil {
		return err
	}

	return nil
}

func byProgram(program string) {
	for _, h := range halpers {
		if h.Program == program {
			fmt.Println(h)
		}
	}
}

func find(program, keywords string) {
	for _, h := range halpers {
		searchString := strings.ToLower(h.Command + " " + h.Explanation)
		if h.Program == program && strings.Contains(searchString, keywords) {
			fmt.Println(h)
		}
	}
}

func add() {
	var h Halp
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter program: ")
	if scanner.Scan() {
		h.Program = scanner.Text()
	}

	fmt.Print("Enter command: ")
	if scanner.Scan() {
		h.Command = scanner.Text()
	}

	fmt.Print("Enter explanation: ")
	if scanner.Scan() {
		h.Explanation = scanner.Text()
	}

	halpers = append(halpers, h)

	f, err := json.Marshal(halpers)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(getConfigPath()+"/config.json", f, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added to library!")
}

func getConfigPath() string {
	// Build path to config file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s/.config/halp/", homeDir)
}
