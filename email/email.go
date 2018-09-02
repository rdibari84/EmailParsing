package main

import (
    "bufio"
    "io/ioutil"
    "os"
    "log"
    "path/filepath"
    "strings"
    "sync"
    "flag"
    "io"
    "regexp"
)

// Define a datastructure to keep parsed email information
type EmailInformation struct {
  From string
  Subject string
  Date string
  FileName string
}

var numberParsedFiles int // define globally as to know how many files to wait for
var wg sync.WaitGroup // VERY IMP to declare this globally, other wise one   //would hit "fatal error: all goroutines are asleep - deadlock!"

//////////////////////////////////////////////////////////////
//////////////////// Helper Methods //////////////////////////
//////////////////////////////////////////////////////////////

func checkError(e error, msg string) {
    if e != nil {
        log.Fatal(msg, e)
        panic(e)
    }
}

func readAndParseFiles(dir string) {
  files, err := ioutil.ReadDir(dir)
  checkError(err, "Issue reading directory " + dir)

  ch := make(chan EmailInformation)
  for _, file := range files {
    if !file.IsDir() && strings.HasSuffix(file.Name(), ".msg") {
        numberParsedFiles++
        wg.Add(1)
        go emailParsing(dir + file.Name(), ch)
    } else if !file.IsDir() {
      log.Printf("Skipping %s. It isn't an email file\n", file.Name())
    }
  }

  log.Printf("found %d email files in directory to parse\n", numberParsedFiles)

  go func(wg sync.WaitGroup, ch chan EmailInformation) {
		log.Println("waiting")
		wg.Wait()
		log.Println("done waiting")
		close(ch)
	}(wg, ch)

  for i := 0; i < numberParsedFiles; i++ {
    s := <-ch
    log.Printf("{from: %s, subject: %s, date: %s, filename: %s }\n", s.From, s.Subject, s.Date, s.FileName)
  }
}

func writeFile() {}

func emailParsing(myfile string, ch chan<-EmailInformation) {
  log.Printf("Opening %s\n", myfile)
  defer wg.Done() // 3

  file, err := os.Open(myfile)
  checkError(err, "Error opening file "+ myfile + "\n")
  defer file.Close()
  filename := filepath.Base(myfile)
  info := parseFile(filename, file)

  ch <- info
}

func parseFile(filename string, file io.Reader) EmailInformation {
  // initialize variables
  var from = ""
  var subject = ""
  var date = ""

  // scan file
  scanner := bufio.NewScanner(file)
	for scanner.Scan() {
    line := scanner.Text()

    if from == "" { // haven't found from yet. keep looking
      from = findRegex("From:", line)
    }
    if date == "" { // haven't found date yet. keep looking
      date = findRegex("Date:", line)
    }
    if subject == "" { // haven't found subject yet. keep looking
      subject = findRegex("Subject:", line)
    }
    if from != "" && date != "" && subject != "" { // if we have all the information, stop checking
      break;
    }
	}
	err := scanner.Err();
  checkError(err, "Issue scanning the file ")

  // if we're missing any information set the feild to ||
  if from == "" {
    from = "||"
  }
  if subject == "" {
    subject = "||"
  }
  if date == "" {
    date = "||"
  }

  s := EmailInformation{from, subject, date, filename}
  log.Printf("{from: %s, subject: %s, date: %s, filename: %s }\n", s.From, s.Subject, s.Date, s.FileName)
  return s
}

func findRegex(pattern string, s string) string {
  matched, err := regexp.MatchString(pattern, s)
  checkError(err, "Issue checking text for " + pattern)
  if matched {
    a := regexp.MustCompile(pattern)
    toreturn := strings.TrimSpace(a.Split(s, -1)[1])
    return toreturn
  }
  return ""
}

//////////////////////////////////////////////////////////////
////////////////////////// Main //////////////////////////////
//////////////////////////////////////////////////////////////

func main() {
  // figure out current directory
  currentdir, err := filepath.Abs(filepath.Dir("."))
  if err != nil {
    log.Fatal(err)
  }

  // define dir command line flag. Default dir to current directory if not set.
  var dir string
  flag.StringVar(&dir, "dir", currentdir, "full directory path that holds the email files to parse")
  flag.Parse()

  // double check directory. Add a trailing '/' if not present
  log.Printf("directory: %s\n", dir)
  if !strings.HasSuffix(dir, "/") {
    dir = dir + "/"
  }

  // parse files
  readAndParseFiles(dir)
}
