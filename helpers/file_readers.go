package helpers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func listdir() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func ReadUidsFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		listdir()
		log.Fatal(fmt.Sprintf("read uids: %v", err))
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) > 0 && txt != "/n" {
			lines = append(lines, txt)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(fmt.Sprintf("read uids scan: %v", err))
	}
	return lines
}
