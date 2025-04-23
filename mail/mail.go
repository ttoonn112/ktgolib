package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
  	"strings"
	//"github.com/kr/pretty"
)

var Mail_SystemName string
var Mail_Server string
var Mail_Port string
var Mail_User string
var Mail_Pass string
var Mail_Admin []string

type Request struct {
	from    string
	to      []string
  	cc      []string
  	bcc     []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func Setup(system_name string,server string, port string, user string, pass string, mail_admin []string){
	Mail_SystemName = system_name
  	Mail_Server = server
  	Mail_Port = port
  	Mail_User = user
  	Mail_Pass = pass
	Mail_Admin = mail_admin
}

func NewRequest(to []string,cc []string,bcc []string, subject string) *Request {
	return &Request{
		to: to,
    	cc: cc,
    	bcc: bcc,
		subject: subject,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() error {
	var to []string = nil
	body := "From: "+Mail_SystemName+"<"+Mail_User+">\r\n"
	if r.to != nil && len(r.to) > 0 && strings.Join(r.to, ",") != "" {
		body += "To: " + strings.Join(r.to, ",") + "\r\n"
		to = r.to
	}
  	if r.cc != nil && len(r.cc) > 0 && strings.Join(r.cc, ",") != "" {
		if r.to == nil || len(r.to) == 0 || strings.Join(r.to, ",") == "" {
			body += "To: "+strings.Join(r.cc, ",")+"\r\n"
			to = r.cc
		}else{
			body += "Cc: " + strings.Join(r.cc, ",") + "\r\n"
			to = append(to, r.cc...)
		}
  	}
	if to == nil {
		return false
	}
  	if r.bcc != nil && len(r.bcc) > 0 {
   	 	body += "Bcc: " + strings.Join(r.bcc, ",") + "\r\n"
		to = append(to, r.bcc...)
  	}
  	body += "Subject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%s", Mail_Server, Mail_Port)

	return smtp.SendMail(SMTP, smtp.PlainAuth("", Mail_User, Mail_Pass, Mail_Server), Mail_User, to, []byte(body))
}

func (r *Request) sendMailWithTemplate(templateName string, items interface{}) error {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		return err 
	}
	return r.sendMail()
}

func Send(template string, rcpt []string, cc []string, subject string, message map[string]string) error {
  	r := NewRequest(rcpt, cc, Mail_Admin, subject)
  	message["subject"] = subject
	return r.sendMailWithTemplate(template, message)
}