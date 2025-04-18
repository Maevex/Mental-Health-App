package models

type SesiKonsultasi struct {
	SesiID     int    `json:"sesi_id"`
	UserID     int    `json:"user_id"`
	Status     string `json:"status"`
	WaktuMulai string `json:"waktu_mulai"`
}
