package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

// Fungsi untuk memeriksa jenis file yang diunggah

func (repo *PeminjamanRepo) CreatePeminjam(c *gin.Context) {
	res := models.JsonResponse{Success: true}

	name := c.PostForm("name")
	email := c.PostForm("email")
	asalOrganisai := c.PostForm("asal_organisasi")
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
	var fileDetails []models.FileDetail
	for _, fileHeader := range c.Request.MultipartForm.File["file"] {
		// Batasan jenis file yang diizinkan
		if !helper.IsFileTypeAllowed(fileHeader) {
			MsgErr := "File type not allowed. Only PDF files are allowed."
			res.Success = false
			res.Error = &MsgErr
			c.JSON(http.StatusBadRequest, res)
			return
		}

		// Batasan ukuran file
		if !helper.IsFileSizeAllowed(fileHeader) {
			MsgErr := "File size exceeds the limit. Maximum file size allowed is 5MB."
			res.Success = false
			res.Error = &MsgErr
			c.JSON(http.StatusBadRequest, res)
			return
		}

		file := fileHeader
		fileUUID := uuid.New().String() // Membuat UUID unik untuk setiap file
		filename := fileUUID + "_" + file.Filename
		filename = strings.ReplaceAll(filename, " ", "_") // Mengganti spasi dengan garis bawah
		filePath := filepath.Join(fileDir, filename)

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
	}

	user := models.Users{
		Name:          name,
		Email:         email,
		AsalOrganisai: asalOrganisai,
		PhoneNumber:   phoneNumber,
		Status:        "pending",
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
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	// Menyimpan data user ke dalam database
	if err := repo.Db.Create(&user).Error; err != nil {
		MsgErr := err.Error()
		res.Success = false
		res.Error = &MsgErr
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	send := helper.SendMail(user.Email, user.Name)
	if send != nil {
		res.Success = false
		MsgErr := err.Error()
		res.Error = &MsgErr
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = user
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

func (repo *PeminjamanRepo) SendApproved(c *gin.Context) {
	res := models.JsonResponse{Success: true}

	userID := c.Query("id")
	status := c.Query("status")

	user := models.Users{}
	err := repo.Db.First(&user, userID).Error
	if err != nil {
		res.Success = false
		ErrMSG := err.Error()
		res.Error = &ErrMSG
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if status == "approve" {
		user.Status = "approve"
		send := helper.SendMailAproved(user.Email, user.Name)
		if send != nil {
			res.Success = false
			MsgErr := send.Error()
			res.Error = &MsgErr
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	} else if status == "reject" {
		user.Status = "reject"
		send := helper.SendMailRejected(user.Email, user.Name)
		if send != nil {
			res.Success = false
			MsgErr := send.Error()
			res.Error = &MsgErr
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	err = repo.Db.Save(&user).Error
	if err != nil {
		res.Success = false
		ErrMSG := err.Error()
		res.Error = &ErrMSG
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	result := map[string]interface{}{
		"Message": "Data berhasil diperbarui",
		"status":  user.Status,
	}

	res.Data = result
	c.JSON(http.StatusOK, res)

}

func SendReminderEmails(repo *PeminjamanRepo) {
	res := models.JsonResponse{Success: true}

	// Mengambil seluruh data pengguna
	users, err := models.GetAllUsers()
	if err != nil {
		res.Success = false
		ErrMSG := "Error GetAllUsers"
		res.Error = &ErrMSG
		fmt.Println(res)
		return
	}

	now := time.Now()

	for _, user := range users {
		var peminjamanDetail models.PeminjamanDetail

		if len(user.Peminjamans) > 0 {
			// Mengambil PeminjamanDetail dari slice pertama
			peminjamanDetail = user.Peminjamans[0]
		} else {
			// Menghandle jika tidak ada PeminjamanDetail yang ditemukan
			res.Success = false
			ErrMSG := "PeminjamanDetail not found"
			res.Error = &ErrMSG
			fmt.Println(res)
			return // Melanjutkan ke pengguna berikutnya jika tidak ada PeminjamanDetail
		}

		// Menghitung selisih waktu antara sekarang dan EndDate
		timeDiff := peminjamanDetail.EndDate.Sub(now)

		// Mengatur durasi reminder sebelum EndDate (misalnya 1 hari sebelumnya)
		reminderDuration := 24 * time.Hour

		if timeDiff > reminderDuration {
			reminderTime := peminjamanDetail.EndDate.Add(-reminderDuration)

			// Mengirim email reminder jika waktu reminder lebih besar dari selisih waktu
			err := helper.SendReminderEmail(user.Email, user.Name, reminderTime)
			if err != nil {
				res.Success = false
				MsgErr := err.Error()
				res.Error = &MsgErr
				fmt.Println(res)
				return // Melanjutkan ke pengguna berikutnya jika terjadi error saat mengirim email
			}

			fmt.Printf("Reminder email sent for user ID %s\n", user.Id, user.Name)
		} else if now.After(peminjamanDetail.EndDate) {
			// Kasus ketika EndDate sudah lewat
			err := helper.SendOverdueEmail(user.Email, user.Name, peminjamanDetail.EndDate)
			if err != nil {
				res.Success = false
				MsgErr := err.Error()
				res.Error = &MsgErr
				fmt.Println(res)
				return // Melanjutkan ke pengguna berikutnya jika terjadi error saat mengirim email
			}

			fmt.Printf("Overdue email sent for user ID %s\n", user.Id, user.Name)
		}
	}
}
