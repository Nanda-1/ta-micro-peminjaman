package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"ta_microservice_peminjaman/app/helper"
	"ta_microservice_peminjaman/app/models"
	"ta_microservice_peminjaman/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeminjamanRepo struct {
	Db *gorm.DB
}

func NewPeminjaman() *PeminjamanRepo {
	db := db.InitDb()
	db.AutoMigrate(&models.Users{}, &models.PeminjamanDetail{}, &models.FileDetail{})
	return &PeminjamanRepo{Db: db}
}

func (repo *PeminjamanRepo) CreatePeminjam(c *gin.Context) {
	res := models.JsonResponse{Success: true}

	name := c.PostForm("name")
	email := c.PostForm("email")
	asalOrganisai := c.PostForm("asal_organisai")
	phoneNumber := c.PostForm("phone_number")
	initialDay := c.PostForm("initial_day")
	startDateStr := c.PostForm("start_date")
	endDateStr := c.PostForm("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		MsgErr := err.Error()
		res.Success = false
		res.Error = &MsgErr
		c.JSON(http.StatusBadRequest, res)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		MsgErr := err.Error()
		res.Success = false
		res.Error = &MsgErr
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Membuat direktori penyimpanan file jika belum ada
	fileDir := "./uploads"
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, os.ModePerm)
	}

	// Memproses setiap file yang diunggah
	for _, file := range c.Request.MultipartForm.File["file"] {
		// Menginisialisasi variabel untuk menyimpan detail file
		var fileDetails []models.FileDetail

		fileUUID := uuid.New().String() // Membuat UUID unik untuk setiap file
		filename := fileUUID + "_" + file.Filename
		filename = strings.ReplaceAll(filename, " ", "_") // Mengganti spasi dengan garis bawah
		fileDirClean := filepath.Clean(fileDir)
		filePath := filepath.Join(fileDirClean, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			MsgErr := err.Error()
			res.Success = false
			res.Error = &MsgErr
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		fileDetails = append(fileDetails, models.FileDetail{
			FilePath:       filePath,
			FileName:       filename,
			UploadComplete: true,
			Created_at:     time.Now(),
			Updated_at:     time.Now(),
		})

		user := models.Users{
			Name:          name,
			Email:         email,
			AsalOrganisai: asalOrganisai,
			PhoneNumber:   phoneNumber,
			Peminjamans: []models.PeminjamanDetail{
				{
					InitialDay:  initialDay,
					StartDate:   startDate,
					EndDate:     endDate,
					FileDetails: fileDetails,
					Created_at:  time.Now(),
					Updated_at:  time.Now(),
				},
			},
		}

		// Menyimpan data user ke dalam database
		if err := repo.Db.Create(&user).Error; err != nil {
			MsgErr := err.Error()
			res.Success = false
			res.Error = &MsgErr
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		send := helper.SendMail(email)
		if send != nil {
			res.Success = false
			MsgErr := err.Error()
			res.Error = &MsgErr
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		res.Data = user
	}

	c.JSON(http.StatusOK, res)

}

func (repo *PeminjamanRepo) GetFile(c *gin.Context) {
	res := models.JsonResponse{Success: true}
	peminjamanDetailID := c.Query("peminjaman_detail_id")

	// Query database untuk mendapatkan file
	var file models.FileDetail
	err := repo.Db.Where("peminjaman_detail_id = ?", peminjamanDetailID).First(&file).Error
	if err != nil {
		res.Success = false
		MsgErr := err.Error()
		res.Error = &MsgErr
		c.JSON(http.StatusInternalServerError, res)
		return

	}

	// Membaca file dari sistem penyimpanan
	data, err := ioutil.ReadFile(file.FilePath)
	if err != nil {
		res.Success = false
		MsgErr := err.Error()
		res.Error = &MsgErr
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Mengatur header untuk men-download file
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
	c.Data(http.StatusOK, "application/octet-stream", data)

}

func (repo *PeminjamanRepo) GetAll(c *gin.Context) {
	res := models.JsonResponse{Success: true}

	Get, err := models.GetAllUsers()
	if err != nil {
		MsgErr := err.Error()
		res.Success = false
		res.Error = &MsgErr
		c.JSON(500, res)
		return
	}

	res.Data = Get
	c.JSON(200, res)
}

func (repo *PeminjamanRepo) CountBorrower(c *gin.Context) {
	res := models.JsonResponse{Success: true}

	count, err := models.CountPeminjam()
	if err != nil {
		MsgErr := err.Error()
		res.Success = false
		res.Error = &MsgErr
		c.JSON(400, res)
		return
	}

	res.Data = count
	c.JSON(200, res)
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
	auth := smtp.PlainAuth("", "Impeesa@perbanas.id", "Sukamantri123", "smtp.gmail.com")

	err = smtp.SendMail(host+":"+port, auth, "Impeesa@perbanas.id", to, []byte(body))
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
