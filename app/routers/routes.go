package routers

import (
	"ta_microservice_peminjaman/app/controllers"
	"ta_microservice_peminjaman/app/middleware"

	"github.com/gin-gonic/gin"
)

type API struct {
	RepoPeminjaman controllers.PeminjamanRepo
}

func SetupRouter(RepoPeminjaman controllers.PeminjamanRepo) *gin.Engine {
	r := gin.New()
	api := API{
		RepoPeminjaman,
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	protectedRouter := r.Group("/api/peminjaman")
	protectedRouter.Use(middleware.ApiKey())
	protectedRouter.POST("/create", api.RepoPeminjaman.CreatePeminjam)
	protectedRouter.GET("/get-borrow", api.RepoPeminjaman.GetAll)
	protectedRouter.GET("/get-file", api.RepoPeminjaman.GetFile)
	protectedRouter.GET("/count", api.RepoPeminjaman.CountBorrower)
	protectedRouter.POST("/email", api.RepoPeminjaman.SendApproved)
	// protectedRouter.POST("/remender", api.RepoPeminjaman.SendReminderEmail)

	return r
}
