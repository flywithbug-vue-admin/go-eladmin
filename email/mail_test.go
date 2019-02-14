package email

import (
	"testing"

	"gopkg.in/gomail.v2"
)

func TestMailVerifyailTest(t *testing.T) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "flywithbug@163.com", "Óscar García Amor")
	m.SetHeader("To", "369495368@qq.com")
	m.SetHeader("Subject", "test")
	m.SetBody("text/html", "Hello <b>Alice</b> and <i>Bob</i>!")

	d := gomail.NewDialer("smtp.163.com", 465, "flywithbug@163.com", "wyflywithbug1112")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
