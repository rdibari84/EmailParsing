# EmailParsing


### Setup
```
echo $GOPATH # should be set to $HOME/go
mkdir -p $GOPATH/src/github.com/rdibari84
cd $GOPATH/src/github.com/rdibari84
git clone https://github.com/rdibari84/EmailParsing.git
```

### Build Code
```
cd $GOPATH/src
go install github.com/rdibari84/EmailParsing/email
```

### Run 
Arguments:
* `-dir` points to the directory with emails (`.msg` files) to parse.  
  - Defaults to the current directory.  
  - Fails if the directory passed to it doesn't exist 
  - Does nothing if it finds no emails to parse.
  - Does not handle recursive folders
* `-outfile` points to the path and name of which to write the output. 
  - Defaults to the current directory. 
  - Fails if the directory passed to it doesn't exist OR if the requested file doesn't end in a `.txt` file format.
```
cd $GOPATH/src
go run github.com/rdibari84/EmailParsing/email/email.go -dir=github.com/rdibari/EmailParsing/data/ -outfile=github.com/rdibari84/EmailParsing/emails.txt
```

### Run Unit Tests
```
cd $GOPATH/src
go test --cover github.com/rdibari84/EmailParsing/email
```
