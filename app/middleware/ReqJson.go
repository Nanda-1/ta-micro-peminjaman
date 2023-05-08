package middleware

import (
	"net/http"
	"ta_microservice_peminjaman/app/models"

	"github.com/gin-gonic/gin"
)

func ReqJson() gin.HandlerFunc {

	return func(c *gin.Context) {
		res := models.JsonResponse{Success: true}
		if c.Request.Header.Get("Content-Type") != "application/json" {
			res.Success = false
			x := "accepted content-type: application/json"
			res.Error = &x
			c.JSON(http.StatusBadRequest, res)
			c.Abort()
			return
		}
	}

}
