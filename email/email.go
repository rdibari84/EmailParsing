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

var wg sync.WaitGroup // VERY IMP to declare this globally, other wise one   //would hit "fatal error: all goroutines are asleep - deadlock!"

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func readFile(dir string) (*os.File,string) {
  file, err := os.Open(dir)
  check(err)
  filename := filepath.Base(dir)
  return file, filename
}

func parseFile(file *os.File) []string {
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
    if from != "" && date != "" && subject != "" { // if we have all the information, stop checking
      break;
    }
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error Reading file %v. error %v \n", file, err)
	}
  // if we're missing any information retun ||
  if from == "" {
    from = "||"
  }
  if subject == "" {
    subject = "||"
  }
  if date == "" {
    date = "||"
  }
  return []string{from, subject, date}
}

func getInformationGo(myfile string, ch chan<-[]string) {
  defer wg.Done() // 3
  fmt.Printf("Opening %s\n", myfile)

  file, filename := readFile(myfile)

  info := parseFile(file)

  ch <- []string{info[0], info[1], info[2], filename}
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
        go getInformationGo(dir + "/" + file.Name(), ch)
    }
	}

  go func(wg sync.WaitGroup, ch chan []string) {
  		log.Println("waiting")
  		wg.Wait()
  		log.Println("done waiting")
  		close(ch)
  	}(wg, ch)

  for i := 0; i < len(files); i++ {
    s := <-ch
    fmt.Printf("{from: %s, subject: %s, date: %s, filename: %s }\n", s[0], s[1], s[2], s[3])
  }

}
