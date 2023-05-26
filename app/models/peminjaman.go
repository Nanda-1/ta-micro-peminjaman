package models

import (
	"errors"
	"io/ioutil"
	"ta_microservice_peminjaman/db"
	"time"
)

type Users struct {
	Id            int                `json:"id" gorm:"primaryKey"`
	Name          string             `json:"name"`
	Email         string             `json:"email" gorm:"unique"`
	AsalOrganisai string             `json:"asal_organisai"`
	PhoneNumber   string             `json:"phone_number"`
	Peminjamans   []PeminjamanDetail `json:"peminjamans" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserId"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
}

type PeminjamanDetail struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	UserId      int       `json:"user_id"`
	InitialDate string    `json:"initial_date" `
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	File        string    `json:"file"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

func (PeminjamanDetail) TableName() string {
	return "peminjaman_details"
}

func CreatePeminjam(user *Users) (*Users, error) {
	result := db.Db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func GetPeminjamByID(id int) (*Users, error) {
	Peminjaman := Users{}
	db.Db.Where("id", id).First(&Peminjaman)
	if Peminjaman.Id == 0 {
		return nil, errors.New("peminjaman tidak ditemukan.")
	}
	return &Peminjaman, nil
}
func GetAllUsers() ([]Users, error) {
	var users []Users
	result := db.Db.Preload("Peminjamans").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func GetFile(id int) ([]byte, error) {
	var file PeminjamanDetail
	result := db.Db.Where("id = ? ", id).Order("id DESC").First(&file)
	if result.Error != nil {
		return nil, result.Error
	}

	// Read the file from disk or any storage location
	// and return it as a byte slice
	fileData, err := ioutil.ReadFile(file.File)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}
