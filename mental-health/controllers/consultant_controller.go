package controllers

import (
	"encoding/json"
	"mental-health/config"
	"mental-health/models"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateConsultant(w http.ResponseWriter, r *http.Request) {

	var consultant models.KonsultanKontak
	err := json.NewDecoder(r.Body).Decode(&consultant)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Insert consultant into database
	_, err = config.DB.Exec("INSERT INTO konsultan_kontak (nama, spesialisasi, pengalaman, no_telepon, email) VALUES (?, ?, ?, ?, ?)",
		consultant.Nama, consultant.Spesialisasi, consultant.Pengalaman, consultant.NoTelepon, consultant.Email)
	if err != nil {
		http.Error(w, "Failed to create consultant", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(consultant)

}

func GetAllConsultants(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT konsultan_id, nama, spesialisasi, pengalaman, no_telepon, email FROM konsultan_kontak")
	if err != nil {

		http.Error(w, "Failed to fetch consultants", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var consultants []models.KonsultanKontak

	for rows.Next() {
		var c models.KonsultanKontak
		err := rows.Scan(&c.ID, &c.Nama, &c.Spesialisasi, &c.Pengalaman, &c.NoTelepon, &c.Email)
		if err != nil {
			http.Error(w, "Failed to parse consultant data", http.StatusInternalServerError)
			return
		}
		consultants = append(consultants, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consultants)
}

func DeleteConsultant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Eksekusi query DELETE
	result, err := config.DB.Exec("DELETE FROM konsultan_kontak WHERE konsultan_id = ?", id)
	if err != nil {
		http.Error(w, "Gagal menghapus data konsultan", http.StatusInternalServerError)
		return
	}

	// Cek apakah ada baris yang dihapus
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Konsultan tidak ditemukan", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Konsultan berhasil dihapus"})
}

func UpdateConsultant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedData models.KonsultanKontak
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Eksekusi query update
	result, err := config.DB.Exec(`
		UPDATE konsultan_kontak 
		SET nama = ?, spesialisasi = ?, pengalaman = ?, no_telepon = ?, email = ? 
		WHERE konsultan_id = ?`,
		updatedData.Nama, updatedData.Spesialisasi, updatedData.Pengalaman, updatedData.NoTelepon, updatedData.Email, id)
	if err != nil {
		http.Error(w, "Gagal mengupdate data konsultan", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Konsultan tidak ditemukan", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Konsultan berhasil diupdate"})
}

func GetConsultantByID(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Query ke DB
	row := config.DB.QueryRow("SELECT id, nama, spesialisasi, pengalaman, no_telepon, email FROM konsultan_kontak WHERE id = ?", id)

	var consultant models.KonsultanKontak
	err := row.Scan(&consultant.ID, &consultant.Nama, &consultant.Spesialisasi, &consultant.Pengalaman, &consultant.NoTelepon, &consultant.Email)
	if err != nil {
		http.Error(w, "Konsultan tidak ditemukan", http.StatusNotFound)
		return
	}

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consultant)
}
