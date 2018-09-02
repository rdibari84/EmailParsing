# EmailParsing


### Setup
```
echo $GOPATH # should be set to $HOME/go
mkdir -p $GOPATH/src/github.com/{{github-user}}
cd $GOPATH/src/github.com/{{github-user}}
git clone https://github.com/rdibari84/EmailParsing.git
```

### Build Code
```
cd $GOPATH/src
go install github.com/rdibari84/EmailParsing/email
```

### Run Unit Tests
- note unit tests use httptest to test api
- also tests concurrent connections
```
cd $GOPATH/src
go test github.com/rdibari84/EmailParsing/email
```
