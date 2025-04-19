package models

type ResponsChatbot struct {
	ResponsID  int    `json:"respons_id"`
	PesanID    int    `json:"pesan_id"`
	IsiRespons string `json:"isi_respons"`
}
