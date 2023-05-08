package main

import (
	"ta_microservice_peminjaman/app/controllers"
	"ta_microservice_peminjaman/app/routers"
)

func main() {

	packagePeminjam := controllers.NewPeminjaman()

	r := routers.SetupRouter(*packagePeminjam)
	_ = r.Run(":8070")
}
