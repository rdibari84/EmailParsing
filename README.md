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
- dir points to the directory with emails to parse. If not provided, it defaults to the current directory. The code looks for any email file (.msg) in the specified dir. It fails if the directory passed to it doesn't exist and does nothing if it finds no emails to parse.
- outfile points to the path and name of which to write the output. If not provided, it defaults to the current directory. It fails if the directory passed to it doesn't exist or if the requested file doesn't end in a .txt.
```
cd $GOPATH/src
go run github.com/rdibari84/EmailParsing/email/email.go -dir=github.com/rdibari/EmailParsing/data/ -outfile=github.com/rdibari84/EmailParsing/emails.txt
```

### Run Unit Tests
```
cd $GOPATH/src
go test github.com/rdibari84/EmailParsing/email
```
