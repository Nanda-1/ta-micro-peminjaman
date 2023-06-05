package helper

import (
	"net/smtp"
	"os"
	"strings"
)

func SendMail(To string) error {
	to := []string{To}
	msg := "Pengajuan peminjaman alat anda sudah disetujui segera datang kembali ke sekretariat mapala impeesa untuk mengambil barang"
	subject := "Pengajuan peminjaman"
	body := "From: " + os.Getenv("CONFIG_SENDER_NAME") + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		msg

	host := "smtp.gmail.com"
	port := "587"
	auth := smtp.PlainAuth("", "Impeesa@perbanas.id", "Sukamantri123", "smtp.gmail.com")

	err := smtp.SendMail(host+":"+port, auth, "Impeesa@perbanas.id", to, []byte(body))
	if err != nil {
		return nil
	}

	return nil
}
