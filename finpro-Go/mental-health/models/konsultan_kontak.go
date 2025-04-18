package models

type KonsultanKontak struct {
	ID           int    `json:"id"`
	Nama         string `json:"nama"`
	Spesialisasi string `json:"spesialisasi"`
	Pengalaman   string `json:"pengalaman"`
	NoTelepon    string `json:"no_telepon"`
	Email        string `json:"email"`
}
