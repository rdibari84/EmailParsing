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

func IsDirectory(path string) (bool) {
    fileInfo, err := os.Stat(path)
    checkError(err, "Invalid outfile. Path '" + path + "' does not exist\n")
    return fileInfo.IsDir()
}

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

func readParseAndWriteFiles(dir string, outfile string) {

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
		wg.Wait()
		log.Println("done waiting")
		close(ch)
	}(wg, ch)

  // only write file if we parsed emails
  if numberParsedFiles > 0 {
    writeFile(ch, outfile)
  }
}

func writeFile(ch chan EmailInformation, outfile string) {
  f, err := os.Create(outfile)
  checkError(err, "Issue creating " + outfile)
  defer f.Close()

  w := bufio.NewWriter(f)
  for i := 0; i < numberParsedFiles; i++ {
    emailinfo := <-ch
    _, err := w.WriteString(emailinfo.From + "|" + emailinfo.Subject + "|" + emailinfo.Date + "|" + emailinfo.FileName + "\n")
    checkError(err, "Issue writting to file")
  }
  w.Flush()
  log.Printf("Succesfully wrote %s", outfile)
}

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
  var subjectIsDone = false
  // scan file
  scanner := bufio.NewScanner(file)
	for scanner.Scan() {
    line := scanner.Text()
    if from == "" { // haven't found from yet. keep looking
      tmp := findRegex("^From:", line)
      if tmp != "" {
          from = extractEmail(tmp)
      }
    }
    if date == "" { // haven't found date yet. keep looking
      date = findRegex("^Date:", line)
    }

    if subject == "" { // haven't found subject yet. keep looking
      subject = findRegex("^Subject:", line)
    } else if !subjectIsDone { // we have a subject
      // if line starts with a space, then its a continuation from above
      matched, _ := regexp.MatchString("^ [\\S\\s]*", line)
      if matched {
        subject += line
      } else { // if we have no more lines that start with a space, then the subject is done
        subjectIsDone = true
      }
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

func extractEmail(s string) string {
  p := regexp.MustCompile("[^<>][\\w.]*@[\\S]*\\.[\\w]{3}")
  toreturn := strings.TrimSpace(p.FindString(s))
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
