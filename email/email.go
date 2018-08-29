package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
    "log"
    "path/filepath"
    "strings"
    "sync"
)

var wg sync.WaitGroup // 1

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func readFileGo(myfile string, ch chan<-[]string) {
  defer wg.Done() // 3
  fmt.Printf("Opening %s\n", myfile)

  file, err := os.Open(myfile)
  check(err)
  defer file.Close()

  scanner := bufio.NewScanner(file)
  var from = ""
  var subject = ""
  var date = ""
	for scanner.Scan() {
    line := scanner.Text()
    if strings.HasPrefix(line, "From: ") {
      from = line[len("From: "):]
    }
    if strings.HasPrefix(line, "Date: ") {
      date = line[len("Date: "):]
    }
    if strings.HasPrefix(line, "Subject: ") {
      subject = line[len("Subject: "):]
    }
    if from != "" && date != "" && subject != "" {
      break;
    }
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error Reading file " + myfile + "\n", err)
	}
  ch <- []string{from, subject, date}
}

func main() {
  dir, err := filepath.Abs(filepath.Dir("./data/"))
  if err != nil {
    log.Fatal(err)
  }

  files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

  ch := make(chan[]string)
  fmt.Printf("found %d files\n", len(files))
	for _, file := range files {
    if !file.IsDir() {
        wg.Add(1)
        go readFileGo(dir + "/" + file.Name(), ch)
    }
	}

  wg.Wait() // 4
  for i := 0; i < len(files); i++ {
    s := <-ch
    fmt.Printf("{from: %s, subject: %s, date: %s }\n", s[0], s[1], s[2])
  }

}
