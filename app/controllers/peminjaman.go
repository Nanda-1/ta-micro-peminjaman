package controllers

import (
	"net/smtp"
	"os"
	"strings"
	"ta_microservice_peminjaman/app/models"
	"ta_microservice_peminjaman/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PeminjamanRepo struct {
	Db *gorm.DB
}

func NewPeminjaman() *PeminjamanRepo {
	db := db.InitDb()
	db.AutoMigrate(&models.Users{}, &models.Peminjaman{})
	return &PeminjamanRepo{Db: db}
}

func (repo *PeminjamanRepo) CreatePeminjam(c *gin.Context) {
	res := models.JsonResponse{Success: true}

	file, err := c.FormFile("file")
	if err != nil {
		res.Success = false
		ErrMSG := err.Error()
		res.Error = &ErrMSG
		c.JSON(400, res)
		return
	}

	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", os.ModePerm)
	}

	if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
		res.Success = false
		ErrMSG := err.Error()
		res.Error = &ErrMSG
		c.JSON(500, res)
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")
	asalOrganisai := c.PostForm("asal_organisai")
	tanggalPinjamStr := c.PostForm("tanggal_pinjam")
	// tanggalPinjam, err := time.Parse(time.RFC3339, tanggalPinjamStr)
	// if err != nil {
	// 	res.Success = false
	// 	ErrMSG := err.Error()
	// 	res.Error = &ErrMSG
	// 	c.JSON(500, res)
	// 	return
	// }

	user := models.Users{
		Name:          name,
		Email:         email,
		AsalOrganisai: asalOrganisai,
		File:          file.Filename,
		Peminjamans: []models.Peminjaman{
			{
				// UserId:         user.Id,
				Tanggal_pinjam: tanggalPinjamStr,
			},
		},
	}

	// Menyimpan data user ke dalam database
	if err := repo.Db.Create(&user).Error; err != nil {
		res.Success = false
		ErrMSG := err.Error()
		res.Error = &ErrMSG
		c.JSON(500, res)
		return
	}

	// pinjam := models.Peminjaman{
	// 	UserId:         user.Id,
	// 	Tanggal_pinjam: tanggalPinjamStr,
	// }

	// // Menyimpan data peminjaman ke dalam database
	// if err := repo.Db.Create(&pinjam).Error; err != nil {
	// 	res.Success = false
	// 	ErrMSG := err.Error()
	// 	res.Error = &ErrMSG
	// 	c.JSON(500, res)
	// 	return
	// }

	result := map[string]interface{}{
		"User": user,
	}

	res.Data = result
	c.JSON(200, res)

}

func (repo *PeminjamanRepo) GetFile(c *gin.Context) {
	// res := models.JsonResponse{Success: true}

}

func (repo *PeminjamanRepo) GetAll(c *gin.Context) {

}

type NotificationRequest struct {
	To string `json:"to"`
}

func (repo *PeminjamanRepo) SendEmail(c *gin.Context) {
	res := models.JsonResponse{Success: true}

	var reqBody NotificationRequest
	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		res.Success = false
		ErrMSG := err.Error()
		res.Error = &ErrMSG
		c.JSON(400, res)
		return
	}

	to := []string{reqBody.To}
	msg := "Pengajuan peminjaman alat anda sudah disetujui segera datang kembali ke sekretariat mapala impeesa untuk mengambil barang"
	subject := "Pengajuan peminjaman"
	body := "From: " + os.Getenv("CONFIG_SENDER_NAME") + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		msg

	host := "smtp.gmail.com"
	port := "587"
	auth := smtp.PlainAuth("", "rinandasoe@gmail.com", "cggudwsusrzaxzfu", "smtp.gmail.com")

	err = smtp.SendMail(host+":"+port, auth, "rinandasoe@gmail.com", to, []byte(body))
	if err != nil {
		res.Success = false
		ErrMSG := err.Error()
		res.Error = &ErrMSG
		c.JSON(400, res)
		return
	}

	res.Data = "Email Berhasil dikirim"

	c.JSON(200, res)
}
