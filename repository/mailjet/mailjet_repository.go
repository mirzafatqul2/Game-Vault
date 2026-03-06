package mailjet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/pobyzaarif/goshortcute"
)

type MailjetConfig struct {
	MailjetBaseURL           string
	MailjetBasicAuthUsername string
	MailjetBasicAuthPassword string
	MailjetSenderEmail       string
	MailjetSenderName        string
}

type MailjetRepository struct {
	logger *slog.Logger
	cfg    MailjetConfig
}

func NewMailjetRepository(logger *slog.Logger, cfg MailjetConfig) *MailjetRepository {
	return &MailjetRepository{logger: logger, cfg: cfg}
}

type payloadSendEmail struct {
	Messages []Messages `json:"Messages"`
}
type From struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}
type To struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}
type Messages struct {
	From     From   `json:"From"`
	To       []To   `json:"To"`
	Subject  string `json:"Subject"`
	TextPart string `json:"TextPart"`
	HTMLPart string `json:"HTMLPart"`
}

func (r *MailjetRepository) SendEmail(toName, toEmail, subject, message string) (err error) {
	url := r.cfg.MailjetBaseURL + "/v3.1/send"
	method := http.MethodPost

	toBody := []To{}
	toBody = append(toBody, To{
		Email: toEmail,
		Name:  toName,
	})

	messageBody := Messages{
		To: toBody,
		From: From{
			Email: r.cfg.MailjetSenderEmail,
			Name:  r.cfg.MailjetSenderName,
		},
		Subject:  subject,
		TextPart: message,
		HTMLPart: message,
	}
	constructMessages := []Messages{}
	constructMessages = append(constructMessages, messageBody)

	payload := payloadSendEmail{
		Messages: constructMessages,
	}

	payloadByte, _ := json.Marshal(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadByte))
	if err != nil {
		r.logger.Error("failed create mailjet request", "error", err)
		return err
	}

	buildABasicAuth := goshortcute.StringtoBase64Encode(r.cfg.MailjetBasicAuthUsername + ":" + r.cfg.MailjetBasicAuthPassword)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+buildABasicAuth)

	res, err := client.Do(req)
	if err != nil {
		r.logger.Error("mailjet request failed", "error", err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("mailjet returned status %d", res.StatusCode)
	}

	r.logger.Info("email sent successfully",
		"to", toEmail,
		"status", res.StatusCode,
	)

	return nil
}
