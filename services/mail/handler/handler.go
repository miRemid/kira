package handler

import (
	"context"
	"log"

	"github.com/go-gomail/gomail"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/proto/pb"
)

var (
	account  string
	password string
	client   *gomail.Dialer
)

func init() {
	account = common.Getenv("MAIL_ACCOUNT", "")
	password = common.Getenv("MAIL_PASSWORD", "")
	log.Println(account, password)
	client = gomail.NewDialer("smtp.office365.com", 587, account, password)
}

func SendMail(ctx context.Context, in *pb.SendMailReq) error {
	// send mail
	m := gomail.NewMessage()
	m.SetHeader("From", account)
	m.SetHeader("To", in.To)
	m.SetHeader("Subject", in.Subject)
	m.SetBody("text/html", in.Content)
	return client.DialAndSend(m)
}
