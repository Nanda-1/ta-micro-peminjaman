package helper

import (
	"net/smtp"
	"os"
	"strings"
	"time"
)

func SendMail(To string, name string) error {
	to := []string{To}
	msg := "Halo, " + name + "\n\n" +
		"Terima kasih telah melakukan pengajuan peminjaman alat. Data Anda telah berhasil kami terima dan sedang dalam proses verifikasi oleh pihak sekretariat Mapala Impeesa Perbanas.\n\n" +
		"Mohon menunggu konfirmasi lebih lanjut melalui email ini setelah tim kami meninjau pengajuan Anda. Jika ada informasi tambahan yang diperlukan, kami akan menghubungi Anda melalui kontak yang telah Anda berikan.\n\n" +
		"Terima kasih atas perhatiannya.\n\n" +
		"Salam,\n" +
		"Tim Sekretariat Mapala Impeesa Perbanas"
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

func SendMailAproved(To string, name string) error {
	to := []string{To}
	msg :=
		`Halo ` + name + `,

Terima kasih telah melakukan pengajuan peminjaman alat. Data Anda telah berhasil kami terima dan telah disetujui oleh pihak sekretariat Mapala Impeesa Perbanas.

Silakan datang ke sekretariat untuk mengambil barang yang Anda ajukan.

Terima kasih atas perhatiannya.

Salam,
Tim Sekretariat Mapala Impeesa Perbanas`

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

func SendMailRejected(To string, name string) error {
	to := []string{To}
	msg :=
		`Halo, ` + name + `

Mohon maaf, pengajuan peminjaman alat Anda telah ditolak oleh pihak sekretariat Mapala Impeesa Perbanas.

Jika Anda memiliki pertanyaan lebih lanjut atau memerlukan informasi tambahan, silakan hubungi kami melalui kontak yang telah Anda berikan.

Terima kasih atas perhatiannya.

Salam,
Tim Sekretariat Mapala Impeesa Perbanas`
	subject := "Pengajuan peminjaman ditolak"
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
		return err
	}

	return nil
}

func SendReminderEmail(To string, name string, tomorrow time.Time) error {
	to := []string{To}
	msg :=
		`Halo ` + name + `,

Ini adalah pengingat bahwa besok, ` + tomorrow.Format("02 January 2006") + `, adalah batas waktu pengembalian alat yang Anda pinjam.

Mohon segera mengembalikan alat tersebut ke sekretariat Mapala Impeesa Perbanas.

Terima kasih atas perhatiannya.

Salam,
Tim Sekretariat Mapala Impeesa Perbanas`
	subject := "Pengingat Pengembalian Alat"
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
		return err
	}

	return nil
}


func SendOverdueEmail(To string, name string, endDate time.Time) error {
	to := []string{To}
	msg :=
		`Halo ` + name + `,

Mohon maaf, pengembalian alat yang Anda pinjam seharusnya sudah dilakukan pada ` + endDate.Format("02 January 2006") + `. Namun, hingga saat ini alat tersebut belum dikembalikan.

Kami harap Anda segera mengembalikan alat tersebut ke sekretariat Mapala Impeesa Perbanas.

Terima kasih atas perhatiannya.

Salam,
Tim Sekretariat Mapala Impeesa Perbanas`
	subject := "Pengingat Pengembalian Alat Telah Melewati Batas Waktu"
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
		return err
	}

	return nil
}
