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

// Make a wait group wrapper to count how many items are in the waitgroup at a specific point in time
type WaitGroupWrapper struct {
  Wg sync.WaitGroup
  NumberInWaitGroup int
}

// wrap the WaitGroup Add in order to increment the NumberInWaitGroup
func (wgw *WaitGroupWrapper) Add(delta int) {
  wgw.Wg.Add(delta)
  wgw.NumberInWaitGroup++
}

// wrap the WaitGroup Done in order to decrement the NumberInWaitGroup
func (wgw *WaitGroupWrapper) Done() {
  wgw.Wg.Done()
  wgw.NumberInWaitGroup--
  log.Printf("Go routine finished. still waiting on %d routines to finish", wgw.NumberInWaitGroup)
}

// define contructor
func NewWaitGroupWrapper() WaitGroupWrapper {
    wgw := WaitGroupWrapper{NumberInWaitGroup: 0}
    return wgw
}

 // VERY IMP to declare this globally, other wise one would hit "fatal error: all goroutines are asleep - deadlock!"
var wgw = NewWaitGroupWrapper()

//////////////////////////////////////////////////////////////
//////////////////// Helper Methods //////////////////////////
//////////////////////////////////////////////////////////////

// Common error handling method
func checkError(e error, msg string) {
    if e != nil {
        log.Fatal(msg, e)
        panic(e)
    }
}


func IsDirectory(path string) (bool) {
    fileInfo, err := os.Stat(path)
    checkError(err, "Invalid outfile. Path '" + path + "' does not exist\n")
    return fileInfo.IsDir()
}

/* validates input into Main
makes sure the directory to read the emails from the the place to write the outfile are valid
throws exception if not */
func validateInput(dir string, outfile string) string {
  // double check directory. Add a trailing '/' if not present
  log.Printf("directory: %s\n", dir)
  if !strings.HasSuffix(dir, "/") {
    dir = dir + "/"
  }

  // checkout outfile. make sure its a .txt file and that the path to file exists
  log.Printf("outfile: %s\n", outfile)
  if !strings.HasSuffix(outfile, ".txt") {
    log.Fatal("Invalid outfile. Must have be a .txt\n")
  }
  IsDirectory(filepath.Dir(outfile))
  return dir
}

/* The brains of the function
reads the input directory for email files
adds to the wait group and kicks off goroutines
waits for the routines to finish */
func readParseAndWriteFiles(dir string, outfile string) {
  // get all files in input dir
  files, err := ioutil.ReadDir(dir)
  checkError(err, "Issue reading directory " + dir)

  ch := make(chan EmailInformation)
  for _, file := range files {
    // check if the files are emails. if so, add to waitgroup and kick off parsing
    if !file.IsDir() && strings.HasSuffix(file.Name(), ".msg") {
        wgw.Add(1)
        go emailParsing(dir + file.Name(), ch)
    } else if !file.IsDir() {   // files are not emails. do nothing
      log.Printf("Skipping %s. It isn't an email file\n", file.Name())
    }
  }

  /* capture the total number of files being parsed.
  wgw.NumberInWaitGroup will decrement as the goroutines finish */
  totalNumberEmails := wgw.NumberInWaitGroup
  log.Printf("found %d email files in directory to parse\n", totalNumberEmails)

  /* wait for the goroutines finish */
  go func(wgw WaitGroupWrapper, ch chan EmailInformation) {
		wgw.Wg.Wait()
		log.Println("done waiting")
		close(ch)
	}(wgw, ch)

  // only write file if we parsed emails
  if totalNumberEmails > 0 {
    writeFile(ch, outfile, totalNumberEmails)
  }
}

func writeFile(ch chan EmailInformation, outfile string, totalNumberEmails int) {
  f, err := os.Create(outfile)
  checkError(err, "Issue creating " + outfile)
  defer f.Close()

  w := bufio.NewWriter(f)
  for i := 0; i < totalNumberEmails; i++ {
    emailinfo := <-ch
    _, err := w.WriteString(emailinfo.From + "|" + emailinfo.Subject + "|" + emailinfo.Date + "|" + emailinfo.FileName + "\n")
    checkError(err, "Issue writting to file")
  }
  w.Flush()
  log.Printf("Succesfully wrote %s", outfile)
}

func emailParsing(myfile string, ch chan<-EmailInformation) {
  log.Printf("Opening %s\n", myfile)
  defer wgw.Done() // 3

  file, err := os.Open(myfile)
  checkError(err, "Error opening file "+ myfile + "\n")
  defer file.Close()
  filename := filepath.Base(myfile)
  info := parseFile(filename, file)

  ch <- info
}

func parseFile(filename string, file io.Reader) EmailInformation {
  // initialize variables
  var from = "||"
  var subject = "||"
  var date = "||"
  var subjectIsDone = false

  // scan file
  scanner := bufio.NewScanner(file)
	for scanner.Scan() {
    line := scanner.Text()

    // haven't found from yet since its equal to its default value. keep looking
    if from == "||" {
      tmp := findRegex("^From:", line)
      if tmp != "||" {
          from = extractEmail(tmp)
      }
    }

    // haven't found date yet since its equal to its default value. keep looking
    if date == "||" {
      date = findRegex("^Date:", line)
    }

    // haven't found subject yet since its equal to its default value. keep looking
    if subject == "||" {
      subject = findRegex("^Subject:", line)
    } else if !subjectIsDone { // we have a subject!
      // if the next line starts with a space, then its a continuation of the subject
      matched, _ := regexp.MatchString("^ [\\S\\s]*", line)
      if matched {
        subject += line // append the line to the subject
      } else { // The following line does not start with a space- then the subject is done
        subjectIsDone = true
      }
    }

    // if we have all the information, stop checking
    if from != "||" && date != "||" && subject != "||" && subjectIsDone {
      break;
    }
	}
	err := scanner.Err();
  checkError(err, "Issue scanning the file ")

  // make the object to return
  s := EmailInformation{from, subject, date, filename}
  log.Printf("{from: %s, subject: %s, date: %s, filename: %s }\n", s.From, s.Subject, s.Date, s.FileName)
  return s
}

/* Checks if the pattern is contained in the string and returns the right most slice.
Returns || if pattern is contained in the string */
func findRegex(pattern string, s string) string {
  matched, err := regexp.MatchString(pattern, s)
  checkError(err, "Issue checking text for " + pattern)
  if matched {
    a := regexp.MustCompile(pattern)
    toreturn := strings.TrimSpace(a.Split(s, -1)[1])
    return toreturn
  }
  return "||"
}

/* Checks if the string contains an email address.
ignores <>@?=utf-8 or anything else that might be in the string
returns || if it can't extract an email */
func extractEmail(s string) string {
  p := 	regexp.MustCompile("[\\w-_.]+@[\\w_-]+?\\.[a-zA-Z-._]+")
  toreturn := strings.TrimSpace(p.FindString(s))
  if toreturn == "" {
    toreturn = "||"
  }
  return toreturn
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
  // define dir command line flag. Default dir to current directory if not set.
  var outfile string
  flag.StringVar(&outfile, "outfile", currentdir + "/email.txt", "full directory path that holds the email files to parse")
  flag.Parse()


  dir = validateInput(dir, outfile)
  // parse files
  readParseAndWriteFiles(dir, outfile)
}
