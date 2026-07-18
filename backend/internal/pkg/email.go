package pkg

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/sheet-platform/backend/internal/config"
)

func SendVerificationEmail(emailCfg *config.EmailConfig, to, code string) error {
	subject := "注册 - 验证码"
	body := fmt.Sprintf("您的验证码是：%s，有效期5分钟。请勿将验证码泄露给他人。", code)

	msg := buildMessage(emailCfg.Username, to, subject, body)

	addr := fmt.Sprintf("%s:%d", emailCfg.Host, emailCfg.Port)

	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: emailCfg.Host})
	if err != nil {
		return fmt.Errorf("tls dial failed: %w", err)
	}
	defer conn.Close()

	host, _, _ := net.SplitHostPort(addr)
	if host == "" {
		host = emailCfg.Host
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client failed: %w", err)
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", emailCfg.Username, emailCfg.Password, emailCfg.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth failed: %w", err)
	}

	if err := client.Mail(emailCfg.Username); err != nil {
		return fmt.Errorf("smtp mail failed: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt failed: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data failed: %w", err)
	}
	defer wc.Close()

	if _, err := wc.Write([]byte(msg)); err != nil {
		return fmt.Errorf("smtp write failed: %w", err)
	}

	return nil
}

func buildMessage(from, to, subject, body string) string {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	}
	return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}
