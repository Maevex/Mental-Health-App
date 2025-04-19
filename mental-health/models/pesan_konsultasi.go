package models

type PesanKonsultasi struct {
	PesanID        int     `json:"pesan_id"`
	SesiID         int     `json:"sesi_id"`
	IsiPesan       string  `json:"isi_pesan"`
	SentimentScore float64 `json:"sentiment_score"`
	UrgencyLevel   string  `json:"urgency_level"`
}
