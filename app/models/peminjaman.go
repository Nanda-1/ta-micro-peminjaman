package models

import (
	"errors"
	"ta_microservice_peminjaman/db"
	"time"
)

type Users struct {
	Id            int          `json:"id" gorm:"primaryKey"`
	Name          string       `json:"name"`
	Email         string       `json:"email" gorm:"unique"`
	AsalOrganisai string       `json:"asal_organisai"`
	File          string       `json:"file"`
	Peminjamans   []Peminjaman `json:"peminjamans" gorm:"foreignKey:UserId"`
	Created_at    time.Time    `json:"created_at"`
	Updated_at    time.Time    `json:"updated_at"`
}

type Peminjaman struct {
	Id             int       `json:"id" gorm:"primarykey"`
	UserId         int       `json:"user_id" `
	Tanggal_pinjam string    `json:"tanggal_pinjam"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
}

func (Peminjaman) TableName() string {
	return "peminjamans"
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
