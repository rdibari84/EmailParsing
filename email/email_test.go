package main

import (
  "testing"
  "strings"
  "bufio"
  "log"
  "path/filepath"
  "os"
  "io"
  "bytes"
)

////////////////////////////////////
//////// extractEmail Tests ////////
////////////////////////////////////

func TestExtractEmailUnwantedChar(t *testing.T){
  email := "<Amway@MagazineLine.com>"
  e := extractEmail(email)
  if e != "Amway@MagazineLine.com" {
    t.Errorf("Expected extractEmail to return 'Amway@MagazineLine.com' but got %s instead", e)
  }
}


func TestExtractEmailMultiUnwantedChar(t *testing.T){
  email := "=?utf-8?q?Aviva?= <aviva@avivaemail.co.uk>"
  e := extractEmail(email)
  if e != "aviva@avivaemail.co.uk" {
    t.Errorf("Expected extractEmail to return 'aviva@avivaemail.co.uk' but got %s instead", e)
  }
}

func TestExtractEmailName(t *testing.T){
  email := " \"AwaySomeTime\" <Amway@MagazineLine.com>"
  e := extractEmail(email)
  if e != "Amway@MagazineLine.com" {
    t.Errorf("Expected extractEmail to return 'Amway@MagazineLine.com' but got %s instead", e)
  }
}

func TestVaildEmailDot(t *testing.T) {
  email := "customer.service@citicorp.com"
  e := extractEmail(email)
  if e != "customer.service@citicorp.com" {
    t.Errorf("Expected extractEmail to return 'customer.service@citicorp.com' but got %s instead", e)
  }
}

func TestVaildEmailMultiDotsAndNet(t *testing.T) {
  email := "alchemyworx@cp.assurance.returnpath.net"
  e := extractEmail(email)
  if e != "alchemyworx@cp.assurance.returnpath.net" {
    t.Errorf("Expected extractEmail to return 'alchemyworx@cp.assurance.returnpath.net' but got %s instead", e)
  }
}

func TestVaildEmailUnderscore(t *testing.T) {
  email := "customer_service@citicorp.com"
  e := extractEmail(email)
  if e != "customer_service@citicorp.com" {
    t.Errorf("Expected extractEmail to return 'customer_service@citicorp.com' but got %s instead", e)
  }
}

func TestVaildEmailNumbers(t *testing.T) {
  email := "3114.352_22@45280.com"
  e := extractEmail(email)
  if e != "3114.352_22@45280.com" {
    t.Errorf("Expected extractEmail to return '3114.352_22@45280.com' but got %s instead", e)
  }
}

func TestBadEmail(t *testing.T) {
  email := "citicorp.com"
  e := extractEmail(email)
  if e != "||" {
    t.Errorf("Expected extractEmail to return '||' but got %s instead", e)
  }
}

////////////////////////////////////
//////// parseFile Tests ///////////
////////////////////////////////////

var badEmail = `
Return-Path: <bo-bx6rw3raupta74au6m0rxbysph09qe@b.avivaemail.co.uk>
X-Original-To: alchemyworx@cp.assurance.returnpath.net
Delivered-To: assurance@localhost.returnpath.net
Received: from mxa-d1.returnpath.net (unknown [10.8.2.117])
	by cpa-d1.returnpath.net (Postfix) with ESMTP id 84EBF198271
	for <alchemyworx@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:14:28 -0600 (MDT)
Received: from mta141.norwichunionemail.com (mta141.norwichunionemail.com [64.37.95.217])
	by mxa-d1.returnpath.net (Postfix) with SMTP id 727BEBCC
	for <alchemyworx@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:14:27 -0600 (MDT)
DomainKey-Signature: a=rsa-sha1; q=dns; c=nofws;
  s=200505; d=avivaemail.co.uk;
  b=E4WJDZHtNvKzks4LOkpPErtbJZwW7nlYg6X21Fu5H66AXXaHgKxs7bVXeE7GQ7Sx4jB4T7fv4CdYp/JW4mIrJpuubCCc8ZzzvPo6PEdJHNT856KIR1Dt/UR+Jx+fcbYQ/NO/jIilt4KQfeNzsm+FvP7QnV6jRXeOq1k+CTc1xHo=;
 h=Date:Message-ID:List-Unsubscribe:From:To:Subject:MIME-Version:Reply-To:Content-type;
Date: Fri, 1 Apr 2011 14:14:49 -0000
Message-ID: <bx6rw3raupta74au6m0rxbysph09qe.0.15@mta141.avivaemail.co.uk>
List-Unsubscribe: <mailto:rm-0bx6rw3raupta74au6m0rxbysph09qe@avivaemail.co.uk>
From: =?utf-8?q?Aviva?= <aviva@@avivaemail.co.uk>
To: alchemyworx@cp.assurance.returnpath.net
Subject: (TEST-Multipart)
MIME-Version: 1.0
Reply-To: =?utf-8?q?Aviva?= <support-bx6rw3raupta74au6m0rxbysph09qe@avivaemail.co.uk>
Content-type: multipart/alternative; boundary="=bx6rw3raupta74au6m0rxbysph09qe"`

func TestParseFileWithBadEmail(t *testing.T){
  r1 := strings.NewReader(badEmail)
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "||") {
    t.Errorf("Expected emailInfo.From to return '||' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != "(TEST-Multipart)") {
    t.Errorf("Expected (TEST-Multipart). Got %s instead", emailInfo.Subject)
  }
  if (emailInfo.Date != "Fri, 1 Apr 2011 14:14:49 -0000") {
    t.Errorf("Expected emailInfo.Date to return 'Date' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

var mulitSubject = `Return-Path: <bo-bx6rw3raupta74au6m0rxbysph09qe@b.avivaemail.co.uk>
X-Original-To: alchemyworx@cp.assurance.returnpath.net
Delivered-To: assurance@localhost.returnpath.net
Received: from mxa-d1.returnpath.net (unknown [10.8.2.117])
	by cpa-d1.returnpath.net (Postfix) with ESMTP id 84EBF198271
	for <alchemyworx@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:14:28 -0600 (MDT)
Received: from mta141.norwichunionemail.com (mta141.norwichunionemail.com [64.37.95.217])
	by mxa-d1.returnpath.net (Postfix) with SMTP id 727BEBCC
	for <alchemyworx@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:14:27 -0600 (MDT)
DomainKey-Signature: a=rsa-sha1; q=dns; c=nofws;
  s=200505; d=avivaemail.co.uk;
  b=E4WJDZHtNvKzks4LOkpPErtbJZwW7nlYg6X21Fu5H66AXXaHgKxs7bVXeE7GQ7Sx4jB4T7fv4CdYp/JW4mIrJpuubCCc8ZzzvPo6PEdJHNT856KIR1Dt/UR+Jx+fcbYQ/NO/jIilt4KQfeNzsm+FvP7QnV6jRXeOq1k+CTc1xHo=;
 h=Date:Message-ID:List-Unsubscribe:From:To:Subject:MIME-Version:Reply-To:Content-type;
Date: Fri, 1 Apr 2011 14:14:49 -0000
Message-ID: <bx6rw3raupta74au6m0rxbysph09qe.0.15@mta141.avivaemail.co.uk>
List-Unsubscribe: <mailto:rm-0bx6rw3raupta74au6m0rxbysph09qe@avivaemail.co.uk>
From: "LA Galaxy" <enews@events.lagalaxy.com>
To: alchemyworx@cp.assurance.returnpath.net
Subject: (TEST-Multipart) =?utf-8?q?=5BRetention_In_Life_ezine=5Fhome=5F050411=5D_Introducing_Your_?=
 =?utf-8?q?Aviva_Essentials=3A_Win_4_tickets_to_the_Aviva_Premiership_Rugb?=
 =?utf-8?q?y_Final=2C_Keep_the_cost_of_driving_down_and_more?=
MIME-Version: 1.0
Reply-To: =?utf-8?q?Aviva?= <support-bx6rw3raupta74au6m0rxbysph09qe@avivaemail.co.uk>
Content-type: multipart/alternative; boundary="=bx6rw3raupta74au6m0rxbysph09qe"`

func TestParseFileWithMultiLineSubject(t *testing.T){
  r1 := strings.NewReader(mulitSubject)
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "enews@events.lagalaxy.com") {
    t.Errorf("Expected emailInfo.From to return 'enews@events.lagalaxy.com' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != `(TEST-Multipart) =?utf-8?q?=5BRetention_In_Life_ezine=5Fhome=5F050411=5D_Introducing_Your_?= =?utf-8?q?Aviva_Essentials=3A_Win_4_tickets_to_the_Aviva_Premiership_Rugb?= =?utf-8?q?y_Final=2C_Keep_the_cost_of_driving_down_and_more?=`) {
    t.Errorf("Expected %v. Got %v ", `(TEST-Multipart) =?utf-8?q?=5BRetention_In_Life_ezine=5Fhome=5F050411=5D_Introducing_Your_?=
     =?utf-8?q?Aviva_Essentials=3A_Win_4_tickets_to_the_Aviva_Premiership_Rugb?=
     =?utf-8?q?y_Final=2C_Keep_the_cost_of_driving_down_and_more?=`, emailInfo.Subject)
  }
  if (emailInfo.Date != "Fri, 1 Apr 2011 14:14:49 -0000") {
    t.Errorf("Expected emailInfo.Date to return 'Date' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

var validMsg = `Return-Path: <bounce-17011457_HTML-295436918-26481642-10016540-37@bounce.sls.enewsaxs.com>
X-Original-To: aeg@cp.delivery.ncrcustomerpower.com
Delivered-To: assurance@localhost.returnpath.net
Received: from mxa-d1.returnpath.net (unknown [10.8.2.117])
	by cpa-d1.returnpath.net (Postfix) with ESMTP id 72714198272
	for <aeg@cp.delivery.ncrcustomerpower.com>; Fri,  1 Apr 2011 12:06:36 -0600 (MDT)
Received: from mta.sls.enewsaxs.com (mta.sls.enewsaxs.com [64.132.89.25])
	by mxa-d1.returnpath.net (Postfix) with ESMTP id 014F1F99
	for <aeg@cp.delivery.ncrcustomerpower.com>; Fri,  1 Apr 2011 12:06:36 -0600 (MDT)
DKIM-Signature: v=1; a=rsa-sha1; c=relaxed/relaxed; s=200608; d=events.lagalaxy.com;
 h=From:To:Subject:Date:List-Unsubscribe:MIME-Version:Reply-To:Message-ID:Content-Type:Content-Transfer-Encoding; i=enews@events.lagalaxy.com;
 bh=FKdYuRbX4MLx8jwGZMeClvG4XuI=;
 b=J3KDT/5dKurpogHe/xXkDrR8tJpjNaStAd3EaVJnnuc49A0NZG9M+meFSJpqwCV/J4blSrHV8Wzq
   E4PrAxOKXHxx0vqIo1OxUY4Z+3Adb5r3BYU78pC+8pwKGIW7IQVxkRRjMlccuTWDe04pxPKLTc/w
   nHYhWZZBrNPezfpBZEI=
DomainKey-Signature: a=rsa-sha1; c=nofws; q=dns; s=200608; d=events.lagalaxy.com;
 b=tzuRFzrkMXsBp9VRhxb23h2oRZLmnk6ITXnm7F3UxGRWy82WWojdhJ+VkL/b9p8S+voO/tjX7dNK
   J5y+5Nvx8onr6x5IaHy5TtSoWXtJQpD/BGOZbmFXlcJ/ovIWj6bYjF/Fv6ZDOCxUwJ5e1446r6jB
   MmzPff1ga34Kb11AOvQ=;
Received: by mta.sls.enewsaxs.com (PowerMTA(TM) v3.5r15) id hioa2o0ie1sa for <aeg@cp.delivery.ncrcustomerpower.com>; Fri, 1 Apr 2011 12:06:25 -0600 (envelope-from <bounce-17011457_HTML-295436918-26481642-10016540-37@bounce.sls.enewsaxs.com>)
From: "LA Galaxy" <enews@events.lagalaxy.com>
To: <aeg@cp.delivery.ncrcustomerpower.com>
Subject: April Fool's Day Offer, Save up to 40% with no fees
Date: Fri, 01 Apr 2011 12:06:22 -0600
List-Unsubscribe: <mailto:leave-fc5115777565057970131f2d205921-fe2016797060077b7c1578-fe5e1676716c057b7116-fe6e1570756502787114-ff951777@leave.sls.enewsaxs.com>
MIME-Version: 1.0
Reply-To: "LA Galaxy" <reply-fe5e1676716c057b7116-17011457_HTML-295436918-10016540-37@events.lagalaxy.com>
x-job: 10016540_26481642
Message-ID: <3709e1a3-663f-464c-a38f-584ae8c9fe24@xtinmta105.xt.local>
Content-Type: text/html;
	charset="us-ascii"
Content-Transfer-Encoding: 7bit`

func TestParseFileWithValidMsg(t *testing.T){
  r1 := strings.NewReader(validMsg)
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "enews@events.lagalaxy.com") {
    t.Errorf("Expected emailInfo.From to return 'enews@events.lagalaxy.com' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != "April Fool's Day Offer, Save up to 40% with no fees") {
    t.Errorf("Got %s", emailInfo.Subject)
  }
  if (emailInfo.Date != "Fri, 01 Apr 2011 12:06:22 -0600") {
    t.Errorf("Expected emailInfo.Date to return 'Fri, 01 Apr 2011 12:06:22 -0600' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

var missingDate = `Return-Path: <infos@contact-darty.com>
X-Original-To: 1000mercis@cp.assurance.returnpath.net
Delivered-To: assurance@localhost.returnpath.net
Received: from mxa-d1.returnpath.net (unknown [10.8.2.117])
	by cpa-d1.returnpath.net (Postfix) with ESMTP id 426E5198271
	for <1000mercis@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:17:45 -0600 (MDT)
Received: from smtp-8-ft1.mm.fr.colt.net (smtp-7-ft1.mm.fr.colt.net [62.23.8.162])
	by mxa-d1.returnpath.net (Postfix) with ESMTP id 2906A1CD
	for <1000mercis@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 08:17:44 -0600 (MDT)
Received: from host.25.62.23.62.rev.coltfrance.com ([62.23.62.25]:62162 helo=contact-darty.com)
	by massmail-ft1.infra.coltfrance.com with esmtp (Exim)
	id 1Q5fAU-00030S-4i
	for <1000mercis@cp.assurance.returnpath.net>; Fri, 01 Apr 2011 16:17:42 +0200
From: "Darty" <infos@contact-darty.com>
To: 1000mercis@cp.assurance.returnpath.net
Subject: Cuit Vapeur 29.90 euros
Message-ID: <20110401161739.E3786358A9D7B977@contact-darty.com>
MIME-Version: 1.0
x-idmail: DartyCRM_322_385774_10000
Content-Type: text/html;
	charset="iso-8859-1"
Content-Transfer-Encoding: 7bit`

func TestParseFileWithMissingDate(t *testing.T){
  r1 := strings.NewReader(missingDate)
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "infos@contact-darty.com") {
    t.Errorf("Expected emailInfo.From to return 'infos@contact-darty.com' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != "Cuit Vapeur 29.90 euros") {
    t.Errorf("Expected emailInfo.Subject to return 'Cuit Vapeur 29.90 euros' but got %s instead", emailInfo.Subject)
  }
  if (emailInfo.Date != "||") {
    t.Errorf("Expected emailInfo.Date to return '||' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

var missingFrom = `Return-Path: <bounce-17011457_HTML-295436918-26481642-10016540-37@bounce.sls.enewsaxs.com>
X-Original-To: aeg@cp.delivery.ncrcustomerpower.com
Delivered-To: assurance@localhost.returnpath.net
Received: from mxa-d1.returnpath.net (unknown [10.8.2.117])
	by cpa-d1.returnpath.net (Postfix) with ESMTP id 72714198272
	for <aeg@cp.delivery.ncrcustomerpower.com>; Fri,  1 Apr 2011 12:06:36 -0600 (MDT)
Received: from mta.sls.enewsaxs.com (mta.sls.enewsaxs.com [64.132.89.25])
	by mxa-d1.returnpath.net (Postfix) with ESMTP id 014F1F99
	for <aeg@cp.delivery.ncrcustomerpower.com>; Fri,  1 Apr 2011 12:06:36 -0600 (MDT)
DKIM-Signature: v=1; a=rsa-sha1; c=relaxed/relaxed; s=200608; d=events.lagalaxy.com;
 h=From:To:Subject:Date:List-Unsubscribe:MIME-Version:Reply-To:Message-ID:Content-Type:Content-Transfer-Encoding; i=enews@events.lagalaxy.com;
 bh=FKdYuRbX4MLx8jwGZMeClvG4XuI=;
 b=J3KDT/5dKurpogHe/xXkDrR8tJpjNaStAd3EaVJnnuc49A0NZG9M+meFSJpqwCV/J4blSrHV8Wzq
   E4PrAxOKXHxx0vqIo1OxUY4Z+3Adb5r3BYU78pC+8pwKGIW7IQVxkRRjMlccuTWDe04pxPKLTc/w
   nHYhWZZBrNPezfpBZEI=
DomainKey-Signature: a=rsa-sha1; c=nofws; q=dns; s=200608; d=events.lagalaxy.com;
 b=tzuRFzrkMXsBp9VRhxb23h2oRZLmnk6ITXnm7F3UxGRWy82WWojdhJ+VkL/b9p8S+voO/tjX7dNK
   J5y+5Nvx8onr6x5IaHy5TtSoWXtJQpD/BGOZbmFXlcJ/ovIWj6bYjF/Fv6ZDOCxUwJ5e1446r6jB
   MmzPff1ga34Kb11AOvQ=;
Received: by mta.sls.enewsaxs.com (PowerMTA(TM) v3.5r15) id hioa2o0ie1sa for <aeg@cp.delivery.ncrcustomerpower.com>; Fri, 1 Apr 2011 12:06:25 -0600 (envelope-from <bounce-17011457_HTML-295436918-26481642-10016540-37@bounce.sls.enewsaxs.com>)
To: "AEGCorp" <aeg@cp.delivery.ncrcustomerpower.com>
Subject: subject!
Date: 01 Apr 2011 16:17:41 +0200
List-Unsubscribe: <mailto:leave-fc5115777565057970131f2d205921-fe2016797060077b7c1578-fe5e1676716c057b7116-fe6e1570756502787114-ff951777@leave.sls.enewsaxs.com>
MIME-Version: 1.0
Reply-To: "LA Galaxy" <reply-fe5e1676716c057b7116-17011457_HTML-295436918-10016540-37@events.lagalaxy.com>
x-job: 10016540_26481642
Message-ID: <3709e1a3-663f-464c-a38f-584ae8c9fe24@xtinmta105.xt.local>
Content-Type: text/html;
	charset="us-ascii"
Content-Transfer-Encoding: 7bit`

func TestParseFileWithMissingFrom(t *testing.T){
  r1 := strings.NewReader(missingFrom)
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "||") {
    t.Errorf("Expected emailInfo.From to return 'From' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != "subject!") {
    t.Errorf("Expected emailInfo.Subject to return 'subject!' but got %s instead", emailInfo.Subject)
  }
  if (emailInfo.Date != "01 Apr 2011 16:17:41 +0200") {
    t.Errorf("Expected emailInfo.Date to return 'Date' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

var missingAll = `Return-Path: <out-582632-B2C71BD37AF148CE9D728B61264F854D@mail.beliefnet.com>
X-Original-To: beliefnet@cp.monitor1.returnpath.net
Delivered-To: assurance@localhost.returnpath.net
Received: from mxa-d1.returnpath.net (unknown [10.8.2.117])
	by cpa-d1.returnpath.net (Postfix) with ESMTP id 4F39C19825C
	for <beliefnet@cp.monitor1.returnpath.net>; Fri,  1 Apr 2011 08:12:01 -0600 (MDT)
Received: from mail1101.mail.beliefnet.com (mail1101.mail.beliefnet.com [64.27.69.101])
	by mxa-d1.returnpath.net (Postfix) with ESMTP id 6275D992
	for <beliefnet@cp.monitor1.returnpath.net>; Fri,  1 Apr 2011 08:12:00 -0600 (MDT)
DKIM-Signature: v=1; a=rsa-sha1; c=relaxed/relaxed; s=default; d=mail.beliefnet.com;
 h=From:Reply-To:To:Message-ID:Subject:MIME-Version:Content-Type:Content-Transfer-Encoding; i=specialoffers@mail.beliefnet.com;
 bh=kwfm9qePWfYmPSIHlEjFYVCYzlg=;
 b=SNIGrGXvBRXp4hSf72CttBCV3dmQd3F1UphY00hAKhNFiOhPCUtwTknlclAzlkYJ2OBGBXpoNxlA
   EIF4Uu4o5Q==
DomainKey-Signature: a=rsa-sha1; c=nofws; q=dns; s=default; d=mail.beliefnet.com;
 b=RvVmJOpSKiPqhWQl7saVvXjcbR7D4c9cUR6T0wQs7AlVzuobcxddkQ0GiFKIYdQNQHUOsh5bBZ/C
   0IuUBJPPxg==;
Received: by mail1101.mail.beliefnet.com (PowerMTA(TM) v3.5r13) id hinej00sgh8s for <beliefnet@cp.monitor1.returnpath.net>; Fri, 1 Apr 2011 10:11:27 -0400 (envelope-from <out-582632-B2C71BD37AF148CE9D728B61264F854D@mail.beliefnet.com>)
Reply-To: Announce - Beliefnet Sponsor <r-ethtsdsrhbqfybrtqfjswzrjlyzvkrtpwcqlzspjhrshtkpbtl@mail.beliefnet.com>
To:  <beliefnet@cp.monitor1.returnpath.net>
Message-ID: <527817310.344.1301667087687.JavaMail.root@mail.beliefnet.com>
MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: quoted-printable
X-ID: slvvmprjvrkvtbrbtbjvbckljmk wrhsbwfvszdbclhgpfsrlwkbrkcdlmcftcrjrbkmptvmbcp bvpjksbvbp
X-RPTags: Seed
X-NLCID: 35149`

func TestParseFileWithMissingAll(t *testing.T){
  r1 := strings.NewReader(missingAll)
  r := bufio.NewReader(r1)
  emailInfo := parseFile("file.msg", r)
  if (emailInfo.From != "||") {
    t.Errorf("Expected emailInfo.From to return '||' but got %s instead", emailInfo.From)
  }
  if (emailInfo.Subject != "||") {
    t.Errorf("Expected emailInfo.Subject to return '||' but got %s instead", emailInfo.Subject)
  }
  if (emailInfo.Date != "||") {
    t.Errorf("Expected emailInfo.Date to return '||' but got %s instead", emailInfo.Date)
  }
  if (emailInfo.FileName != "file.msg") {
    t.Errorf("Expected emailInfo.FileName to return 'file.msg' but got %s instead", emailInfo.FileName)
  }
}

////////////////////////////////////
//////// Integration Tests ////////
////////////////////////////////////
var validMsg1 = `Return-Path: <suncoast@boydgaming.net>
X-Original-To: boydgamingcorporation@cp.assurance.returnpath.net
Delivered-To: assurance@localhost.returnpath.net
Received: from mxa-d1.returnpath.net (unknown [10.8.2.117])
	by cpa-d1.returnpath.net (Postfix) with ESMTP id 0937F198271
	for <boydgamingcorporation@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 11:36:27 -0600 (MDT)
Received: from smtp2.boydgaming.net (smtp2.boydgaming.net [64.79.129.191])
	by mxa-d1.returnpath.net (Postfix) with ESMTP id C120D126
	for <boydgamingcorporation@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 11:36:26 -0600 (MDT)
Received: from [10.251.20.13] (unknown [10.251.20.13])
	by smtp2.boydgaming.net (Postfix) with ESMTP id 3093227844F
	for <boydgamingcorporation@cp.assurance.returnpath.net>; Fri,  1 Apr 2011 10:36:26 -0700 (PDT)
X-DKIM: Sendmail DKIM Filter v2.8.3 smtp2.boydgaming.net 3093227844F
DKIM-Signature: v=1; a=rsa-sha1; c=simple/simple; d=boydgaming.net; s=smtp;
	t=1301679386; bh=DyM8CohSbjCN3K1AByz8K3FHMsA=;
	h=To:From:Reply-To:Subject:Date:MIME-Version:Content-Type:
	 Content-Transfer-Encoding:List-Unsubscribe:Message-ID;
	b=XtHwSYl7M0uQt3U8XI/0f18VYVe/ta7nARZB1GN4TWbczqD8jIOf7cLIp93IX/Q27
	 flggpcloqkcxlrCmJ4qWYstPrVGhpHfCdVJdTQMf9CmnTjjWfqsJ1mPL6QLOdcKxRB
	 XmkjNNLtheRsBRLE1/e6wGHNutpri6zNCZzx7xgA=
To: Lisa Marshall <boydgamingcorporation@cp.assurance.returnpath.net>
From: Suncoast Hotel & Casino - Las Vegas <suncoast@boydgaming.net>
Reply-To: Suncoast Hotel & Casino - Las Vegas <suncoast@boydgaming.net>
Subject: See What's Happening with our Table Games!
Date: Fri, 01 Apr 2011 10:36:26 -0700
X-LibVersion: 3.3.2
MIME-Version: 1.0
Content-Type: multipart/alternative;
 boundary="_=_swift-13596511954d960d1a312b33.37345773_=_"
Content-Transfer-Encoding: 7bit
X-Ninja-Mailer-ID: 6752
List-Unsubscribe: <http://www.boydgaming.com/unsubscribe?6752XXLUX0X0XTEST>
Message-ID: <20110401173626.15575.2089030531.swift@webadmin.boydgaming.net>`

func WriteTestFile(name string) string{
  dir, err := filepath.Abs(filepath.Dir("."))
  if err != nil {
      log.Fatal(err)
  }

  outfile := dir + "/" + name

  // write file
  f, err := os.Create(outfile)
  checkError(err, "Issue creating " + outfile)
  defer f.Close()
  w := bufio.NewWriter(f)

  w.Write([]byte(validMsg1))
  w.Flush()
  log.Printf("Succesfully wrote %s", outfile)

  return outfile
}

/*
Use bytes.count because its faster than scanning.
takes awaysextra logic and buffering required to return whole lines
and takes advantage of some assembly optimized functions offered by the bytes package to search characters in a byte slice.
32k buffer should be good for large files
*/
func lineCounter(r io.Reader) (int, error) {
    buf := make([]byte, 32*1024)
    count := 0
    lineSep := []byte{'\n'}

    for {
        c, err := r.Read(buf)
        count += bytes.Count(buf[:c], lineSep)

        switch {
        case err == io.EOF:
            return count, nil

        case err != nil:
            return count, err
        }
    }
}


func TestEmailParsing(t *testing.T){
  outfile := WriteTestFile("test.msg")

  // Make channel
  ch := make(chan EmailInformation)
  files := []string{outfile,outfile,outfile,outfile,outfile}
  for _, file := range files {
    wgw.Add(1)
    go emailParsing(file, ch) // call test function
  }

  go func(wgw WaitGroupWrapper, ch chan EmailInformation) {
		log.Println("waiting")
		wgw.Wg.Wait()
		log.Println("done waiting")
		close(ch)
	}(wgw, ch)

  for i := 0; i < len(files); i++ {
    s := <-ch
    if (s.From != "suncoast@boydgaming.net") {
      t.Errorf("Expected emailInfo.From to return 'suncoast@boydgaming.net' but got %s instead", s.From)
    }
    if (s.Subject != "See What's Happening with our Table Games!") {
      t.Errorf("Expected emailInfo.Subject to return 'See What's Happening with our Table Games!' but got %s instead", s.Subject)
    }
    if (s.Date != "Fri, 01 Apr 2011 10:36:26 -0700") {
      t.Errorf("Expected emailInfo.Date to return 'Fri, 01 Apr 2011 10:36:26 -0700' but got %s instead", s.Date)
    }
    if (s.FileName != "test.msg") {
      t.Errorf("Expected emailInfo.FileName to return 'test.msg' but got %s instead", s.FileName)
    }
  }

  os.Remove(outfile)
}

func TestReadParseAndWriteFiles(t *testing.T){
  outfile := WriteTestFile("test.msg")
  outfile1 := WriteTestFile("test1.msg")
  outfile2 := WriteTestFile("test2.msg")

  dir, err := filepath.Abs(filepath.Dir("."))
  if err != nil {
      log.Fatal(err)
  }
  dir = validateInput(dir, dir + "test.txt")
  readParseAndWriteFiles(dir, dir + "test.txt")

  // assert that a file was actually written
  file, err := os.Open(dir + "test.txt")
  if err != nil {
      t.Errorf("Did not expect an error opening %s, but got one. error: %v", dir + "test.txt", err)
  }
  lineNum, err := lineCounter(file)
  if err != nil {
      t.Errorf("Did not expect an error reading the number of lines in %s but got one. error: %v", dir + "test.txt", err)
  }
  if lineNum != 3 {
    t.Errorf("Expected three lines")
  }
  file.Close()

  os.Remove(outfile)
  os.Remove(outfile1)
  os.Remove(outfile2)
  os.Remove(dir + "test.txt")
}
