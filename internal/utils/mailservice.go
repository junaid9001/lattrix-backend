package utils

// "errors"
// "fmt"
// "log"
// "net/smtp"

// "github.com/junaid9001/lattrix-backend/internal/config"

func SendEmail(to, subject, body string) error {
	// from := config.AppConfig.EMAIL
	// password := config.AppConfig.EMAIL_PASSWORD

	// if from == "" || password == "" {
	// 	log.Println("Error: Email environment variables not set")
	// 	return errors.New("email configuration missing")
	// }

	// smtpHost := "smtp.gmail.com"
	// smtpPort := "587"

	// senderName := "Lattrix"
	// fromHeader := fmt.Sprintf("%s <%s>", senderName, from)

	// msg := []byte(
	// 	"From: " + fromHeader + "\r\n" +
	// 		"To: " + to + "\r\n" +
	// 		"Subject: " + subject + "\r\n" +
	// 		"MIME-Version: 1.0\r\n" +
	// 		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
	// 		"\r\n" +
	// 		body + "\r\n",
	// )

	// auth := smtp.PlainAuth("", from, password, smtpHost)
	// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	// if err != nil {
	// 	return err
	// }

	return nil
}
