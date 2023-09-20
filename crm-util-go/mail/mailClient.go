package mail

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"mime/multipart"
	"net/smtp"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type MailClient struct {
	SmtpServer     string
	SmtpPortNo     string
	Username       string
	Password       string
	CryptoProtocol string
}

type Message struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

func (mc *MailClient) sendEmailSSL(m *Message) error {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         mc.SmtpServer,
	}

	conn, err := tls.Dial("tcp", mc.SmtpServer+":"+mc.SmtpPortNo, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	smtpClient, err := smtp.NewClient(conn, mc.SmtpServer)
	if err != nil {
		return err
	}
	defer smtpClient.Close()

	if len(mc.Username) > 0 && len(mc.Password) > 0 {
		auth := smtp.PlainAuth("", mc.Username, mc.Password, mc.SmtpServer)
		smtpClient.Auth(auth)
	}

	// Set the sender
	err = smtpClient.Mail(mc.Username)
	if err != nil {
		return err
	}

	// Set recipient
	for _, toAddr := range m.To {
		err = smtpClient.Rcpt(toAddr)
		if err != nil {
			return err
		}
	}

	// Send the email body.
	writer, err := smtpClient.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write(m.ToBytes())
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	return smtpClient.Quit()
}

func (mc *MailClient) sendEmailTLS(m *Message) error {
	// Connect to the remote SMTP server.
	smtpClient, err := smtp.Dial(mc.SmtpServer + ":" + mc.SmtpPortNo)
	if err != nil {
		return err
	}
	defer smtpClient.Close()

	supportedExt, _ := smtpClient.Extension("STARTTLS")
	if supportedExt {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         mc.SmtpServer,
		}

		err = smtpClient.StartTLS(tlsConfig)
		if err != nil {
			return err
		}
	}

	if len(mc.Username) > 0 && len(mc.Password) > 0 {
		auth := smtp.PlainAuth("", mc.Username, mc.Password, mc.SmtpServer)
		err = smtpClient.Auth(auth)
		if err != nil {
			return err
		}
	}

	// Set the sender
	err = smtpClient.Mail(mc.Username)
	if err != nil {
		return err
	}

	// Set recipient
	for _, toAddr := range m.To {
		err = smtpClient.Rcpt(toAddr)
		if err != nil {
			return err
		}
	}

	// Send the email body.
	writer, err := smtpClient.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write(m.ToBytes())
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	return smtpClient.Quit()
}

func (mc *MailClient) SendEmail(m *Message) error {
	if strings.EqualFold(mc.CryptoProtocol, "SSL") {
		return mc.sendEmailSSL(m)
	} else {
		return mc.sendEmailTLS(m)
	}
}

func (mc *MailClient) NewMessage(subject string, body string) *Message {
	return &Message{Subject: subject, Body: body, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(pathFileName string) error {
	b, err := os.ReadFile(pathFileName)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(pathFileName)
	m.Attachments[fileName] = b
	return nil
}

func getMimeContentType(fileName string) (mimeContentType string) {
	fileExtension := path.Ext(fileName)

	switch fileExtension {
	case ".txt":
		mimeContentType = "text/plain"
	case ".htm", ".html":
		mimeContentType = "text/html"
	case ".xhtml":
		mimeContentType = "application/xhtml+xml"
	case ".xml":
		mimeContentType = "text/xml"
	case ".css":
		mimeContentType = "text/css"
	case ".js":
		mimeContentType = "text/javascript"
	case ".csv":
		mimeContentType = "text/csv"
	case ".bmp":
		mimeContentType = "image/bmp"
	case ".gif":
		mimeContentType = "image/gif"
	case ".ico":
		mimeContentType = "image/vnd.microsoft.icon"
	case ".jpeg", ".jpg":
		mimeContentType = "image/jpeg"
	case ".png":
		mimeContentType = "image/png"
	case ".tif", ".tiff":
		mimeContentType = "image/tiff"
	case ".webp":
		mimeContentType = "image/webp"
	case ".svg":
		mimeContentType = "image/svg+xml"
	case ".pdf":
		mimeContentType = "application/pdf"
	case ".doc":
		mimeContentType = "application/msword"
	case ".docx":
		mimeContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		mimeContentType = "application/vnd.ms-excel"
	case ".xlsx":
		mimeContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		mimeContentType = "application/vnd.ms-powerpoint"
	case ".pptx":
		mimeContentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".vsd":
		mimeContentType = "application/vnd.visio"
	case ".rar":
		mimeContentType = "application/vnd.rar"
	case ".zip":
		mimeContentType = "application/zip"
	case ".7z":
		mimeContentType = "application/x-7z-compressed"
	case ".gz":
		mimeContentType = "application/gzip"
	default:
		contentType, _ := mime.ExtensionsByType(fileExtension)
		if contentType != nil {
			mimeContentType = contentType[0]
		}
	}

	return mimeContentType
}

func (m *Message) ToBytes() []byte {
	withAttachments := len(m.Attachments) > 0

	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\r\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))

		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
		buf.WriteString("Content-Transfer-Encoding: base64\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(base64.StdEncoding.EncodeToString([]byte(m.Body)))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
		buf.WriteString(m.Body)
	}

	if withAttachments {
		for fileName, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", getMimeContentType(fileName)))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)

			buf.WriteString("\r\n")
			buf.Write(b)
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}
