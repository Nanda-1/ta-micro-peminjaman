package main

import (
	"ta_microservice_peminjaman/app/controllers"
	"ta_microservice_peminjaman/app/routers"

	"github.com/robfig/cron"
)

func main() {
	packagePeminjam := controllers.NewPeminjaman()
	r := routers.SetupRouter(*packagePeminjam)

	c := cron.New()
	_ = c.AddFunc("0 12 * * *", func() {
		controllers.SendReminderEmails(packagePeminjam) // Memanggil SendReminderEmails dengan instance PeminjamanRepo
	})
	c.Start()

	_ = r.Run(":8070")
}
