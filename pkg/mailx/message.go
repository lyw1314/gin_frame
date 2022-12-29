// author by lipengfei5 @2022-03-28

package mailx

import (
	"io"

	"gopkg.in/gomail.v2"
)

type Message struct {
	*gomail.Message
}

// 新建消息
func NewMessage() *Message {
	return &Message{gomail.NewMessage(gomail.SetEncoding(gomail.Base64))}
}

// 收件人地址
func (msg *Message) To(tomails []string) *Message {
	msg.SetHeader("To", tomails...)
	return msg
}

// 邮件主题
func (msg *Message) Subject(subject string) *Message {
	msg.SetHeader("Subject", subject)
	return msg
}

// 邮件体
func (msg *Message) Body(body string) *Message {
	msg.SetBody("text/html", body)
	return msg
}

// 发件人地址
func (msg *Message) From(addr string, name string) *Message {
	msg.SetAddressHeader("From", addr, name)
	return msg
}

// 抄送
func (msg *Message) Cc(ccmails []string) *Message {
	msg.SetHeader("Cc", ccmails...)
	return msg
}

// 密抄
func (msg *Message) Bcc(bccmails []string) *Message {
	msg.SetHeader("Bcc", bccmails...)
	return msg
}

// 可能出现发送时删除的情况
func (msg *Message) AttachReader(name string, r io.Reader, settings ...gomail.FileSetting) *Message {
	msg.Message.Attach(name, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := io.Copy(w, r)
		return err
	}))
	return msg
}
