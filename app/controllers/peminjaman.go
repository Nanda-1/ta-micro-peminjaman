package controllers

import (
	"net/smtp"
	"os"
	"strings"
	"ta_microservice_peminjaman/app/models"
	"ta_microservice_peminjaman/db"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PeminjamanRepo struct {
	Db *gorm.DB
}

func NewPeminjaman() *PeminjamanRepo {
	db := db.InitDb()
	db.AutoMigrate(&models.Users{}, &models.PeminjamanDetail{})
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
	phoneNumber := c.PostForm("phone_number")
	initialDay := c.PostForm("initial_day")
	startDateStr := c.PostForm("start_date")
	endDateStr := c.PostForm("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		// handle the error if the startDateStr is not in the correct format
		res.Success = false
		ErrMSG := err.Error()
		res.Error = &ErrMSG
		c.JSON(400, res)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		// handle the error if the endDateStr is not in the correct format
		res.Success = false
		ErrMSG := err.Error()
		res.Error = &ErrMSG
		c.JSON(400, res)
		return
	}

	user := models.Users{
		Name:          name,
		Email:         email,
		AsalOrganisai: asalOrganisai,
		PhoneNumber:   phoneNumber,
		Peminjamans: []models.PeminjamanDetail{
			{
				InitialDate: initialDay,
				StartDate:   startDate,
				EndDate:     endDate,
				File:        file.Filename,
				Created_at:  time.Time{},
				Updated_at:  time.Time{},
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

	// result := map[string]interface{}{
	// 	"User": err,
	// }

	res.Data = user
	c.JSON(200, res)

}

func (repo *PeminjamanRepo) GetFile(c *gin.Context) {
	// res := models.JsonResponse{Success: true}

	// id := c.Param("id")
	// if id == "" {
	// 	res.Success = false
	// 	errMsg := "Invalid user ID"
	// 	res.Error = &errMsg
	// 	c.JSON(http.StatusBadRequest, res)
	// 	return
	// }

	// idInt, err := strconv.Atoi(id)
	// if err != nil {
	// 	MsgErr := err.Error()
	// 	res.Success = false
	// 	res.Error = &MsgErr
	// 	c.JSON(http.StatusBadRequest, res)
	// 	return
	// }

	// file, err := models.GetFile(idInt)
	// if err != nil {
	// 	MsgErr := err.Error()
	// 	res.Success = false
	// 	res.Error = &MsgErr
	// 	c.JSON(http.StatusBadRequest, res)
	// 	return
	// }

	// // Set the response headers
	// c.Header("Content-Type", "application/pdf")
	// c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))

	// // Create a new PDF document
	// pdf := gofpdf.New("P", "mm", "A4", "")
	// pdf.AddPage()

	// // Register the file data as a PDF image
	// pdf.RegisterImageOptionsReader("", gofpdf.ImageOptions{
	// 	ImageType: "pdf",
	// }, bytes.NewReader(file))

	// // Use the image on the PDF document
	// pdf.ImageOptions("", 0, 0, 0, 0, false, gofpdf.ImageOptions{}, 0, "")

	// // Generate the PDF output
	// var buffer bytes.Buffer
	// err = pdf.Output(&buffer)
	// if err != nil {
	// 	MsgErr := err.Error()
	// 	res.Success = false
	// 	res.Error = &MsgErr
	// 	c.JSON(http.StatusBadRequest, res)
	// 	return
	// }

	// // Send the PDF as the response
	// c.Data(http.StatusOK, "application/pdf", buffer.Bytes())

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
