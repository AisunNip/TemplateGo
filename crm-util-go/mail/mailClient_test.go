package mail

import (
	"fmt"
	"testing"
)

/*
	gmail: Less secure app access (ON)
	Port: 587  TLS   Transport Layer Security (TLS) has replaced SSL
	Port: 465  SSL
*/

func TestSendEmailSSL(t *testing.T) {
	mc := new(MailClient)
	mc.SmtpServer = "smtp.gmail.com"
	mc.SmtpPortNo = "465"
	mc.Username = "paravit.tun@gmail.com"
	mc.Password = "xxx"
	mc.CryptoProtocol = "SSL"

	mailMessage := mc.NewMessage("Test subject SSL", "Test body SSL")
	mailMessage.To = []string{"paravit_tun@truecorp.co.th"}
	err := mailMessage.AttachFile("D:/PB-210308-071_resolve.docx")

	if err != nil {
		fmt.Println("AttachFile fail", err.Error())
		return
	}

	err = mc.SendEmail(mailMessage)

	if err != nil {
		fmt.Println("SendEmail fail", err.Error())
		t.Errorf("TestSendEmailSSL Error %v", err.Error())
	} else {
		fmt.Println("TestSendEmailSSL success")
	}
}

func TestSendEmailTLS(t *testing.T) {
	mc := new(MailClient)
	mc.SmtpServer = "smtp.gmail.com"
	mc.SmtpPortNo = "587"
	mc.Username = "paravit.tun@gmail.com"
	mc.Password = "xxx"
	mc.CryptoProtocol = "TLS"

	mailMessage := mc.NewMessage("Test subject TLS", "Test body TLS")
	mailMessage.To = []string{"paravit_tun@truecorp.co.th"}
	err := mailMessage.AttachFile("D:/PB-210308-071_resolve.docx")

	if err != nil {
		fmt.Println("AttachFile fail", err.Error())
		return
	}

	err = mc.SendEmail(mailMessage)

	if err != nil {
		fmt.Println("SendEmail fail", err.Error())
		t.Errorf("TestSendEmailTLS Error %v", err.Error())
	} else {
		fmt.Println("TestSendEmailTLS success")
	}
}

func TestSendEmailTrueCorpTLS(t *testing.T) {
	mc := new(MailClient)
	mc.SmtpServer = "172.19.3.75"
	// 172.19.3.75  appsmtpr1.true.th
	mc.SmtpPortNo = "25"
	mc.Username = "noreply@truecorp.co.th"
	mc.CryptoProtocol = "TLS"

	mailMessage := mc.NewMessage("Test subject TLS", "Test body TLS")
	mailMessage.To = []string{"paravit_tun@truecorp.co.th"}

	err := mc.SendEmail(mailMessage)

	if err != nil {
		fmt.Println("SendEmail fail", err.Error())
		t.Errorf("TestSendEmailTLS Error %v", err.Error())
	} else {
		fmt.Println("TestSendEmailTLS success")
	}
}
