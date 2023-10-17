package handler

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

func SendMailHandler(w http.ResponseWriter, r *http.Request) {
	allowedMimeTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/pdf": true,
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if _, ok := allowedMimeTypes[fileHeader.Header.Get("Content-Type")]; !ok {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	emails := r.FormValue("emails")
	emailList := strings.Split(emails, ",")

	boundary := multipart.NewWriter(nil).Boundary()
	var message bytes.Buffer
	headers := map[string]string{
		"From":         os.Getenv("SMTP_EMAIL"),
		"To":           strings.Join(emailList, ", "),
		"Subject":      "Attached File",
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", boundary),
	}

	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")

	writer := multipart.NewWriter(&message)
	writer.SetBoundary(boundary)
	part, err := writer.CreateFormFile("attachment", fileHeader.Filename)
	if err != nil {
		http.Error(w, "Error attaching the file", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		http.Error(w, "Error writing to buffer", http.StatusInternalServerError)
		return
	}

	writer.Close()

	smtpServer := "smtp.mail.ru"
	smtpPort := "25"
	auth := smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), smtpServer)
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, os.Getenv("SMTP_EMAIL"), emailList, message.Bytes())
	if err != nil {
		http.Error(w, "Error sending the email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
