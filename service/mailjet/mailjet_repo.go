package mailjet

type MailjetRepository interface {
	SendEmail(toEmail, toName, subject, message string) error
}
