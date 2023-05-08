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
	protectedRouter.Use(middleware.ApiKey(), middleware.Jwt())
	protectedRouter.POST("/create", api.RepoPeminjaman.CreatePeminjam)
	protectedRouter.POST("/email", api.RepoPeminjaman.SendEmail)

	return r
}
