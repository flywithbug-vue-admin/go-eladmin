package email

import (
	"errors"
	"fmt"
	"go-eladmin/config"
	"regexp"
	"strings"

	"gopkg.in/gomail.v2"
)

var (
	Mail *gomail.Dialer
)

func sendMail(to, title, subject, body, from string) error {
	if !MailVerify(to) {
		return fmt.Errorf("mail not right")
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, title)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	//log4go.Info(" %s %s ", to, m.GetHeader("To"))
	if Mail == nil {
		return errors.New("mail is nil")
	}
	return Mail.DialAndSend(m)
}

func SendMail(title, subject, content, mail string) error {
	return sendMail(mail, title, subject, content, config.Conf().MailConfig.Username)
}

var routerRe = regexp.MustCompile(`^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`)

func MailVerify(mail string) bool {
	match := routerRe.FindString(mail)
	return strings.EqualFold(match, mail)
}

func SendVerifyCode(_, code, mail string) error {
	body := fmt.Sprintf("您的验证码是： %s ", code)
	return sendMail(mail, "FlyWithBug", "邮箱验证", body, config.Conf().MailConfig.Username)
}

//func ReDialer(host string, port int, username, password string) (*gomail.Dialer, error) {
//	Mail = gomail.NewDialer(host, port, username, password)
//	return Mail, nil
//}
