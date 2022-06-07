package report

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/xhit/go-simple-mail/v2"
)

func SendMail(subject string, body string, to ...string) {
	server := mail.NewSMTPClient()
	setConfig(server)

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}
	email := mail.NewMSG()
	// 获取留作业的教师的邮箱

	email.SetFrom(viper.GetString("mail.from")).
		AddTo(to...).
		// AddCc() 暂时不需要抄送，可在需要留存电子邮件档案时对其进行设置
		SetSubject(subject)

	email.SetBody(mail.TextHTML, body)

	if email.Error != nil {
		log.Fatal(email.Error)
	}

	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent")
	}
}

func setConfig(server *mail.SMTPServer) {
	server.Host = viper.GetString("mail.host")
	server.Port = viper.GetInt("mail.port")
	server.Username = viper.GetString("mail.username")
	server.Password = viper.GetString("mail.password")
	// server.Authentication=mail.AuthPlain
	// 单次发送，不必保活
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	// server.TLSConfig = nil
}
