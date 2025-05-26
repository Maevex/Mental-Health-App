package models

type KeluhanAwalRequest struct {
	Kategori string `json:"kategori"`
	Keluhan  string `json:"keluhan"`
}

type KeluhanAwal struct {
	ID             int     `json:"id"`
	UserID         int     `json:"user_id"`
	Kategori       string  `json:"kategori"`
	IsiKeluhan     string  `json:"isi_keluhan"`
	SentimentScore float64 `json:"sentiment_score"`
	Kesimpulan     string  `json:"kesimpulan"`
	SubKategori    string  `json:"sub_kategori"`
	Waktu          string  `json:"waktu"`
}
