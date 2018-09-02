package main

import (
  "testing"
  "strings"
  "bufio"
  "log"
  "sync"
  "path/filepath"
)

var validEmail = "Return-Path: <infos@contact-darty.com>\n "+
"X-Original-To: 1000mercis@cp.assurance.returnpath.net\n "+
"Delivered-To: assurance@localhost.returnpath.net\n "+
"Received: from mxa-d1.returnpath.net (unknown [10.8.2.117])\n "+
	"by cpa-d1.returnpath.net (Postfix) with ESMTP id 426E5198271\n "+
	"for <1000mercis@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:17:45 -0600 (MDT)\n "+
"Received: from smtp-8-ft1.mm.fr.colt.net (smtp-7-ft1.mm.fr.colt.net [62.23.8.162])\n "+
	"by mxa-d1.returnpath.net (Postfix) with ESMTP id 2906A1CD\n "+
	"for <1000mercis@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:17:44 -0600 (MDT)\n "+
"Received: from host.25.62.23.62.rev.coltfrance.com ([62.23.62.25]:62162 helo=contact-darty.com)\n "+
	"by massmail-ft1.infra.coltfrance.com with esmtp (Exim)\n "+
	"id 1Q5fAU-00030S-4i\n "+
	"for <1000mercis@cp.assurance.returnpath.net>; Fri, 01 Apr 2011 16:17:42 +0200\n "+
"From: \"Darty\" <infos@contact-darty.com>\n "+
"To: 1000mercis@cp.assurance.returnpath.net\n "+
"Subject: Cuit Vapeur 29.90 euros, Nintendo 3DS 239 euros, GPS TOM TOM 139 euros... decouvrez VITE tous les bons plans du weekend !\n "+
"Date: 01 Apr 2011 16:17:41 +0200\n "+
"Message-ID: <20110401161739.E3786358A9D7B977@contact-darty.com>\n" +
"MIME-Version: 1.0\n" +
"x-idmail: DartyCRM_322_385774_10000\n" +
"Content-Type: text/html;\n" +
	"charset=\"iso-8859-1\"\n" +
"Content-Transfer-Encoding: 7bit\n"


var missingFrom = "Return-Path: <infos@contact-darty.com>\n "+
"X-Original-To: 1000mercis@cp.assurance.returnpath.net\n "+
"Delivered-To: assurance@localhost.returnpath.net\n "+
"Received: from mxa-d1.returnpath.net (unknown [10.8.2.117])\n "+
	"by cpa-d1.returnpath.net (Postfix) with ESMTP id 426E5198271\n "+
	"for <1000mercis@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:17:45 -0600 (MDT)\n "+
"Received: from smtp-8-ft1.mm.fr.colt.net (smtp-7-ft1.mm.fr.colt.net [62.23.8.162])\n "+
	"by mxa-d1.returnpath.net (Postfix) with ESMTP id 2906A1CD\n "+
	"for <1000mercis@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:17:44 -0600 (MDT)\n "+
"Received: from host.25.62.23.62.rev.coltfrance.com ([62.23.62.25]:62162 helo=contact-darty.com)\n "+
	"by massmail-ft1.infra.coltfrance.com with esmtp (Exim)\n "+
	"id 1Q5fAU-00030S-4i\n "+
	"for <1000mercis@cp.assurance.returnpath.net>; Fri, 01 Apr 2011 16:17:42 +0200\n "+
"To: 1000mercis@cp.assurance.returnpath.net\n "+
"Subject: Cuit Vapeur 29.90 euros, Nintendo 3DS 239 euros, GPS TOM TOM 139 euros... decouvrez VITE tous les bons plans du weekend !\n "+
"Date: 01 Apr 2011 16:17:41 +0200\n "+
"Message-ID: <20110401161739.E3786358A9D7B977@contact-darty.com>\n"

var missingAll = "Return-Path: <infos@contact-darty.com>\n "+
"X-Original-To: 1000mercis@cp.assurance.returnpath.net\n "+
"Delivered-To: assurance@localhost.returnpath.net\n "+
"Received: from mxa-d1.returnpath.net (unknown [10.8.2.117])\n "+
	"by cpa-d1.returnpath.net (Postfix) with ESMTP id 426E5198271\n "+
	"for <1000mercis@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:17:45 -0600 (MDT)\n "+
"Received: from smtp-8-ft1.mm.fr.colt.net (smtp-7-ft1.mm.fr.colt.net [62.23.8.162])\n "+
	"by mxa-d1.returnpath.net (Postfix) with ESMTP id 2906A1CD\n "+
	"for <1000mercis@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:17:44 -0600 (MDT)\n "+
"Received: from host.25.62.23.62.rev.coltfrance.com ([62.23.62.25]:62162 helo=contact-darty.com)\n "+
	"by massmail-ft1.infra.coltfrance.com with esmtp (Exim)\n "+
	"id 1Q5fAU-00030S-4i\n "+
	"for <1000mercis@cp.assurance.returnpath.net>; Fri, 01 Apr 2011 16:17:42 +0200\n "+
"Message-ID: <20110401161739.E3786358A9D7B977@contact-darty.com>\n"

func TestParseFileWithSpecificFields(t *testing.T){
  r1 := strings.NewReader("From: From\nSubject: Subject\nDate: Date\n")
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "From") {
    t.Errorf("Expected emailInfo.From to return 'From' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != "Subject") {
    t.Errorf("Expected emailInfo.Subject to return 'Subject' but got %s instead", emailInfo.Subject)
  }
  if (emailInfo.Date != "Date") {
    t.Errorf("Expected emailInfo.Date to return 'Date' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

func TestParseFileWithSpaces(t *testing.T){
  r1 := strings.NewReader(" From: From \n Subject: Subject \n Date:  Date \n")
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "From") {
    t.Errorf("Expected emailInfo.From to return 'From' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != "Subject") {
    t.Errorf("Expected emailInfo.Subject to return 'Subject' but got %s instead", emailInfo.Subject)
  }
  if (emailInfo.Date != "Date") {
    t.Errorf("Expected emailInfo.Date to return 'Date' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

func TestParseFileWithValidEmail(t *testing.T){
  r1 := strings.NewReader(validEmail)
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "\"Darty\" <infos@contact-darty.com>") {
    t.Errorf("Expected emailInfo.From to return 'From' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != "Cuit Vapeur 29.90 euros, Nintendo 3DS 239 euros, GPS TOM TOM 139 euros... decouvrez VITE tous les bons plans du weekend !") {
    t.Errorf("Expected emailInfo.Subject to return 'Subject' but got %s instead", emailInfo.Subject)
  }
  if (emailInfo.Date != "01 Apr 2011 16:17:41 +0200") {
    t.Errorf("Expected emailInfo.Date to return 'Date' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

func TestParseFileWithMissingFrom(t *testing.T){
  r1 := strings.NewReader(missingFrom)
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "||") {
    t.Errorf("Expected emailInfo.From to return 'From' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != "Cuit Vapeur 29.90 euros, Nintendo 3DS 239 euros, GPS TOM TOM 139 euros... decouvrez VITE tous les bons plans du weekend !") {
    t.Errorf("Expected emailInfo.Subject to return 'Subject' but got %s instead", emailInfo.Subject)
  }
  if (emailInfo.Date != "01 Apr 2011 16:17:41 +0200") {
    t.Errorf("Expected emailInfo.Date to return 'Date' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

func TestParseFileWithMissingAll(t *testing.T){
  r1 := strings.NewReader(missingAll)
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "||") {
    t.Errorf("Expected emailInfo.From to return 'From' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != "||") {
    t.Errorf("Expected emailInfo.Subject to return 'Subject' but got %s instead", emailInfo.Subject)
  }
  if (emailInfo.Date != "||") {
    t.Errorf("Expected emailInfo.Date to return 'Date' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

func TestEmailParsing(t *testing.T){
  dir, err := filepath.Abs(filepath.Dir("."))
  if err != nil {
      log.Fatal(err)
  }

  ch := make(chan EmailInformation)
  files := []string{dir + "/test.msg",dir + "/test.msg",dir + "/test.msg",dir + "/test.msg",dir + "/test.msg",dir + "/test.msg"}
  for _, file := range files {
    wg.Add(1)
    go emailParsing(file, ch)
  }

  go func(wg sync.WaitGroup, ch chan EmailInformation) {
		log.Println("waiting")
		wg.Wait()
		log.Println("done waiting")
		close(ch)
	}(wg, ch)

  for i := 0; i < len(files); i++ {
    s := <-ch
    //log.Printf("{from: %s, subject: %s, date: %s, filename: %s }\n", s.From, s.Subject, s.Date, s.FileName)
    if (s.From != "\"Darty\" <infos@contact-darty.com>") {
      t.Errorf("Expected emailInfo.From to return 'From' but got %s instead", s.From)
    }
    if (s.Subject != "Hello World!") {
      t.Errorf("Expected emailInfo.Subject to return 'Hello World!' but got %s instead", s.Subject)
    }
    if (s.Date != "01 Apr 2011 16:17:41 +0200") {
      t.Errorf("Expected emailInfo.Date to return 'Date' but got %s instead", s.Date)
    }
    if (s.FileName != "test.msg") {
      t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", s.FileName)
    }
  }
}

func TestReadAndParseFiles(t *testing.T){
  dir, err := filepath.Abs(filepath.Dir("."))
  if err != nil {
      log.Fatal(err)
  }
  readAndParseFiles(dir)
}
